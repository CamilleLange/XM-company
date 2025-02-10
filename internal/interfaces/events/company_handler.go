package events

import (
	"github.com/CamilleLange/XM-company/internal/features/company"
	"github.com/nats-io/nats.go"
)

type CompagnyEventHandler struct {
	companyFeatures company.CompagnyFeatures
}

func NewCompagnyEventHandler(companyFeatures company.CompagnyFeatures) *CompagnyEventHandler {
	return &CompagnyEventHandler{
		companyFeatures: companyFeatures,
	}
}

func (h *CompagnyEventHandler) RegisterEvents(nats *nats.Conn) error {
	return nil
}
