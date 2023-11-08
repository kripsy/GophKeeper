//nolint:gosec
package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	Register(login, password string) error
	Login(login, password string) error

	DownloadFile(ctx context.Context, fileName string, fileHash string, syncKey string) (chan []byte, error)
	UploadFile(ctx context.Context, fileName string, hash string, syncKey string, data chan []byte) error

	BlockStore(ctx context.Context, syncKey string, guidChan chan string) error
	ApplyChanges(ctx context.Context, id string) error

	IsNotAvailable() bool
	IsAvailable() bool
	TryToConnect() bool
}

type Grpc struct {
	client        pb.GophKeeperServiceClient
	token         string
	serverAddress string
	conn          *grpc.ClientConn
	log           zerolog.Logger
	isAvailable   bool
}

func NewClient(serverAddress string, log zerolog.Logger) *Grpc {
	return &Grpc{
		serverAddress: serverAddress,
		log:           log,
	}
}

func (c *Grpc) Register(login, password string) error {
	resp, err := c.client.Register(context.Background(), &pb.AuthRequest{
		Username: login,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}
	c.token = resp.GetToken()

	return nil
}

func (c *Grpc) Login(login, password string) error {
	resp, err := c.client.Login(context.Background(), &pb.AuthRequest{
		Username: login,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	c.token = resp.GetToken()

	return nil
}

func (c *Grpc) BlockStore(ctx context.Context, syncKey string, guidChan chan string) error {
	stream, err := c.client.BlockStore(c.getCtx(ctx, c.token))
	if err != nil {
		return fmt.Errorf("BlockStore: %w", err)
	}

	err = stream.Send(&pb.BlockStoreRequest{Guid: syncKey})
	if err != nil {
		c.log.Err(err).Msg("err send block store req")

		return fmt.Errorf("BlockStore.Send: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("BlockStore.Recv: %w", err)
	}

	guidChan <- resp.GetGuid()

	return nil
}

func (c *Grpc) DownloadFile(ctx context.Context,
	fileName string,
	fileHash string,
	syncKey string,
) (chan []byte, error) {
	req := &pb.MultipartDownloadFileRequest{
		FileName: fileName,
		Guid:     syncKey,
		Hash:     fileHash,
	}

	stream, err := c.client.MultipartDownloadFile(c.getCtx(ctx, c.token), req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	data := make(chan []byte, 1)
	go func() {
	loop:
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				close(data)

				break loop
			}
			if err != nil {
				c.log.Err(err).Msg("failed download file")
			}

			data <- resp.GetContent()
		}
	}()

	return data, nil
}

func (c *Grpc) UploadFile(ctx context.Context, fileName string, hash string, syncKey string, data chan []byte) error {
	stream, err := c.client.MultipartUploadFile(c.getCtx(ctx, c.token))
	if err != nil {
		return fmt.Errorf("MultipartUploadFile: %w", err)
	}

	for chunk := range data {
		if err = stream.Send(&pb.MultipartUploadFileRequest{
			FileName: fileName,
			Hash:     hash,
			Guid:     syncKey,
			Content:  chunk,
		}); err != nil && !errors.Is(err, io.EOF) {
			c.log.Err(err).Msg("failed upload")

			return fmt.Errorf("UploadFile: %w", err)
		}
	}
	//nolint:ineffassign,nolintlint
	_, err = stream.CloseAndRecv()
	if err != nil {
		c.log.Err(err).Msg("failed upload")

		return fmt.Errorf("UploadFile: %w", err)
	}

	return nil
}

func (c *Grpc) ApplyChanges(ctx context.Context, id string) error {
	_, err := c.client.ApplyChanges(c.getCtx(ctx, c.token), &pb.ApplyChangesRequest{Guid: id})
	if err != nil {
		return fmt.Errorf("ApplyChanges: %w", err)
	}

	return nil
}

func (c *Grpc) IsNotAvailable() bool {
	return !c.isAvailable
}

func (c *Grpc) IsAvailable() bool {
	return c.isAvailable
}

func (c *Grpc) TryToConnect() bool {
	conn, err := grpc.Dial(c.serverAddress,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		c.log.Info().Err(err).Msg("failed connect to server")

		return false
	}

	c.conn = conn
	c.client = pb.NewGophKeeperServiceClient(conn)

	_, err = c.client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.log.Info().Err(err).Msg("failed connect to server")

		return false
	}

	c.isAvailable = true

	return true
}

func (c *Grpc) getCtx(ctx context.Context, jwt string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": jwt,
	})

	newCtx := metadata.NewOutgoingContext(ctx, md)

	return newCtx
}
