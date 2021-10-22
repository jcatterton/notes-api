package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"notes-api/pkg/models"
	"notes-api/pkg/testhelper/mocks"
)

func TestService_Ping_ShouldReturnErrorIfDaoErrors(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("Ping", mock.Anything).Return(errors.New("test"))

	service := NotesService{
		Dao: mockDao,
	}

	err := service.Ping(context.TODO())
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_Ping_ShouldReturnNilIfNoErrorOccurs(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("Ping", mock.Anything).Return(nil)

	service := NotesService{
		Dao: mockDao,
	}

	require.Nil(t, service.Ping(context.TODO()))
}

func TestService_GetNotes_ShouldReturnErrorIfIDIsNotEmptyAndIsNotValidHex(t *testing.T) {
	service := NotesService{}

	notes, err := service.GetNotes(context.TODO(), "test")
	require.Nil(t, notes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "encoding/hex: invalid byte")
}

func TestService_GetNotes_ShouldReturnErrorOnDaoError(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("GetNotes", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

	service := NotesService{
		Dao: mockDao,
	}

	notes, err := service.GetNotes(context.TODO(), "")
	require.Nil(t, notes)
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_GetNotes_ShouldReturnNoErrorIfNoErrorOccurs(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{}, nil)

	service := NotesService{
		Dao: mockDao,
	}

	notes, err := service.GetNotes(context.TODO(), "")
	require.Nil(t, err)
	require.NotNil(t, notes)
}

func TestService_UpdateNote_ShouldReturnErrorIfIDIsNotValidHex(t *testing.T) {
	service := NotesService{}

	err := service.UpdateNote(context.TODO(), "test", models.NoteRequest{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "encoding/hex: invalid byte")
}

func TestService_UpdateNote_ShouldReturnErrorOnDaoError(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

	service := NotesService{
		Dao: mockDao,
	}

	err := service.UpdateNote(context.TODO(), "000000000000000000000000", models.NoteRequest{})
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_UpdateNote_ShouldReturnNoErrorIfNoErrorOccurs(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NotesService{
		Dao: mockDao,
	}

	require.Nil(t, service.UpdateNote(context.TODO(), "000000000000000000000000", models.NoteRequest{}))
}

func TestService_DeleteNote_ShouldReturnErrorIfIDIsNotValidHex(t *testing.T) {
	service := NotesService{}

	err := service.DeleteNote(context.TODO(), "test")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "encoding/hex: invalid byte")
}

func TestService_DeleteNote_ShouldReturnErrorOnDaoError(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("DeleteNote", mock.Anything, mock.Anything).Return(errors.New("test"))

	service := NotesService{
		Dao: mockDao,
	}

	err := service.DeleteNote(context.TODO(), "000000000000000000000000")
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_DeleteNote_ShouldReturnNoErrorIfNoErrorOccurs(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("DeleteNote", mock.Anything, mock.Anything).Return(nil)

	service := NotesService{
		Dao: mockDao,
	}

	require.Nil(t, service.DeleteNote(context.TODO(), "000000000000000000000000"))
}

func TestService_CreateNote_ShouldReturnErrorOnDaoError(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("CreateNote", mock.Anything, mock.Anything).Return(errors.New("test"))

	service := NotesService{
		Dao: mockDao,
	}

	id, err := service.CreateNote(context.TODO(), models.NoteRequest{})
	require.Equal(t, "", id)
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_CreateNote_ShouldReturnNoteIDAndNilErrorIfNoErrorOccurs(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("CreateNote", mock.Anything, mock.Anything).Return(nil)

	service := NotesService{
		Dao: mockDao,
	}

	id, err := service.CreateNote(context.TODO(), models.NoteRequest{})
	require.NotEqual(t, "", id)
	require.Nil(t, err)
}

func TestService_SendToContentService_ShouldReturnErrorOnDaoError(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{}, errors.New("test"))

	service := NotesService{
		Dao: mockDao,
	}

	err := service.SendToContentService(context.TODO(), "")
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_SendToContentService_ShouldReturnErrorOnExtHandlerError(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{{}}, nil)

	mockExt := &mocks.ExtAPIHandler{}
	mockExt.On("SendToContentService", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

	service := NotesService{
		Dao: mockDao,
		Ext: mockExt,
	}

	err := service.SendToContentService(context.TODO(), "")
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_SendToContentService_ShouldReturnNoErrorIfNoErrorOccurs(t *testing.T) {
	mockDao := &mocks.NoteDaoHandler{}
	mockDao.On("GetNotes", mock.Anything, mock.Anything).Return([]models.Note{{}}, nil)

	mockExt := &mocks.ExtAPIHandler{}
	mockExt.On("SendToContentService", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NotesService{
		Dao: mockDao,
		Ext: mockExt,
	}

	require.Nil(t, service.SendToContentService(context.TODO(), ""))
}

func TestService_ValidateToken_ShouldReturnErrorOnExtHandlerError(t *testing.T) {
	mockExt := &mocks.ExtAPIHandler{}
	mockExt.On("ValidateToken", mock.Anything, mock.Anything).Return(errors.New("test"))

	service := NotesService{
		Ext: mockExt,
	}

	err := service.ValidateToken(context.TODO(), "test")
	require.NotNil(t, err)
	require.Equal(t, "test", err.Error())
}

func TestService_ValidateToken_ShouldReturnNoErrorIfNoErrorOccurs(t *testing.T) {
	mockExt := &mocks.ExtAPIHandler{}
	mockExt.On("ValidateToken", mock.Anything, mock.Anything).Return(nil)

	service := NotesService{
		Ext: mockExt,
	}

	require.Nil(t, service.ValidateToken(context.TODO(), "test"))
}

func TestService_SetToken_ShouldSetToken(t *testing.T) {
	mockExt := &mocks.ExtAPIHandler{}
	mockExt.On("SetToken", mock.Anything).Return()

	service := NotesService{
		Ext: mockExt,
	}

	service.SetToken("test")
}
