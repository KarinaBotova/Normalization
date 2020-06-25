package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/KarinaBotova/Normalization/models"
	postgresdriver "github.com/KarinaBotova/Normalization/pkg/postgres-driver"
)

func postStudent(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get body")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to read body"))
		return
	}
	defer r.Body.Close()
	var Students []models.Student
	if err = json.Unmarshal(data, &Students); err != nil {
		log.Error().Err(err).Msg("Failed to decode body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("Failed to decode body"))
		return
	}

	// Send data in DB
	if err = postgresdriver.SaveStudents(Students); err != nil {
		log.Error().Err(err).Msg("Could not send Students to DB")
	}
	log.Info().Msgf("Students %v was saved via REST", Students)
	w.WriteHeader(http.StatusOK)
}

// запрос-ответ
func StartServer(host, port string, done chan interface{}) {
	// создаем роутер
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Yesss"))
		w.WriteHeader(200)
	})

	r.Post("/", postStudent)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
		// TODO: TLS
		TLSConfig:    nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		// TODO Shutdown?
	}

	// Start server
	log.Info().Msgf("Server is listening on %s:%s", host, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
