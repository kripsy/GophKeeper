//nolint:testpackage
package usecase

import (
	"testing"
)

func TestClientUsecase_about(t *testing.T) {
	tests := []struct {
		name     string
		aboutMsg string
	}{
		{
			name:     "ok=)",
			aboutMsg: "I cover what I can",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				aboutMsg: tt.aboutMsg,
			}
			c.about()
		})
	}
}
