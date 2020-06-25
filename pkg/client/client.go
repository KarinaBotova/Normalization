package client

import (
	"context"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	models "github.com/KarinaBotova/Normalization/models"
	pubsubclient "github.com/KarinaBotova/Normalization/pkg/transport/mq/client"
	restclient "github.com/KarinaBotova/Normalization/pkg/transport/rest/client"
	wsclient "github.com/KarinaBotova/Normalization/pkg/transport/ws/client"
)

var opts struct {
	Server    string `long:"server" env:"SERVER" required:"true"`
	Port      string `long:"port" env:"PORT" description:"Server port" required:"true"`
	ProjectID string `long:"projectID" env:"PROJECT_ID" required:"true" default:"trrp-virus"`
	RESTPort  string `long:"rest_port" env:"REST_PORT" description:"Server port" required:"true"`
	WSPort    string `long:"ws_port" env:"WS_PORT" description:"Server port" required:"true"`
	LogLevel  string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
	Topic     string `long:"topic" env:"TOPIC" description:"Topic" required:"true"`
}

var (
	defaultStudents []models.Student
	client          *pubsubclient.Client
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Client is ready"))
		return
	case http.MethodPost: // считывание данных из запроса
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Failed to read data"))
			return
		}
		// преобразование данных
		var Students []models.Student
		if err = json.Unmarshal(data, &Students); err != nil {
			log.Error().Err(err).Msg("Failed to decode body")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Failed to decode body"))
			return
		}
		// извлекаем из запроса тип канала передачи сообщения. Очередь или сокет?
		t := r.Header.Get("Type")
		switch t {
		// запрос-ответ взаимодействие
		case "REST":
			if err := restclient.SendStudents(Students, opts.Server, opts.RESTPort); err != nil {
				log.Error().Err(err).Msg("Failed to send Students via REST")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send Students via REST"))
				return
			}
		case "WS": // web-сокет
			if err := wsclient.SendStudents(Students, opts.Server, opts.WSPort); err != nil {
				log.Error().Err(err).Msg("Failed to send Students via WS")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send Students via WS"))
				return
			}
		case "MQ": // очередь
			if err := client.SendStudents(Students); err != nil {
				log.Error().Err(err).Msg("Failed to send Students via MQ")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send Students via MQ"))
				return
			}

		default:
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Unrecognized header"))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method is not allowed"))
		return
	}

	log.Info().Msg("Students sent successfully")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Students sent successfully"))
}

// запуск клиента
func Start(done chan interface{}) {
	// инициализация логгера
	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	zerolog.ErrorFieldName = "ERROR"
	zerolog.TimestampFieldName = "TIME"
	zerolog.CallerFieldName = "CALLER"
	log.Logger = log.Output(os.Stderr).With().Str("PROGRAM", "easy-normalization").Caller().Logger()

	// считывание переменных в командной строке
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse flags")
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// инициализация соединения с очередью
	if client, err = pubsubclient.NewClient(opts.ProjectID, opts.Topic); err != nil {
		log.Fatal().Err(err).Msg("Failed to create new client")
	}
	// создаем сервер клиента
	r := http.NewServeMux()

	r.HandleFunc("/", handler)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%s", opts.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		ctx := context.Background()
		<-done
		_ = srv.Shutdown(ctx)
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
