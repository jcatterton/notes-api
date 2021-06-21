package api

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"notes-api/pkg/models"
	"notes-api/pkg/testhelper/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAPI_CheckHealth_ShouldRespondWith500IfErrorOccursPingingDatabase(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("Ping", mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(checkHealth(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_CheckHealth_ShouldRespondWith200IfNoErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("Ping", mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(checkHealth(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), "API is running and connected to database")
}

func TestAPI_GetNotes_ShouldRespondWith400IfErrorOccursRetrievingAuthToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNotes(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "no authorization header found")
}

func TestAPI_GetNotes_ShouldRespondWith400IfTokenFormatIsInvalid(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNotes(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "authorization header must be in format 'Bearer'")
}

func TestAPI_GetNotes_ShouldRespondWith401IfErrorOccursValidatingToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNotes(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_GetNotes_ShouldRespondWith500IfServiceErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("GetNotes", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNotes(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_GetNotes_ShouldRespondWith200OnSuccess(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{}, nil)

	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNotes(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAPI_GetNote_ShouldRespondWith400IfErrorOccursRetrievingAuthToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "no authorization header found")
}

func TestAPI_GetNote_ShouldRespondWith400IfTokenFormatIsInvalid(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "authorization header must be in format 'Bearer'")
}

func TestAPI_GetNote_ShouldRespondWith401IfErrorOccursValidatingToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_GetNote_ShouldRespondWith500IfServiceErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("GetNotes", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_GetNote_ShouldRespondWith500IfMoreThanOneNoteIsReturned(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{{}, {}}, nil)

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "more than one note returned for given ID")
}

func TestAPI_GetNote_ShouldRespondWith204IfNoNotesAreReturned(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{}, nil)

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestAPI_GetNote_ShouldRespondWith200OnSuccess(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{{}}, nil)

	req, err := http.NewRequest(http.MethodGet, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAPI_EditNote_ShouldRespondWith400IfErrorOccursRetrievingAuthToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodPut, "/note", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(editNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "no authorization header found")
}

func TestAPI_EditNote_ShouldRespondWith400IfTokenFormatIsInvalid(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodPut, "/note", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(editNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "authorization header must be in format 'Bearer'")
}

func TestAPI_EditNote_ShouldRespondWith401IfErrorOccursValidatingToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodPut, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(editNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_EditNote_ShouldRespondWith400IfErrorOccursDecodingRequestBody(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodPut, "/note", ioutil.NopCloser(strings.NewReader("")))
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(editNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "EOF")
}

func TestAPI_EditNote_ShouldRespondWith500IfServiceErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodPut, "/note", ioutil.NopCloser(strings.NewReader("{}")))
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(editNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_EditNote_ShouldRespondWith200OnSuccess(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodPut, "/note", ioutil.NopCloser(strings.NewReader("{}")))
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(editNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAPI_CreateNote_ShouldRespondWith400IfErrorOccursRetrievingAuthToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(createNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "no authorization header found")
}

func TestAPI_CreateNote_ShouldRespondWith400IfTokenFormatIsInvalid(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(createNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "authorization header must be in format 'Bearer'")
}

func TestAPI_CreateNote_ShouldRespondWith401IfErrorOccursValidatingToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(createNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_CreateNote_ShouldRespondWith400IfErrorOccursDecodingRequestBody(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodPost, "/note", ioutil.NopCloser(strings.NewReader("")))
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(createNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "EOF")
}

func TestAPI_CreateNote_ShouldRespondWith500IfServiceErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("CreateNote", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))

	req, err := http.NewRequest(http.MethodPost, "/note", ioutil.NopCloser(strings.NewReader("{}")))
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(createNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_CreateNote_ShouldRespondWith200OnSuccess(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("CreateNote", mock.Anything, mock.Anything, mock.Anything).Return("", nil)

	req, err := http.NewRequest(http.MethodPost, "/note", ioutil.NopCloser(strings.NewReader("{}")))
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(createNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAPI_DeleteNote_ShouldRespondWith400IfErrorOccursRetrievingAuthToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodDelete, "/note", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(deleteNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "no authorization header found")
}

func TestAPI_DeleteNote_ShouldRespondWith400IfTokenFormatIsInvalid(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodDelete, "/note", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(deleteNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "authorization header must be in format 'Bearer'")
}

func TestAPI_DeleteNote_ShouldRespondWith401IfErrorOccursValidatingToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodDelete, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(deleteNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_DeleteNote_ShouldRespondWith500IfServiceErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("DeleteNote", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodDelete, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(deleteNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_DeleteNote_ShouldRespondWith200OnSuccess(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("DeleteNote", mock.Anything, mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(deleteNote(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAPI_SendToContentService_ShouldRespondWith400IfErrorOccursRetrievingAuthToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodPost, "/save", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sendToContentService(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "no authorization header found")
}

func TestAPI_SendToContentService_ShouldRespondWith400IfTokenFormatIsInvalid(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sendToContentService(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "authorization header must be in format 'Bearer'")
}

func TestAPI_SendToContentService_ShouldRespondWith401IfErrorOccursValidatingToken(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sendToContentService(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_SendToContentService_ShouldRespondWith500IfServiceErrorOccurs(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("SetToken", mock.Anything).Return()
	mockSvc.On("SendToContentService", mock.Anything, mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sendToContentService(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "test")
}

func TestAPI_SendToContentService_ShouldRespondWith200OnSuccess(t *testing.T) {
	mockSvc := &mocks.NoteServiceHandler{}
	mockSvc.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)
	mockSvc.On("SetToken", mock.Anything).Return()
	mockSvc.On("SendToContentService", mock.Anything, mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodPost, "/note", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer test")

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sendToContentService(context.TODO(), mockSvc))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}
