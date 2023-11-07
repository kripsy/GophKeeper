package controller_test

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/kripsy/GophKeeper/internal/server/controller/mocks"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcServerMultipartUploadFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStream := mocks.NewMockGophKeeperService_MultipartUploadFileServer(mockCtrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(mockCtrl)
	mockSyncStatus := mocks.NewMockSyncStatus(mockCtrl)

	logger := zap.NewNop()

	grpcServer := controller.InitGrpcServiceServer(
		nil,
		mockSecretUseCase,
		"secret",
		logger,
		mockSyncStatus,
	)

	testCases := []struct {
		name           string
		setupMocks     func()
		expectedError  error
		expectedFileID string
	}{
		{
			name: "Success",
			setupMocks: func() {
				//nolint:staticcheck
				newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
				//nolint:staticcheck
				newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
				mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
				mockSyncStatus.EXPECT().IsSyncExists(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
				mockStream.EXPECT().Recv().Return(&pb.MultipartUploadFileRequest{}, io.EOF)
				mockSecretUseCase.EXPECT().MultipartUploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockStream.EXPECT().SendAndClose(gomock.Any()).Return(nil)
			},
			expectedError:  nil,
			expectedFileID: "test-file-id",
		},
		{
			name: "Error upload in usecase",
			setupMocks: func() {
				//nolint:staticcheck
				newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
				//nolint:staticcheck
				newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
				mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
				mockStream.EXPECT().Recv().Return(nil, models.NewUnionError("stream receive error")).AnyTimes()
				mockSecretUseCase.EXPECT().MultipartUploadFile(gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(false, models.NewUnionError("")).AnyTimes()
			},
			expectedError: status.Error(codes.Internal, ""),
		},
		// {
		// 	name: "Error recv data",
		// 	setupMocks: func() {
		// 		newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
		// 		newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
		// 		mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
		// 		mockStream.EXPECT().Recv().Return(nil, errors.New("stream receive error")).AnyTimes()
		// 		mockSecretUseCase.EXPECT().MultipartUploadFile(gomock.Any(),
		// gomock.Any(),
		// gomock.Any()).Return(false, errors.New("")).AnyTimes()
		// 	},
		// 	expectedError: status.Error(codes.Internal, ""),
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := grpcServer.MultipartUploadFile(mockStream)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tc.expectedError), status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBlockStore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStream := mocks.NewMockGophKeeperService_BlockStoreServer(mockCtrl)
	logger := zap.NewNop()
	mockSyncStatus := mocks.NewMockSyncStatus(mockCtrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(mockCtrl)

	grpcServer := controller.InitGrpcServiceServer(
		nil,
		mockSecretUseCase,
		"secret",
		logger,
		mockSyncStatus,
	)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		//nolint:containedctx
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "Success",
			setup: func() {
				guid := uuid.New().String()
				//nolint:staticcheck
				newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
				//nolint:staticcheck
				newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
				bucketName, _ := utils.FromUser2BucketName(newCtx, "user", 1)
				mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
				gomock.InOrder(
					mockStream.EXPECT().Recv().Return(&pb.BlockStoreRequest{Guid: guid}, nil).Times(1),
					mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1), // Симулируем окончание стрима
				)
				mockStream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
				mockSyncStatus.EXPECT().AddSync(1, gomock.Any()).Return(true, nil).AnyTimes()
				mockSecretUseCase.EXPECT().DiscardChanges(gomock.Any(), bucketName).Return(true, nil).AnyTimes()
				mockSyncStatus.EXPECT().RemoveClientSync(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "RecvError",
			setup: func() {
				//nolint:staticcheck
				newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
				//nolint:staticcheck
				newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
				bucketName, _ := utils.FromUser2BucketName(newCtx, "user", 1)
				mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
				mockStream.EXPECT().Recv().Return(nil, models.NewUnionError("error")).Times(1)
				mockSecretUseCase.EXPECT().DiscardChanges(gomock.Any(), bucketName).Return(true, nil).AnyTimes()
				mockSyncStatus.EXPECT().RemoveClientSync(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := grpcServer.BlockStore(mockStream)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplyChanges(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSyncStatus := mocks.NewMockSyncStatus(mockCtrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(mockCtrl)
	logger := zap.NewNop()

	grpcServer := controller.InitGrpcServiceServer(
		nil,
		mockSecretUseCase,
		"secret",
		logger,
		mockSyncStatus,
	)

	tests := []struct {
		name    string
		req     *pb.ApplyChangesRequest
		setup   func(ctx context.Context)
		wantErr bool
		//nolint:containedctx
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "Success",
			req: &pb.ApplyChangesRequest{
				Guid: uuid.New().String(),
			},
			setup: func(ctx context.Context) {
				userID := 1
				username := "user"
				bucketName, _ := utils.FromUser2BucketName(ctx, username, userID)
				mockSyncStatus.EXPECT().IsSyncExists(userID, gomock.Any()).Return(true, nil).AnyTimes()
				mockSecretUseCase.EXPECT().ApplyChanges(gomock.Any(), bucketName).Return(true, nil).AnyTimes()
			},
			wantErr: false,
			ctx: context.WithValue(context.WithValue(context.Background(),
				//nolint:staticcheck
				utils.USERNAMECONTEXTKEY, "user"), utils.USERIDCONTEXTKEY, 1),
		},
		{
			name: "ValidationError",
			req: &pb.ApplyChangesRequest{
				Guid: "",
			},
			setup: func(ctx context.Context) {

			},
			wantErr: true,
			ctx: context.WithValue(context.WithValue(context.Background(),
				//nolint:staticcheck
				utils.USERNAMECONTEXTKEY, "user"), utils.USERIDCONTEXTKEY, 1),
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.ctx)

			resp, err := grpcServer.ApplyChanges(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Guid, resp.Guid)
			}
		})
	}
}

func TestMultipartDownloadFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStream := mocks.NewMockGophKeeperService_MultipartDownloadFileServer(mockCtrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(mockCtrl)
	mockSyncStatus := mocks.NewMockSyncStatus(mockCtrl)
	logger := zap.NewNop()

	grpcServer := controller.InitGrpcServiceServer(
		nil,
		mockSecretUseCase,
		"secret",
		logger,
		mockSyncStatus,
	)

	tests := []struct {
		name    string
		req     *pb.MultipartDownloadFileRequest
		setup   func(ctx context.Context)
		wantErr bool
		//nolint:containedctx
		ctx     context.Context
		errCode codes.Code
	}{
		// {
		// 	name: "Success",
		// 	req: &pb.MultipartDownloadFileRequest{
		// 		FileName: "testfile.txt",
		// 		Guid:     uuid.New().String(),
		// 		Hash:     "somehash",
		// 	},
		// 	setup: func(ctx context.Context) {
		// 		userID := 1
		// 		username := "user"
		// 		bucketName, _ := utils.FromUser2BucketName(ctx, username, userID)
		// 		dataChan := make(chan *models.MultipartDownloadFileResponse, 1)
		// 		errChan := make(chan error, 1)

		// 		dataChan <- &models.MultipartDownloadFileResponse{
		// 			FileName: "testfile.txt",
		// 			Guid:     uuid.New().String(),
		// 			Content:  []byte("file content"),
		// 			Hash:     "somehash",
		// 		}
		// 		close(dataChan)
		// 		newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
		// 		newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
		// 		bucketName, _ = utils.FromUser2BucketName(newCtx, "user", 1)
		// 		mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
		// 		mockSyncStatus.EXPECT().IsSyncExists(userID, gomock.Any()).Return(true, nil).AnyTimes()
		// 		mockStream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
		// 		mockSecretUseCase.EXPECT().MultipartDownloadFile(gomock.Any(),
		// gomock.Any(), bucketName).Return(dataChan, errChan).AnyTimes()
		// 	},
		// 	wantErr: false,
		// 	ctx:     context.WithValue(context.WithValue(context.Background(),
		// utils.USERNAMECONTEXTKEY, "user"), utils.USERIDCONTEXTKEY, 1),
		// },
		{
			name: "Error",
			req: &pb.MultipartDownloadFileRequest{
				FileName: "testfile.txt",
				Guid:     uuid.New().String(),
				Hash:     "somehash",
			},
			setup: func(ctx context.Context) {
				userID := 1
				username := "user"
				bucketName, _ := utils.FromUser2BucketName(ctx, username, userID)
				dataChan := make(chan *models.MultipartDownloadFileResponse)
				errChan := make(chan error, 1)

				errChan <- models.NewUnionError("test error")
				close(errChan)
				//nolint:staticcheck
				newCtx := context.WithValue(context.Background(), utils.USERNAMECONTEXTKEY, "user")
				//nolint:staticcheck
				newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, 1)
				mockStream.EXPECT().Context().Return(newCtx).AnyTimes()
				mockSyncStatus.EXPECT().IsSyncExists(userID,
					gomock.Any()).Return(true, nil).AnyTimes()
				mockSecretUseCase.EXPECT().MultipartDownloadFile(gomock.Any(),
					gomock.Any(),
					bucketName).Return(dataChan, errChan).AnyTimes()
			},
			wantErr: true,
			ctx: context.WithValue(context.WithValue(context.Background(),
				//nolint:staticcheck
				utils.USERNAMECONTEXTKEY, "user"), utils.USERIDCONTEXTKEY, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.ctx)

			err := grpcServer.MultipartDownloadFile(tt.req, mockStream)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
