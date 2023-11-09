package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/kripsy/GophKeeper/internal/utils/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const secret = "testsecret"

func TestGetHash(t *testing.T) {
	// Create a logger for testing
	logger, _ := zap.NewProduction()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid Password",
			password: "testpassword",
			wantErr:  false,
		},
		{
			name:     "Empty Password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := auth.GetHash(context.Background(), tt.password, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHash() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !tt.wantErr {
				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.password))
				if err != nil {
					t.Errorf("Failed to compare hashed password with original password: %v", err)
				}
			}
		})
	}
}

func TestBuildJWTString(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		username  string
		secretKey string
		tokenExp  time.Duration
		wantErr   bool
	}{
		{
			name:      "Valid Input",
			userID:    1,
			username:  "testuser",
			secretKey: "testsecret",
			tokenExp:  time.Hour * 1,
			wantErr:   false,
		},
		{
			name:      "Empty Secret Key",
			userID:    1,
			username:  "testuser",
			secretKey: "",
			tokenExp:  time.Hour * 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.BuildJWTString(tt.userID, tt.username, tt.secretKey, tt.tokenExp)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}
		})
	}
}

func TestIsPasswordCorrect(t *testing.T) {
	logger, _ := zap.NewProduction()
	correctPasswordHash, _ := auth.GetHash(context.Background(), "testpassword", logger)
	tests := []struct {
		name         string
		password     string
		hashPassword string
		wantErr      bool
	}{
		{
			name:         "Valid Password",
			password:     "testpassword",
			hashPassword: correctPasswordHash,
			wantErr:      false,
		},
		{
			name:         "Invalid Password",
			password:     "wrongpassword",
			hashPassword: correctPasswordHash,
			wantErr:      true,
		},
		{
			name:         "Empty Password",
			password:     "",
			hashPassword: correctPasswordHash,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.IsPasswordCorrect(context.Background(), []byte(tt.password), []byte(tt.hashPassword), logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPasswordCorrect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		secret    string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "Valid Token",
			token:     returnToken(1, "testuser", secret, time.Hour*1),
			secret:    secret,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Invalid Secret",
			token:     returnToken(1, "testuser", secret, time.Hour*1),
			secret:    "wrongsecret",
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "Expired Token",
			token:     returnToken(1, "testuser", secret, -time.Hour*1),
			secret:    secret,
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := auth.IsValidToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidToken() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if valid != tt.wantValid {
				t.Errorf("IsValidToken() valid = %v, wantValid %v", valid, tt.wantValid)
			}
		})
	}
}

//nolint:dupl
func TestGetUsernameFromToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		secret   string
		wantUser string
		wantErr  bool
	}{
		{
			name:     "Valid Token",
			token:    returnToken(1, "testuser", secret, time.Hour*1),
			secret:   secret,
			wantUser: "testuser",
			wantErr:  false,
		},
		{
			name:     "Invalid Secret",
			token:    returnToken(1, "testuser", secret, time.Hour*1),
			secret:   "wrongsecret",
			wantUser: "",
			wantErr:  true,
		},
		{
			name:     "Expired Token",
			token:    returnToken(1, "testuser", secret, -time.Hour*1),
			secret:   secret,
			wantUser: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := auth.GetUsernameFromToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsernameFromToken() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if username != tt.wantUser {
				t.Errorf("GetUsernameFromToken() username = %v, wantUser %v", username, tt.wantUser)
			}
		})
	}
}

//nolint:dupl
func TestGetUseIDFromToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		secret  string
		wantID  int
		wantErr bool
	}{
		{
			name:    "Valid Token",
			token:   returnToken(123, "testuser", secret, time.Hour*1),
			secret:  secret,
			wantID:  123,
			wantErr: false,
		},
		{
			name:    "Invalid Secret",
			token:   returnToken(123, "testuser", secret, time.Hour*1),
			secret:  "wrongsecret",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Expired Token",
			token:   returnToken(123, "testuser", secret, -time.Hour*1),
			secret:  secret,
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := auth.GetUseIDFromToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUseIDFromToken() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if userID != tt.wantID {
				t.Errorf("GetUseIDFromToken() userID = %v, wantID %v", userID, tt.wantID)
			}
		})
	}
}

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		name     string
		password string
		salt     string
		wantErr  bool
	}{
		{
			name:     "Valid Password and Salt",
			password: "testpassword",
			salt:     "testsalt",
			wantErr:  false,
		},
		{
			name:     "Empty Password",
			password: "",
			salt:     "testsalt",
			wantErr:  false,
		},
		{
			name:     "Empty Salt",
			password: "testpassword",
			salt:     "",
			wantErr:  false,
		},
		{
			name:     "Empty Password and Salt",
			password: "",
			salt:     "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := auth.DeriveKey(tt.password, tt.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeriveKey() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !tt.wantErr && len(key) == 0 {
				t.Errorf("DeriveKey() key should not be nil or empty")
			}
			if !tt.wantErr && len(key) != 32 { // keyLen = 32
				t.Errorf("DeriveKey() key length = %v, want %v", len(key), 32)
			}
		})
	}
}

//nolint:dupl
func TestExtractTokenFromContext(t *testing.T) {
	tests := []struct {
		name string
		//nolint:containedctx
		ctx      context.Context
		expected string
		expectOk bool
	}{
		{
			name: "Valid Token",
			//nolint:staticcheck
			ctx:      context.WithValue(context.Background(), auth.TOKENCONTEXTKEY, "testtoken"),
			expected: "testtoken",
			expectOk: true,
		},
		{
			name:     "No Token",
			ctx:      context.Background(),
			expected: "",
			expectOk: false,
		},
		{
			name: "Invalid Type",
			//nolint:staticcheck
			ctx:      context.WithValue(context.Background(), auth.TOKENCONTEXTKEY, 12345),
			expected: "",
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, ok := auth.ExtractTokenFromContext(tt.ctx)
			if ok != tt.expectOk {
				t.Errorf("ExtractTokenFromContext() ok = %v, expectOk %v", ok, tt.expectOk)
			}
			if token != tt.expected {
				t.Errorf("ExtractTokenFromContext() token = %v, expected %v", token, tt.expected)
			}
		})
	}
}

//nolint:dupl
func TestExtractUsernameFromContext(t *testing.T) {
	tests := []struct {
		name string
		//nolint:containedctx
		ctx      context.Context
		expected string
		expectOk bool
	}{
		{
			name: "Valid Username",
			//nolint:staticcheck
			ctx:      context.WithValue(context.Background(), auth.USERNAMECONTEXTKEY, "testuser"),
			expected: "testuser",
			expectOk: true,
		},
		{
			name:     "No Username",
			ctx:      context.Background(),
			expected: "",
			expectOk: false,
		},
		{
			name: "Invalid Type",
			//nolint:staticcheck
			ctx:      context.WithValue(context.Background(), auth.USERNAMECONTEXTKEY, 12345),
			expected: "",
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, ok := auth.ExtractUsernameFromContext(tt.ctx)
			if ok != tt.expectOk {
				t.Errorf("ExtractUsernameFromContext() ok = %v, expectOk %v", ok, tt.expectOk)
			}
			if username != tt.expected {
				t.Errorf("ExtractUsernameFromContext() username = %v, expected %v", username, tt.expected)
			}
		})
	}
}

//nolint:dupl
func TestExtractUserIDFromContext(t *testing.T) {
	tests := []struct {
		name string
		//nolint:containedctx
		ctx      context.Context
		expected int
		expectOk bool
	}{
		{
			name: "Valid UserID",
			//nolint:staticcheck
			ctx:      context.WithValue(context.Background(), auth.USERIDCONTEXTKEY, 12345),
			expected: 12345,
			expectOk: true,
		},
		{
			name:     "No UserID",
			ctx:      context.Background(),
			expected: 0,
			expectOk: false,
		},
		{
			name: "Invalid Type",
			//nolint:staticcheck
			ctx:      context.WithValue(context.Background(), auth.USERIDCONTEXTKEY, "testuser"),
			expected: 0,
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, ok := auth.ExtractUserIDFromContext(tt.ctx)
			if ok != tt.expectOk {
				t.Errorf("ExtractUserIDFromContext() ok = %v, expectOk %v", ok, tt.expectOk)
			}
			if userID != tt.expected {
				t.Errorf("ExtractUserIDFromContext() userID = %v, expected %v", userID, tt.expected)
			}
		})
	}
}

//nolint:unparam
func returnToken(userID int, username string, secret string, duration time.Duration) string {
	token, _ := auth.BuildJWTString(userID, username, secret, duration)

	return token
}
