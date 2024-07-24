package database

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	type args struct {
		ctx         context.Context
		maxDuration time.Duration
		retryFunc   func() error
	}

	tests := map[string]struct {
		args    args
		wantErr bool
		errMsg  string
	}{
		"Successful retry": {
			args: args{
				ctx:         context.Background(),
				maxDuration: 5 * time.Second,
				retryFunc: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		"Retry until success": {
			args: args{
				ctx:         context.Background(),
				maxDuration: 5 * time.Second,
				retryFunc: func() error {
					if rand.Float64() < 0.8 {
						return errors.New("temporary error")
					}
					return nil
				},
			},
			wantErr: false,
		},
		"Retry exceeds maxDuration": {
			args: args{
				ctx:         context.Background(),
				maxDuration: 1 * time.Second,
				retryFunc: func() error {
					return errors.New("temporary error")
				},
			},
			wantErr: true,
			errMsg:  "temporary error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := Retry(tt.args.ctx, tt.args.maxDuration, tt.args.retryFunc)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRetryResult(t *testing.T) {
	type args struct {
		ctx         context.Context
		maxDuration time.Duration
		retryFunc   func() (any, error)
	}

	tests := map[string]struct {
		args    args
		wantErr bool
		errMsg  string
	}{
		"Successful retry result": {
			args: args{
				ctx:         context.Background(),
				maxDuration: 5 * time.Second,
				retryFunc: func() (any, error) {
					return "success", nil
				},
			},
			wantErr: false,
		},
		"Retry result until success": {
			args: args{
				ctx:         context.Background(),
				maxDuration: 5 * time.Second,
				retryFunc: func() (any, error) {
					if rand.Float64() < 0.8 {
						return nil, errors.New("temporary error")
					}
					return "success", nil
				},
			},
			wantErr: false,
		},
		"Retry result exceeds maxDuration": {
			args: args{
				ctx:         context.Background(),
				maxDuration: 1 * time.Second,
				retryFunc: func() (any, error) {
					return nil, errors.New("temporary error")
				},
			},
			wantErr: true,
			errMsg:  "temporary error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := RetryResult(tt.args.ctx, tt.args.maxDuration, tt.args.retryFunc)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
