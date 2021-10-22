package external

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Requester interface {
	Do(r *http.Request) (*http.Response, error)
}

type ExtAPI struct {
	Client            Requester
	LoginServiceURL   string
	ContentServiceURL string
	Token             string
}

func (ext *ExtAPI) SetToken(token string) {
	ext.Token = token
}

func (ext *ExtAPI) ValidateToken(ctx context.Context, token string) error {
	if ext.LoginServiceURL == "" {
		return errors.New("login service url cannot be empty")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/token", ext.LoginServiceURL), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	resp, err := ext.Client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("non-200 status code received: %v", resp.StatusCode))
	}

	return nil
}

func (ext *ExtAPI) SendToContentService(ctx context.Context, body bytes.Buffer, contentType string) error {
	if ext.ContentServiceURL == "" {
		return errors.New("content service url cannot be empty")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/upload", ext.ContentServiceURL), &body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ext.Token))
	req.Header.Add("Content-Type", contentType)

	res, err := ext.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status code received: %v", res.StatusCode)
	}

	return nil
}
