package external

import (
	"bytes"
	"context"
)

type ExtAPIHandler interface {
	SetToken(token string)
	ValidateToken(ctx context.Context, token string) error
	SendToContentService(ctx context.Context, body bytes.Buffer, contentType string) error
}
