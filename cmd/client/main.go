package main

import (
	"fmt"
	"github.com/kripsy/GophKeeper/internal/client/app"
	"github.com/kripsy/GophKeeper/internal/client/config"
	//	"log"
	//pb "github.com/kripsy/GophKeeper/gen/pkg/api/gophkeeper/v1"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
)

func main() {
	cfg := config.GetConfig()

	a, err := app.NewApplication(cfg)
	a.PrepareApp()
	a.Run()

	//creds := credentials.NewTLS(&tls.Config{
	//	InsecureSkipVerify: true,
	//})

	//fm.AddToStorage("Note", filemanager.Note{Text: "text"}, models.DataInfo{DataType: 0, Description: "просто текст"})
	//fm.AddToStorage("Binance", filemanager.BasicAuth{Login: "L", Password: "P"}, models.DataInfo{DataType: 1, Description: "binance.com"})
	//fm.AddToStorage("Google", filemanager.BasicAuth{Login: "L", Password: "P"}, models.DataInfo{DataType: 1, Description: "google.com"})
	//fm.AddToStorage("SberBank", filemanager.CardData{Number: "11231231", Date: "02/24", CVV: "456"}, models.DataInfo{DataType: 2, Description: "visa"})

	//err = fm.UpdateDataAndInfoByName("note", Note{Text: "tetextxt"}, DataInfo{DataType: 1, Description: "текст"})
	//err = fm.DeleteByName("not")

	//conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println("did not connect: %v", err)
	}
	//defer conn.Close()

	//client := pb.NewShortenerServiceClient(conn)
	//
	//
	//
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//
	//r, err := client.Hello(ctx, &pb.HelloRequest{Url: "your_url_here"})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//
	//log.Printf("Response: %s", r.GetResult())

}
