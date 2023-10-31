package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kripsy/GophKeeper/cmd/testClient/client"
	"github.com/kripsy/GophKeeper/cmd/testClient/utils"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
)

func main() {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk5NjQ4OTIsIlVzZXJJRCI6MSwiVXNlcm5hbWUiOiJIZWxsbyJ9.RUt-iEhM1Dq4SEDwKt7rRxGG4T9q8SQVLchSkAqi1jU"
	ctx := utils.AddJwtIntoContext(context.Background(), jwt)
	filename := "test.txt"
	hash := "Hello"
	client, err := client.InitClient()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	defer client.Close()
	fmt.Println("all ok")

	var guid string
	guidChan := make(chan string, 1)
	errChan := make(chan error, 1)

	ctxBlockStream, cancelBlockStream := context.WithCancel(ctx)
	defer cancelBlockStream()

	// start in goroutine with context

	go func() {
		err = blockStream(ctxBlockStream, client, guidChan)
		if err != nil {
			errChan <- err
			return
		}
		fmt.Println("stop goroutine blockStream")
	}()
	select {
	case newGuid := <-guidChan:
		guid = newGuid
		fmt.Println("We got guid from blockStream without error ", guid)
		break
	case err := <-errChan:
		fmt.Println("We got error from blockStream ", err.Error())
	}

	fmt.Println("we got guid", guid)

	err = downloadFile(ctxBlockStream, client, filename, hash, guid)
	if err != nil {
		fmt.Println("We got error in downloadFile", err.Error())
		return
	}
	fmt.Println("No error =)")

	return
}

func blockStream(ctx context.Context, client *client.Client, guidChan chan string) error {
	guid := utils.GenerateUUID()
	defer close(guidChan)
	blockStoreStream, err := client.ClientGrpcService.BlockStore(ctx)
	if err != nil {
		return err
	}
	defer blockStoreStream.CloseSend()
	fmt.Println(guid)
	err = blockStoreStream.Send(&pb.BlockStoreRequest{Guid: guid})
	if err != nil {
		fmt.Println("We have error in blockStoreStream.Send(&pb.BlockStoreRequest{Guid: guid})", err.Error())
		return err
	}

	// Ожидаем ответ от сервера
	resp, err := blockStoreStream.Recv()
	if err != nil {
		fmt.Println("Error receiving response from server:", err.Error())
		return err
	}

	// Проверяем ответ от сервера (если у вас есть какие-либо критерии для проверки)
	if resp.Guid != guid {
		return fmt.Errorf("unexpected response from server: expected %s, got %s", guid, resp.Guid)
	}

	// Отправляем GUID в канал
	guidChan <- guid
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("blockStream ctx canceled")
			break loop
		}
	}
	fmt.Println("go out from loop in blockStream")
	return nil
}

func downloadFile(ctx context.Context, client *client.Client, fileName, hash, guid string) error {
	req := &pb.MultipartDownloadFileRequest{
		FileName: fileName,
		Guid:     guid,
		Hash:     hash,
	}
	fmt.Println(req)
	stream, err := client.ClientGrpcService.MultipartDownloadFile(ctx, req)
	if err != nil {
		fmt.Println("We have error in client.ClientGrpcService.MultipartDownloadFile(ctx, req)", err.Error())
		return fmt.Errorf("%w", err)
	}

	outFile, err := os.Create("tempFiles/" + fileName)
	if err != nil {
		fmt.Println("We have error in os.Create(fileName)", err.Error())
		return fmt.Errorf("%w", err)
	}

	defer outFile.Close()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("we got EOF")
			break
		}
		if err != nil {
			fmt.Println("We have error in stream.Recv()", err.Error())
			return fmt.Errorf("%w", err)
		}
		fmt.Println("new data")
		fmt.Println(string(resp.Content))

		_, writeErr := outFile.Write(resp.Content)
		if writeErr != nil {
			fmt.Println(writeErr.Error())
			return fmt.Errorf("%w", err)
		}
	}
	fmt.Println("File downloaded successfully!")
	return nil
}
