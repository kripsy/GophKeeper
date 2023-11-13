// Package grpc provides the implementation of the gRPC client for the GophKeeper application.
// It includes methods for user registration, login, file uploading and downloading, and other
// client-server interactions over gRPC.
//
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

// Client interface defines the set of methods that the gRPC client must implement.
// It includes functionalities for user authentication, file operations, and server availability checks.
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

// Grpc struct implements the Client interface and provides methods to interact with the gRPC server.
// It manages server communication for user authentication, file operations, and maintaining the connection state.
type Grpc struct {
	client        pb.GophKeeperServiceClient
	token         string
	serverAddress string
	conn          *grpc.ClientConn
	log           zerolog.Logger
	isAvailable   bool
}

// NewClient creates a new Grpc client instance.
// It initializes the gRPC client with the given server address and logger.
func NewClient(serverAddress string, log zerolog.Logger) *Grpc {
	return &Grpc{
		serverAddress: serverAddress,
		log:           log,
	}
}

// Register handles user registration through the gRPC server.
// It sends a registration request and stores the returned token.
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

// Login handles user login through the gRPC server.
// It sends a login request and stores the returned token.
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

// BlockStore initiates a block store operation on the server.
// It sends the sync key and receives a GUID in response.
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

// DownloadFile handles the file download process from the server.
// It sends a request and streams the file content back to the client.
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

// UploadFile handles the file upload process to the server.
// It streams the file content from the client to the server.
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

// ApplyChanges sends a request to apply changes on the server for the given ID.
func (c *Grpc) ApplyChanges(ctx context.Context, id string) error {
	_, err := c.client.ApplyChanges(c.getCtx(ctx, c.token), &pb.ApplyChangesRequest{Guid: id})
	if err != nil {
		return fmt.Errorf("ApplyChanges: %w", err)
	}

	return nil
}

// IsNotAvailable checks if the server is not available.
func (c *Grpc) IsNotAvailable() bool {
	return !c.isAvailable
}

// IsAvailable checks if the server is available.
func (c *Grpc) IsAvailable() bool {
	return c.isAvailable
}

// TryToConnect attempts to establish a connection with the gRPC server.
// It sets up the connection and checks the server's availability.
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

// getCtx generates a new context with metadata containing the JWT for authorization.
func (c *Grpc) getCtx(ctx context.Context, jwt string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": jwt,
	})

	newCtx := metadata.NewOutgoingContext(ctx, md)

	return newCtx
}
