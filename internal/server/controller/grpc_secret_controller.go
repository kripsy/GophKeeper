package controller

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SecretUseCase interface {
	MiltipartUploadFile(context.Context, <-chan []byte, chan<- string, *string) error
}

func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
	s.logger.Debug("Start UploadFile")
	deadlineDuration := time.Second * 4

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timer := time.NewTimer(deadlineDuration)
	defer timer.Stop()

	dataChan := make(chan []byte)
	doneChan := make(chan bool)
	errChan := make(chan error, 1)
	fileIdChan := make(chan string, 1)
	reqChan := make(chan *pb.MiltipartUploadFileRequest)
	var fileName string
	var once sync.Once
	defer func() {
		s.logger.Debug("send exit to doneChan")
		<-doneChan
	}()

	go func() {
		err := s.secretUseCase.MiltipartUploadFile(ctx, dataChan, fileIdChan, &fileName)
		s.logger.Debug("s.secretUseCase.MiltipartUploadFile ended", zap.Any("error", err))
		if err != nil {
			errChan <- err
		}
		doneChan <- true
	}()

	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			reqChan <- req
		}
	}()
loop:
	for {
		s.logger.Debug("it's loop")
		select {
		case <-timer.C:
			s.logger.Debug("get end context in handler")
			cancel()
			return status.Error(codes.DeadlineExceeded, "No data from timeout")
		case err := <-errChan:
			if err == io.EOF {
				close(dataChan)
				s.logger.Debug("Got EOF of MiltipartUploadFile")
				select {
				case err := <-errChan:
					return err

				case fileName = <-fileIdChan:
					break loop
				}
			}
			return fmt.Errorf("%w", err)
		case req := <-reqChan:
			s.logger.Debug("update timer")
			timer.Reset(deadlineDuration)
			once.Do(func() {
				fileName = req.FileName
			})
			dataChan <- req.Content
		}
	}

	return stream.SendAndClose(&pb.MiltipartUploadFileResponse{
		FileId: fileName,
	})
}
