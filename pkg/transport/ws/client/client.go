package client

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/KarinaBotova/Normalization/models"
)

// сокет
func SendStudents(Students []models.Student, host, port string) error {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", host, port), Path: "/"}
	log.Info().Msgf("Connecting to %s", u.String())
	// инициализируем соединение
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to dial")
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer c.Close()
	// отправляем всех студентов по сокету
	for _, newStudent := range Students {
		if err = c.WriteJSON(newStudent); err != nil {
			log.Error().Err(err).Msg("Failed to write msg")
			return fmt.Errorf("failed to write msg: %v", err)
		}
	}
	// закрываем соединение
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Error().Err(err).Msg("Failed to write close msg")
		return fmt.Errorf("failed to write close msg: %v", err)
	}

	return nil
}
