package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DEADLINEDURATION is a maximum period between stream request batches
const DEADLINEDURATION = 15 * time.Second
const SYNCDEADLINEDURATION = 5 * time.Second

type SecretUseCase interface {
	MiltipartUploadFile(context.Context, <-chan []byte, chan<- string, *string) error
}

func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
	s.logger.Debug("Start UploadFile")
	deadlineDuration := DEADLINEDURATION

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
		s.logger.Debug("got doneChan")
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

func (s *GrpcServer) SyncClient(stream pb.GophKeeperService_SyncClientServer) error {
	s.logger.Debug("Start Sync client")

	deadlineDuration := SYNCDEADLINEDURATION
	errChan := make(chan error, 1)
	reqChan := make(chan *pb.SyncRequest)
	timer := time.NewTimer(deadlineDuration)

	var guid uuid.UUID
	var userID int
	userID = 0
	var once sync.Once
	var syncEnable bool
	defer func() {
		defer timer.Stop()
		if syncEnable {
			s.syncStatus.RemoveClientSync(userID, guid)
		}
	}()

	go func() {
		defer func() {
			s.logger.Debug("Close goroutine in SyncClient")
		}()
		for {
			req, err := stream.Recv()
			fmt.Println(syncEnable)

			once.Do(func() {
				val, err := uuid.Parse(req.GUID)
				if err != nil {
					syncEnable = false
				}
				guid = val
				syncEnable, _ = s.syncStatus.AddSync(userID, val)
				fmt.Println(syncEnable)
			})
			if !syncEnable {
				errChan <- errors.New("Sync for this user already exists")
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			reqChan <- req
		}
	}()
loop:
	for {
		select {
		case <-timer.C:
			return status.Error(codes.DeadlineExceeded, "Timeout reached")
		case err := <-errChan:
			if err == io.EOF {
				s.logger.Debug("Got EOF of SyncClient")
				select {
				case err := <-errChan:
					s.logger.Debug("We have error in SyncClient", zap.Error(err))

					return status.Error(codes.Internal, err.Error())
				default:
					s.logger.Debug("Finish sync, go out from main loop")
					break loop
				}
			}
			s.logger.Error("Error in receive data", zap.Error(err))

			return status.Error(codes.Internal, err.Error())
		case req := <-reqChan:
			if !syncEnable {

				return status.Error(codes.ResourceExhausted, "Sync not enable")
			}
			fmt.Println(req.GUID)
			timer.Reset(deadlineDuration)
			s.logger.Debug("update timer for sync", zap.Int("userID", userID), zap.String("GUID", req.GUID), zap.Bool("sync is finish?", req.IsFinish))

		}
	}
	return stream.SendAndClose(&pb.SyncResponse{Status: "All ok"})
}
