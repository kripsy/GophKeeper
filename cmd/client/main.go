package main

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/gophkeeper/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	address := "127.0.0.1:50051"

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewShortenerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Hello(ctx, &pb.HelloRequest{Url: "your_url_here"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Response: %s", r.GetResult())
}
