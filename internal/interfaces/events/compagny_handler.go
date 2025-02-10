package events

import (
	"github.com/CamilleLange/XM-compagny/internal/features/compagny"
	"github.com/nats-io/nats.go"
)

type CompagnyEventHandler struct {
	compagnyFeatures compagny.CompagnyFeatures
}

func NewCompagnyEventHandler(compagnyFeatures compagny.CompagnyFeatures) *CompagnyEventHandler {
	return &CompagnyEventHandler{
		compagnyFeatures: compagnyFeatures,
	}
}

func (h *CompagnyEventHandler) RegisterEvents(nats *nats.Conn) error {
	return nil
}
