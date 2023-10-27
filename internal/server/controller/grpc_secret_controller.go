package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SecretUseCase interface {
	MiltipartUploadFile(context.Context, <-chan *pb.MiltipartUploadFileRequest) (bool, error)
}

func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
	s.logger.Debug("Start MiltipartUploadFile")
	userID := 0

	reqChan := make(chan *pb.MiltipartUploadFileRequest, 1)

	errChan := make(chan error, 1)

	go func() {
		defer close(reqChan)
		defer close(errChan)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			val, err := uuid.Parse(req.Guid)
			if err != nil {
				errChan <- err
				return
			}
			if isEnabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEnabled {
				errChan <- errors.New("You should block resource before use it")
				return
			}
			reqChan <- req
		}
	}()

	success, err := s.secretUseCase.MiltipartUploadFile(stream.Context(), reqChan)
	if err != nil {
		s.logger.Error("Error in s.secretUseCase.MiltipartUploadFile", zap.Error(err))

		return err
	}

	select {
	case err := <-errChan:
		return err
	default:
	}

	if !success {
		return status.Errorf(codes.Internal, "Failed to upload file")
	}

	return stream.SendAndClose(&pb.MiltipartUploadFileResponse{
		FileId: "123",
	})
}

// func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
// 	s.logger.Debug("Start UploadFile")
// 	userID := 0
// 	streamCtx := stream.Context()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	doneChan := make(chan bool)
// 	errChan := make(chan error, 1)
// 	reqChan := make(chan *pb.MiltipartUploadFileRequest)
// 	var resultUsecase bool

// 	defer func() {
// 		close(errChan)
// 		close(doneChan)

// 	}()

// 	go func() {
// 		res, err := s.secretUseCase.MiltipartUploadFile(ctx, reqChan)
// 		resultUsecase = res
// 		s.logger.Debug("s.secretUseCase.MiltipartUploadFile ended", zap.Any("error", err))
// 		if err != nil {
// 			errChan <- err

// 			return
// 		}
// 		s.logger.Debug("Write to done chan")
// 		doneChan <- true
// 		s.logger.Debug("Wrote to done chan")
// 	}()

// 	go func() {
// 		for {
// 			req, err := stream.Recv()
// 			if err != nil {
// 				errChan <- err
// 				return
// 			}
// 			val, err := uuid.Parse(req.Guid)
// 			if err != nil {
// 				errChan <- err
// 				return
// 			}
// 			if isEbabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEbabled {
// 				errChan <- errors.New("You should block resource before use it")
// 				return
// 			}
// 			reqChan <- req
// 		}
// 	}()
// loop:
// 	for {
// 		s.logger.Debug("it's loop")
// 		select {
// 		case <-streamCtx.Done():
// 			s.logger.Debug("get end context in handler")
// 			cancel()

// 			return status.Error(codes.DeadlineExceeded, "No data from timeout")
// 		case err := <-errChan:
// 			if err == io.EOF {

// 				s.logger.Debug("Got EOF of MiltipartUploadFile")
// 				close(reqChan)
// 				s.logger.Debug("Waiting done chan")
// 				<-doneChan
// 				s.logger.Debug("Got done chan")
// 				s.logger.Debug("Result usecase", zap.Bool("Succes?", resultUsecase))
// 				break loop
// 			}
// 			cancel()

// 			return fmt.Errorf("%w", err)
// 		case req := <-reqChan:
// 			s.logger.Debug("got new stream data")
// 			reqChan <- req
// 		}
// 	}

// 	return stream.SendAndClose(&pb.MiltipartUploadFileResponse{
// 		FileId: "123",
// 	})
// }

func (s *GrpcServer) BlockStore(stream pb.GophKeeperService_BlockStoreServer) error {

	s.logger.Debug("start Register")
	streamCtx := stream.Context()
	errChan := make(chan error, 1)
	reqChan := make(chan *pb.BlockStoreRequest)

	var guid uuid.UUID
	var userID int
	userID = 0
	var once sync.Once
	var syncEnable bool
	defer func() {
		s.logger.Debug("Close main goroutine BlockStore")
		if syncEnable {
			s.syncStatus.RemoveClientSync(userID, guid)
		}
	}()

	go func() {
		defer func() {
			s.logger.Debug("Close goroutine in BlockStore receivong data stream")
			s.logger.Debug("Close errChan")
			close(errChan)
			close(reqChan)
		}()
		v, err := protovalidate.New()
		if err != nil {
			s.logger.Error("Error init validator", zap.Error(err))
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
				s.logger.Error("Sync", zap.Bool("msg", syncEnable))
			})
			if !syncEnable {
				s.logger.Error("Sync not enabled")
				errChan <- errors.New("Sync for this user already exists")

				return
			}

			reqChan <- req
		}
	}()
loop:
	for {
		select {
		case <-streamCtx.Done():
			s.logger.Debug("stream context canceled")

			return status.Error(codes.Aborted, "stream context canceled")
		case err := <-errChan:
			if err == io.EOF {
				s.logger.Debug("Got EOF of BlockStore")
				break loop
			}
			s.logger.Error("Error in receive data", zap.Error(err))

			return status.Error(codes.Internal, err.Error())
		case req := <-reqChan:
			if !syncEnable {
				s.logger.Error("Sync not enabled")

				return status.Error(codes.ResourceExhausted, "Sync not enable")
			}

			s.logger.Debug("update timer for sync", zap.Int("userID", userID), zap.String("GUID", req.Guid), zap.Bool("sync is finish?", req.IsFinish))

		}
	}
	return stream.SendAndClose(&pb.BlockStoreResponse{Status: "All ok"})
}
