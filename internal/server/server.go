package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andreykaipov/goobs"
	"github.com/google/uuid"

	"github.com/dreamsofcode-io/obs-remote/internal/middleware"
	"github.com/dreamsofcode-io/obs-remote/internal/recording"
)

type validationError struct {
	Reason string `json:"reason"`
}

type UpdateBody struct {
	Filename string `json:"filename"`
}

func Start(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := recording.NewRepository()

	obsClient, err := goobs.New("localhost:4455", goobs.WithPassword("ANPtuxVGio4BpvZp"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", func(_ http.ResponseWriter, _ *http.Request) {
		// Health handler, don't need to do anything currently
	})

	router.HandleFunc("POST /record/start", func(w http.ResponseWriter, r *http.Request) {
		obsClient.Record.StartRecord()
	})

	router.HandleFunc("POST /record/stop", func(w http.ResponseWriter, r *http.Request) {
		res, err := obsClient.Record.StopRecord()
		if err != nil {
			slog.Error("failed to stop recording", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		id, err := uuid.NewV7()
		if err != nil {
			slog.Error("failed to generate uuid", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rec := recording.Recording{
			ID:   id,
			Path: res.OutputPath,
		}

		if err = repo.Insert(r.Context(), rec); err != nil {
			slog.Error(
				"failed to write recording",
				slog.Any("error", err),
				slog.String("outputpath", res.OutputPath),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(rec); err != nil {
			slog.Error(
				"failed to marshal on json",
				slog.Any("error", err),
				slog.Any("recording", rec),
			)
		}
	})

	router.HandleFunc("PUT /recordings/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			json.NewEncoder(w).Encode(validationError{
				Reason: "id should be a valid uuid",
			})
			w.WriteHeader(http.StatusBadRequest)
		}

		var body UpdateBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		if body.Filename == "" {
			json.NewEncoder(w).Encode(validationError{
				Reason: "filename cannot be empty",
			})

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		rec, err := repo.FindByID(ctx, id)
		if errors.Is(err, recording.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		s, err := os.Stat(rec.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if s.IsDir() {
			logger.Error("file is a directory", slog.String("file", rec.Path))
			w.WriteHeader(http.StatusInternalServerError)
		}

		dir := filepath.Dir(rec.Path)
		ext := filepath.Ext(rec.Path)

		newPath := fmt.Sprintf("%s/%s%s", dir, body.Filename, ext)

		if err = os.Rename(rec.Path, newPath); err != nil {
			logger.Error("file is a directory", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rec.Path = newPath

		if err := repo.Update(r.Context(), rec); err != nil {
			logger.Error("failed to update the record", slog.Any("error", err))
			return
		}
	})

	srv := http.Server{
		Addr:    ":3000",
		Handler: middleware.Logging(logger)(router),
	}

	fmt.Println("Listening and serving")
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start server", slog.Any("error", err))
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}
