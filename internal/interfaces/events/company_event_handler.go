package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	CREATE EventType = "create"
	UPDATE EventType = "update"
	DELETE EventType = "delete"
)

// EventType is an enum for the event type.
type EventType string

// Event is the data model for all event on Kafka.
type Event struct {
	UUID    uuid.UUID `json:"uuid"`
	Type    EventType `json:"type"`
	TS      time.Time `json:"ts"`
	Payload any       `json:"payload"`
}

// CompanyEventHandler encapsulate all operation on Kafka for the company feature.
type CompanyEventHandler struct {
	writer *kafka.Writer
}

// WriteMessage on Kafka.
func (h *CompanyEventHandler) WriteMessage(e Event) error {
	eventData, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("can't marshal event in JSON : %w", err)
	}

	err = h.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(e.UUID.String()),
		Value: eventData,
	})
	if err != nil {
		return fmt.Errorf("can't write message to kafka : %w", err)
	}

	return nil
}

// Close Kafka connection.
func (h *CompanyEventHandler) Close() error {
	if err := h.writer.Close(); err != nil {
		return fmt.Errorf("can't close kafka writer : %w", err)
	}

	return nil
}

// NewCompanyEventHandler is a factory method for company event handler.
func NewCompanyEventHandler(conf Config) (*CompanyEventHandler, error) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(conf.BrockerAddr),
		Topic:    conf.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	fmt.Println(conf)

	return &CompanyEventHandler{
		writer: w,
	}, nil
}

// NewEvent is a factory method to create an Event.
// This function will generate an event UUID and set the timestamp of the Event.
func NewEvent(payload any, eventType EventType) Event {
	return Event{
		UUID:    uuid.New(),
		Type:    eventType,
		TS:      time.Now().UTC(),
		Payload: payload,
	}
}
