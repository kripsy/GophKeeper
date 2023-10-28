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
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SecretUseCase interface {
	MiltipartUploadFile(context.Context, <-chan *models.MiltipartUploadFileData, string) (bool, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
}

func (s *GrpcServer) MiltipartUploadFile(stream pb.GophKeeperService_MiltipartUploadFileServer) error {
	s.logger.Debug("Start MiltipartUploadFile")

	ctx := stream.Context()
	userID, ok := utils.ExtractUserIDFromContext(ctx)
	if !ok {
		s.logger.Error("cannot get userID from context")

		return status.Errorf(codes.Internal, "Failed to extract userID")
	}

	username, ok := utils.ExtractUsernameFromContext(ctx)
	if !ok {
		s.logger.Error("cannot get username from context")

		return status.Errorf(codes.Internal, "Failed to extract username")
	}

	bucketName, err := utils.FromUser2BucketName(ctx, username, userID)
	if err != nil {
		s.logger.Error("cannot get bucketName")

		return status.Errorf(codes.Internal, "Failed to extract bucketName")
	}
	s.logger.Debug("bucket name", zap.String("msg", bucketName))

	var once sync.Once
	s.logger.Debug("", zap.Any("FUCK", stream.Context()))

	var fileID string

	reqChan := make(chan *models.MiltipartUploadFileData, 1)

	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(errChan)
		defer close(reqChan)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				s.logger.Error("Error get data from stream", zap.Error(err))
				errChan <- err
				return
			}
			val, err := uuid.Parse(req.Guid)
			if err != nil {
				s.logger.Error("Couldn't parse GUID", zap.Error(err))

				errChan <- err
				return
			}
			if isEnabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEnabled {
				s.logger.Error("You should block resource before use it", zap.Error(err))
				errChan <- errors.New("You should block resource before use it")
				return
			}
			once.Do(func() {
				fileID = req.GetFileName()

			})

			reqChan <- &models.MiltipartUploadFileData{
				Content:  req.GetContent(),
				FileName: req.GetFileName(),
				Guid:     req.GetGuid(),
				Hash:     req.GetHash(),
				Username: username,
			}
		}
	}()

	success, err := s.secretUseCase.MiltipartUploadFile(stream.Context(), reqChan, bucketName)
	if err != nil {
		s.logger.Error("Error in s.secretUseCase.MiltipartUploadFile", zap.Error(err))

		return err
	}
	if !success {
		s.logger.Error("failed upload file", zap.Error(err))

		return status.Errorf(codes.Internal, "Failed to upload file")
	}
	s.logger.Debug("end upload")

	return stream.SendAndClose(&pb.MiltipartUploadFileResponse{
		FileId: fileID,
	})
}

func (s *GrpcServer) BlockStore(stream pb.GophKeeperService_BlockStoreServer) error {

	s.logger.Debug("start Register")
	streamCtx := stream.Context()
	errChan := make(chan error, 1)
	reqChan := make(chan *pb.BlockStoreRequest)

	var guid uuid.UUID
	ctx := stream.Context()
	userID, ok := utils.ExtractUserIDFromContext(ctx)
	if !ok {
		s.logger.Error("cannot get userID from context")

		return status.Errorf(codes.Internal, "Failed to extract userID")
	}

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
				s.logger.Debug("Sync", zap.Bool("msg", syncEnable))
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
