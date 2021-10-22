package external

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"notes-api/pkg/testhelper/mocks"
	"testing"
)

func TestExternal_SetToken_ShouldSetToken(t *testing.T) {
	ext := ExtAPI{
		Token: "",
	}

	ext.SetToken("test")

	require.Equal(t, "test", ext.Token)
}

func TestExternal_ValidateToken_ShouldReturnErrorIfLoginServiceURLIsBlank(t *testing.T) {
	ext := ExtAPI{
		LoginServiceURL: "",
	}

	err := ext.ValidateToken(context.TODO(), "test")
	require.NotNil(t, err)
	require.Equal(t, "login service url cannot be empty", err.Error())
}

func TestExternal_ValidateToken_ShouldReturnErrorIfErrorOccursPerformingRequest(t *testing.T) {
	mockRequester := &mocks.Requester{}
	mockRequester.On("Do", mock.Anything).Return(nil, errors.New("test"))

	ext := ExtAPI{
		LoginServiceURL: "test",
		Client: mockRequester,
	}

	err := ext.ValidateToken(context.TODO(), "test")
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestExternal_ValidateToken_ShouldReturnErrorIfResponseStatusCodeIsNot200(t *testing.T) {
	mockRequester := &mocks.Requester{}
	mockRequester.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusTeapot}, nil)

	ext := ExtAPI{
		LoginServiceURL: "test",
		Client: mockRequester,
	}

	err := ext.ValidateToken(context.TODO(), "test")
	require.NotNil(t, err)
	require.Equal(t, "non-200 status code received: 418", err.Error())
}

func TestExternal_ValidateToken_ShouldReturnNoErrorIfResponseIs200(t *testing.T) {
	mockRequester := &mocks.Requester{}
	mockRequester.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	ext := ExtAPI{
		LoginServiceURL: "test",
		Client: mockRequester,
	}

	require.Nil(t, ext.ValidateToken(context.TODO(), "test"))
}

func TestExternal_SendToContentService_ShouldReturnErrorIfContentServiceURLIsBlank(t *testing.T) {
	ext := ExtAPI{
		ContentServiceURL: "",
	}

	err := ext.SendToContentService(context.TODO(), *bytes.NewBuffer(nil), "")
	require.NotNil(t, err)
	require.Equal(t, "content service url cannot be empty", err.Error())
}

func TestExternal_SendToContentService_ShouldReturnErrorIfErrorOccursPerformingRequest(t *testing.T) {
	mockRequester := &mocks.Requester{}
	mockRequester.On("Do", mock.Anything).Return(nil, errors.New("test"))

	ext := ExtAPI{
		ContentServiceURL: "test",
		Client: mockRequester,
	}

	err := ext.SendToContentService(context.TODO(), *bytes.NewBuffer(nil), "test")
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestExternal_SendToContentService_ShouldReturnErrorIfResponseStatusCodeIsNot200(t *testing.T) {
	mockRequester := &mocks.Requester{}
	mockRequester.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusTeapot}, nil)

	ext := ExtAPI{
		ContentServiceURL: "test",
		Client: mockRequester,
	}

	err := ext.SendToContentService(context.TODO(), *bytes.NewBuffer(nil), "test")
	require.NotNil(t, err)
	require.Equal(t, "non-200 status code received: 418", err.Error())
}

func TestExternal_SendToContentService_ShouldReturnNoErrorIfResponseIs200(t *testing.T) {
	mockRequester := &mocks.Requester{}
	mockRequester.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	ext := ExtAPI{
		ContentServiceURL: "test",
		Client: mockRequester,
	}

	err := ext.SendToContentService(context.TODO(), *bytes.NewBuffer(nil), "test")
	require.Nil(t, err)
}
