package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

type Client interface {
	Register(ctx context.Context, in *pb.AuthRequest) error
	Login(ctx context.Context, in *pb.AuthRequest) error

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

func (c *Grpc) Register(ctx context.Context, in *pb.AuthRequest) error {
	resp, err := c.client.Register(ctx, in)
	if err != nil {
		return err
	}
	c.token = resp.Token

	return nil
}

func (c *Grpc) Login(ctx context.Context, in *pb.AuthRequest) error {
	resp, err := c.client.Login(ctx, in)
	if err != nil {
		return err
	}
	c.token = resp.Token

	return nil
}

func (c *Grpc) BlockStore(ctx context.Context, syncKey string, guidChan chan string) error {
	stream, err := c.client.BlockStore(c.getCtx(ctx, c.token))
	if err != nil {
		return err
	}

	err = stream.Send(&pb.BlockStoreRequest{Guid: syncKey})
	if err != nil {
		c.log.Err(err).Msg("err send block store req")

		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		fmt.Println("Error receiving response from server:", err.Error())
		return err
	}

	guidChan <- resp.Guid
	//err = stream.CloseSend()
	//if err != nil {
	//	c.log.Err(err).Msg("failed BlockStore CloseSend")
	//	return err
	//}

	return nil
}

func (c *Grpc) DownloadFile(ctx context.Context, fileName string, fileHash string, syncKey string) (chan []byte, error) {
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
			if err == io.EOF {
				close(data)

				break loop
			}
			if err != nil {
				return //nil, fmt.Errorf("%w", err)
			}

			data <- resp.Content

		}
	}()

	return data, err
}

func (c *Grpc) UploadFile(ctx context.Context, fileName string, hash string, syncKey string, data chan []byte) error {
	stream, err := c.client.MultipartUploadFile(c.getCtx(ctx, c.token))
	if err != nil {
		return err
	}
	go func() {
		for chunk := range data {
			if err := stream.Send(&pb.MultipartUploadFileRequest{
				FileName: fileName,
				Hash:     hash,
				Guid:     syncKey,
				Content:  chunk,
			}); err != nil && err != io.EOF {
				c.log.Err(err).Msg("upload")
			}
		}
	}()

	return nil
}

func (c *Grpc) ApplyChanges(ctx context.Context, id string) error {
	_, err := c.client.ApplyChanges(c.getCtx(ctx, c.token), &pb.ApplyChangesRequest{Guid: id})
	if err != nil {
		return err
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
