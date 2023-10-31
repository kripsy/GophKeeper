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
	MultipartUploadFile(context.Context, <-chan *models.MultipartUploadFileData, string) (bool, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
	MultipartDownloadFile(context.Context,
		*models.MultipartDownloadFileRequest,
		string) (chan *models.MultipartDownloadFileResponse, chan error)
	ApplyChanges(ctx context.Context, bucketName string) (bool, error)
	DiscardChanges(ctx context.Context, bucketName string) (bool, error)
}

func (s *GrpcServer) MultipartUploadFile(stream pb.GophKeeperService_MultipartUploadFileServer) error {
	s.logger.Debug("Start MultipartUploadFile")

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
	var fileID string

	reqChan := make(chan *models.MultipartUploadFileData, 1)
	errChan := make(chan error, 1)
	doneChan := make(chan bool)

	go func() {
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

			reqChan <- &models.MultipartUploadFileData{
				Content:  req.GetContent(),
				FileName: req.GetFileName(),
				Guid:     req.GetGuid(),
				Hash:     req.GetHash(),
				Username: username,
			}
		}
	}()

	go func() {
		success, err := s.secretUseCase.MultipartUploadFile(stream.Context(), reqChan, bucketName)
		if err != nil {
			s.logger.Error("Error in s.secretUseCase.MultipartUploadFile", zap.Error(err))
			errChan <- err
			return
		}
		if !success {
			s.logger.Error("failed upload file")
			errChan <- errors.New("Failed to upload file")
			return
		}
		doneChan <- true
	}()

	select {
	case <-doneChan:
		s.logger.Debug("end upload")

	case err := <-errChan:
		if err != nil {
			s.logger.Debug("was some error", zap.Any("msg", err))
			return status.Error(codes.Internal, err.Error())
		}
	}
	return stream.SendAndClose(&pb.MultipartUploadFileResponse{
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
	// clear temp files
	defer func(ctx context.Context) {
		userID, ok := utils.ExtractUserIDFromContext(ctx)
		if !ok {
			s.logger.Error("cannot get userID from context")
			return
		}

		username, ok := utils.ExtractUsernameFromContext(ctx)
		if !ok {
			s.logger.Error("cannot get username from context")
			return
		}

		bucketName, err := utils.FromUser2BucketName(ctx, username, userID)
		if err != nil {
			s.logger.Error("cannot get bucketName")
			return
		}
		s.logger.Debug("bucket name", zap.String("msg", bucketName))
		s.logger.Debug("Start discard changes")
		newCtx := context.Background()
		success, err := s.secretUseCase.DiscardChanges(newCtx, bucketName)
		if err != nil {
			s.logger.Error("Error discard changes", zap.Error(err))

			return
		}
		if !success {
			s.logger.Debug("Fail discard changes", zap.Error(err))

			return
		}
		return
	}(ctx)

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
			stream.Send(&pb.BlockStoreResponse{
				Guid: req.Guid,
			})
		}
	}
	return nil
}

func (s *GrpcServer) MultipartDownloadFile(req *pb.MultipartDownloadFileRequest, stream pb.GophKeeperService_MultipartDownloadFileServer) error {
	s.logger.Debug("Start MultipartDownloadFile")

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

	val, err := uuid.Parse(req.Guid)
	if err != nil {
		s.logger.Error("Couldn't parse GUID", zap.Error(err))

		return status.Errorf(codes.Internal, "Couldn't parse GUID")

	}
	if isEnabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEnabled {
		s.logger.Error("You should block resource before use it", zap.Error(err))

		return status.Errorf(codes.Internal, "You should block resource before use it")
	}

	multipartDownloadFileRequest := &models.MultipartDownloadFileRequest{
		FileName: req.FileName,
		Guid:     req.Guid,
		Hash:     req.Hash,
	}
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	dataChan, errChan := s.secretUseCase.MultipartDownloadFile(newCtx, multipartDownloadFileRequest, bucketName)
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				s.logger.Debug("dataChan reading false")
				dataChan = nil
				continue
			}
			response := &pb.MultipartDownloadFileResponse{
				FileName: data.FileName,
				Guid:     data.Guid,
				Content:  data.Content,
				Hash:     data.Hash,
			}
			if err := stream.Send(response); err != nil {
				s.logger.Error("Error send response", zap.Error(err))

				return err
			}
		case err := <-errChan:
			if err == io.EOF {
				s.logger.Debug("We got EOF")
				errChan = nil
				continue
			} else {
				s.logger.Error("Error in secretUseCase.MultipartDownloadFile", zap.Error(err))
				return err
			}
		case <-ctx.Done():
			s.logger.Debug("stream context canceled")

			return status.Error(codes.Aborted, "stream context canceled")

		}
		if dataChan == nil && errChan == nil {
			s.logger.Debug("All channel is nil")
			break
		}
	}

	s.logger.Debug("end download")

	return nil
}

func (s *GrpcServer) ApplyChanges(ctx context.Context,
	req *pb.ApplyChangesRequest) (*pb.ApplyChangesResponse, error) {
	s.logger.Debug("start ApplyChanges")
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
	}

	userID, ok := utils.ExtractUserIDFromContext(ctx)
	if !ok {
		s.logger.Error("cannot get userID from context")
		return nil, status.Errorf(codes.Internal, "Failed to extract userID")
	}

	username, ok := utils.ExtractUsernameFromContext(ctx)
	if !ok {
		s.logger.Error("cannot get username from context")
		return nil, status.Errorf(codes.Internal, "Failed to extract username")
	}

	bucketName, err := utils.FromUser2BucketName(ctx, username, userID)
	if err != nil {
		s.logger.Error("cannot get bucketName")
		return nil, status.Errorf(codes.Internal, "Failed to extract bucketName")
	}
	s.logger.Debug("bucket name", zap.String("msg", bucketName))

	val, err := uuid.Parse(req.Guid)
	if err != nil {
		s.logger.Error("Couldn't parse GUID", zap.Error(err))

		return nil, status.Errorf(codes.Internal, "Couldn't parse GUID")

	}
	if isEnabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEnabled {
		s.logger.Error("You should block resource before use it", zap.Error(err))

		return nil, status.Errorf(codes.Internal, "You should block resource before use it")
	}
	s.logger.Debug("Start apply changes")
	ctx, cancel := context.WithCancel(ctx)
	doneChan := make(chan interface{})
	go func() {
		for {
			select {
			case <-doneChan:
				s.logger.Debug("main func done, closing goruitine")
				return
			default:
				if isEnabled, _ := s.syncStatus.IsSyncExists(userID, val); !isEnabled {
					s.logger.Error("You should block resource before use it", zap.Error(err))
					cancel()
					return
				}
			}
		}
	}()
	success, err := s.secretUseCase.ApplyChanges(ctx, bucketName)
	if err != nil {
		s.logger.Error("Error s.secretUseCase.ApplyChanges", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if !success {
		s.logger.Debug("Failed s.secretUseCase.ApplyChanges", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "Failed to apply changes"))
	}

	return &pb.ApplyChangesResponse{
		Guid: req.Guid,
	}, nil
}
