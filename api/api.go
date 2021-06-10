package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"notes-api/dao"
	"notes-api/models"
	"notes-api/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListenAndServe(ctx context.Context) error {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	router, err := route(ctx)
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler:      handlers.CORS(headers, origins, methods)(router),
		Addr:         ":8006",
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}
	shutdownGracefully(server)

	logrus.WithContext(ctx).Info("Starting API server...")
	return server.ListenAndServe()
}

func route(ctx context.Context) (*mux.Router, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		logrus.WithError(err).Error("Error creating mongo client")
		return nil, err
	}

	notesDao := dao.NotesDao{
		Client:     client,
		Database:   os.Getenv("DATABASE"),
		Collection: os.Getenv("COLLECTION"),
	}

	notesService := service.NotesService{
		Dao: notesDao,
	}

	router := mux.NewRouter()
	router.Handle("/health", checkHealth(ctx, &notesService)).Methods(http.MethodGet)
	router.Handle("/notes", getNotes(ctx, &notesService)).Methods(http.MethodGet)
	router.Handle("/note/{id}", getNote(ctx, &notesService)).Methods(http.MethodGet)
	router.Handle("/note/{id}", editNote(ctx, &notesService)).Methods(http.MethodPut)
	router.Handle("/note/{id}", deleteNote(ctx, &notesService)).Methods(http.MethodDelete)
	router.Handle("/note", createNote(ctx, &notesService)).Methods(http.MethodPost)

	return router, nil
}

func checkHealth(ctx context.Context, svc service.NoteServiceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithContext(ctx)
		defer closeRequestBody(ctx, r)

		if err := svc.Ping(ctx); err != nil {
			logger.WithError(err).Error("Error connecting to database")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithSuccess(ctx, w, http.StatusOK, "API is running and connected to database")
	}
}

func getNotes(ctx context.Context, svc service.NoteServiceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithContext(ctx)
		defer closeRequestBody(ctx, r)

		notes, err := svc.GetNotes(ctx, "")
		if err != nil {
			logger.WithError(err).Error("Error retrieving notes")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithSuccess(ctx, w, http.StatusOK, notes)
	}
}

func getNote(ctx context.Context, svc service.NoteServiceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithContext(ctx)
		defer closeRequestBody(ctx, r)

		id := mux.Vars(r)["id"]

		notes, err := svc.GetNotes(ctx, id)
		if err != nil {
			logger.WithError(err).Error("Error retrieving notes")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		} else if len(notes) > 1 {
			err := errors.New("more than one note returned for given ID")
			logger.WithError(err).Error("Invalid note results")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		} else if len(notes) == 0 {
			respondWithSuccess(ctx, w, http.StatusNoContent, nil)
			return
		}

		respondWithSuccess(ctx, w, http.StatusOK, notes[0])
	}
}

func editNote(ctx context.Context, svc service.NoteServiceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithContext(ctx)
		defer closeRequestBody(ctx, r)

		id := mux.Vars(r)["id"]

		var note models.NoteRequest
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			logger.WithError(err).Error("Error decoding request body")
			respondWithError(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := svc.UpdateNote(ctx, id, note); err != nil {
			logger.WithError(err).Error("Error updating note")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithSuccess(ctx, w, http.StatusOK, fmt.Sprintf("Note with ID '%v' updated successfully", id))
	}
}

func createNote(ctx context.Context, svc service.NoteServiceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithContext(ctx)
		defer closeRequestBody(ctx, r)

		var note models.NoteRequest
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			logger.WithError(err).Error("Error decoding request body")
			respondWithError(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		id, err := svc.CreateNote(ctx, note)
		if err != nil {
			logger.WithError(err).Error("Error creating note")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithSuccess(ctx, w, http.StatusOK, fmt.Sprintf("Note with ID '%v' created successfuly", id))
	}
}

func deleteNote(ctx context.Context, svc service.NoteServiceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithContext(ctx)
		defer closeRequestBody(ctx, r)

		id := mux.Vars(r)["id"]

		if err := svc.DeleteNote(ctx, id); err != nil {
			logger.WithError(err).Error("Error deleting note")
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithSuccess(ctx, w, http.StatusOK, fmt.Sprintf("Note with ID '%v' deleted successfully", id))
	}
}

func respondWithError(ctx context.Context, w http.ResponseWriter, code int, message string) {
	logger := logrus.WithContext(ctx)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if message == "" {
		logger.WithError(errors.New("message body is nil")).Error("Unable to write response")
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		logger.WithError(err).Error("Error encoding response")
	}
}

func respondWithSuccess(ctx context.Context, w http.ResponseWriter, code int, body interface{}) {
	logger := logrus.WithContext(ctx)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger.WithError(err).Error("Error encoding response")
		}
	}
}

func closeRequestBody(ctx context.Context, r *http.Request) {
	logger := logrus.WithContext(ctx)

	if r.Body == nil {
		return
	}

	if err := r.Body.Close(); err != nil {
		logger.WithError(err).Error("Error closing request body")
	}
}

func shutdownGracefully(server *http.Server) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals

		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(c); err != nil {
			logrus.WithError(err).Error("Error shutting down server")
		}

		<-c.Done()
		os.Exit(0)
	}()
}
