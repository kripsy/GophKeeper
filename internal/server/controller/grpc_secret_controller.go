package controller

import (
	"fmt"
	"io"
	"sync"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/utils"
)

type SecretUseCase interface {
	MiltipartUploadFile(<-chan []byte, chan<- string, *string) error
}

func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
	s.logger.Debug("Start UploadFile")
	dataChan := make(chan []byte)
	doneChan := make(chan bool)
	errChan := make(chan error, 1)
	fileIdChan := make(chan string, 1)
	var fileName string
	var once sync.Once
	defer func() {
		if !utils.IsClosed(dataChan) {
			close(dataChan)
		}
		s.logger.Debug("send exit to doneChan")
		<-doneChan
	}()

	go func() {
		err := s.secretUseCase.MiltipartUploadFile(dataChan, fileIdChan, &fileName)
		if err != nil {
			errChan <- err
		}
		doneChan <- true
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.logger.Debug("Got EOF of MiltipartUploadFile")
			select {
			case err := <-errChan:
				return err
			default:
				close(dataChan)
				fileId := <-fileIdChan
				return stream.SendAndClose(&pb.MiltipartUploadFileResponse{
					FileId: fileId,
				})
			}
		}
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		once.Do(func() {
			fileName = req.FileName
		})
		dataChan <- req.Content
	}
}
