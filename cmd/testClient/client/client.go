package client

import (
	"fmt"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"google.golang.org/grpc"
)

type Client struct {
	ClientGrpcService pb.GophKeeperServiceClient
	ClientConn        *grpc.ClientConn
}

func InitClient() (*Client, error) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	client := pb.NewGophKeeperServiceClient(conn)

	return &Client{
		ClientGrpcService: client,
		ClientConn:        conn,
	}, nil
}

func (c *Client) Close() error {
	err := c.ClientConn.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
