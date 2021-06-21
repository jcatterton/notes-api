package external

import (
	"bytes"
	"context"
)

type ExtAPIHandler interface {
	ValidateToken(ctx context.Context, token string) error
	SendToContentService(ctx context.Context, body bytes.Buffer, contentType string) error
}
