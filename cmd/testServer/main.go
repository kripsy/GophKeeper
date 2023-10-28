package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
)

const BUCKETNAME = "secrets"
const ENDPOINT = "localhost:9000"
const ACCESSKEYID = "masoud"
const SECRETACCESSKEY = "Strong#Pass#2022"
const ISUSESSL = false

func initBucket(ctx context.Context, bucketName string, client *minio.Client) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func initMinioRepo(endpoint, accessKeyID, secretAccessKey string, isUseSSL bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: isUseSSL,
	})
	if err != nil {
		fmt.Printf("Error init client minio: %s", err.Error())

		return nil, fmt.Errorf("%w", err)
	}

	return minioClient, nil
}

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// client := pb.NewGophKeeperServiceClient(conn)

	// stream, err := client.BlockStore(context.Background())
	// if err != nil {
	// 	log.Fatalf("Error opening stream: %v", err)
	// }

	// go func(stream pb.GophKeeperService_BlockStoreClient) {
	// 	for {
	// 		select {
	// 		case <-stream.Context().Done():
	// 			fmt.Println("stream canceled")

	// 			return
	// 		default:
	// 			connState := conn.GetState()
	// 			if connState != connectivity.Ready {
	// 				fmt.Println("stream canceled or not ready")

	// 				return
	// 			}
	// 			fmt.Println("all ok")
	// 			fmt.Println(conn.GetState())
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	}
	// }(stream)

	go func() {
		for {
			fmt.Println(conn.GetState())
			time.Sleep(1 * time.Second)
		}
	}()

	// // Отправка нескольких BlockStoreRequest сообщений
	// for i := 0; i < 10000; i++ {
	// 	req := &pb.BlockStoreRequest{
	// 		Guid:     "guid-" + string(i),
	// 		IsFinish: false,
	// 	}
	// 	if err := stream.Send(req); err != nil {
	// 		log.Fatalf("Failed to send request: %v", err)
	// 	}
	// 	time.Sleep(2 * time.Second) // Задержка между отправкой сообщений
	// }

	// // Завершение стрима с is_finish = true
	// if err := stream.Send(&pb.BlockStoreRequest{Guid: "final-guid", IsFinish: true}); err != nil {
	// 	log.Fatalf("Failed to send final request: %v", err)
	// }

	// // Получение ответа
	// resp, err := stream.CloseAndRecv()
	// if err != nil {
	// 	log.Fatalf("Failed to receive response: %v", err)
	// }
	// log.Printf("Response: %s", resp.Status)

	// Для демонстрации: ожидание, чтобы увидеть, как "keepalive" работает
	time.Sleep(3000 * time.Second)
}

func main_minio() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	m, err := initMinioRepo(ENDPOINT, ACCESSKEYID, SECRETACCESSKEY, ISUSESSL)
	if err != nil {
		return
	}

	err = initBucket(ctx, BUCKETNAME, m)
	if err != nil {
		return
	}

	uuidFile, err := UploadFile(ctx, []byte{12, 31, 41, 55, 33, 13, 23, 55, 33}, "test File", m)
	if err != nil {
		return
	}

	downloadData, err := DownloadFile(ctx, uuidFile, m)
	if err != nil {
		return
	}

	fmt.Println(downloadData)
}

func UploadFile(ctx context.Context, data []byte, filename string, client *minio.Client) (uuid.UUID, error) {
	// generate uuid for create relation between filename and uuid
	// we should save secrets: GUID (UUID) user_id (int)
	// filename (string) hash256(data) ([]byte) lastUpdate (datetime)

	uuidWithHyphen := uuid.New()
	_, err := client.PutObject(ctx, BUCKETNAME, uuidWithHyphen.String(), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w", err)
	}

	return uuidWithHyphen, nil
}

func DownloadFile(ctx context.Context, minioFileID uuid.UUID, client *minio.Client) (*[]byte, error) {
	object, err := client.GetObject(ctx, BUCKETNAME, minioFileID.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &data, nil
}
