package grpc

import (
	"context"
	"crypto/tls"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client interface {
	Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error)
	Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error)

	IsNotAvailable() bool
	IsAvailable() bool
	TryToConnect() bool
}

type Grpc struct {
	Client        pb.GophKeeperServiceClient
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

func (c *Grpc) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	return c.Client.Register(ctx, in)
}

func (c *Grpc) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	return c.Client.Login(ctx, in)
}

func (c *Grpc) IsNotAvailable() bool {
	return !c.isAvailable
}

func (c *Grpc) IsAvailable() bool {
	return c.isAvailable
}

func (c *Grpc) TryToConnect() bool {
	conn, err := grpc.Dial(c.serverAddress, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		c.log.Info().Err(err).Msg("failed connect to server")

		return false
	}

	//ping

	c.conn = conn
	c.Client = pb.NewGophKeeperServiceClient(conn)
	c.isAvailable = true

	return true
}
