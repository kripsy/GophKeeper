package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/bufbuild/protovalidate-go"
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
	MiltipartUploadFile(context.Context, <-chan *pb.MiltipartUploadFileRequest, chan<- string, *string) error
}

func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
	s.logger.Debug("Start UploadFile")
	userID := 0
	deadlineDuration := DEADLINEDURATION

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timer := time.NewTimer(deadlineDuration)
	defer timer.Stop()

	doneChan := make(chan bool)
	errChan := make(chan error, 1)
	fileIdChan := make(chan string, 1)
	reqChan := make(chan *pb.MiltipartUploadFileRequest)
	var fileName string
	var once sync.Once
	defer func() {
		s.logger.Debug("defer in ")
		<-doneChan
		s.logger.Debug("got doneChan")
		close(errChan)
	}()

	go func() {
		err := s.secretUseCase.MiltipartUploadFile(ctx, reqChan, fileIdChan, &fileName)
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
			val, err := uuid.Parse(req.Guid)
			if err != nil {
				errChan <- err
				return
			}
			if isEbabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEbabled {
				errChan <- errors.New("You should block resource before use it")
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
				close(reqChan)
				s.logger.Debug("Got EOF of MiltipartUploadFile")
				select {
				case err := <-errChan:
					return err

				case fileName = <-fileIdChan:
					break loop
				}
			}

			cancel()
			return fmt.Errorf("%w", err)
		case req := <-reqChan:
			s.logger.Debug("update timer")
			timer.Reset(deadlineDuration)
			once.Do(func() {
				fileName = req.FileName
			})
			reqChan <- req
		}
	}

	return stream.SendAndClose(&pb.MiltipartUploadFileResponse{
		FileId: fileName,
	})
}

func (s *GrpcServer) BlockStore(stream pb.GophKeeperService_BlockStoreServer) error {

	s.logger.Debug("start Register")

	deadlineDuration := SYNCDEADLINEDURATION
	errChan := make(chan error, 1)
	reqChan := make(chan *pb.BlockStoreRequest)
	timer := time.NewTimer(deadlineDuration)

	var guid uuid.UUID
	var userID int
	userID = 0
	var once sync.Once
	var syncEnable bool
	defer func() {
		s.logger.Debug("Close main goroutine BlockStore")
		defer timer.Stop()
		if syncEnable {
			s.syncStatus.RemoveClientSync(userID, guid)
		}
	}()

	go func() {
		defer func() {
			s.logger.Debug("Close goroutine in BlockStore receivong data stream")
		}()
		v, err := protovalidate.New()
		if err != nil {
			errChan <- fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
			return
		}

		for {
			req, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			if err = v.Validate(req); err != nil {
				s.logger.Error("No valid req data", zap.Any("msg", req))
				errChan <- fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
				return
			}

			once.Do(func() {
				val, err := uuid.Parse(req.Guid)
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
				s.logger.Debug("Got EOF of BlockStore")
				break loop
			}
			s.logger.Error("Error in receive data", zap.Error(err))

			return status.Error(codes.Internal, err.Error())
		case req := <-reqChan:
			if !syncEnable {

				return status.Error(codes.ResourceExhausted, "Sync not enable")
			}
			fmt.Println(req.Guid)
			timer.Reset(deadlineDuration)
			s.logger.Debug("update timer for sync", zap.Int("userID", userID), zap.String("GUID", req.Guid), zap.Bool("sync is finish?", req.IsFinish))

		}
	}
	return stream.SendAndClose(&pb.BlockStoreResponse{Status: "All ok"})
}
