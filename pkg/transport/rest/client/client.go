package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"

	pb "github.com/kolya59/easy_normalization/proto"
)

func SendStudents(Students []models.Student, host, port string) error {
	// Отправляет студентов
	data, err := json.Marshal(Students)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal data")
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	u := url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%s", host, port), Path: "/"}
	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(data))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send request")
		return fmt.Errorf("failed to send request: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		log.Info().Msgf("Response status is OK: %v", resp.Status)
		return nil
	} else {
		log.Error().Msgf("Response status is not OK: %v", resp.Status)
		return fmt.Errorf("response status is not OK: %v", resp.Status)
	}
}
