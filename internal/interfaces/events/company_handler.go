package events

import (
	"github.com/CamilleLange/XM-company/internal/features/company"
	"github.com/nats-io/nats.go"
)

type CompanyEventHandler struct {
	companyFeatures company.CompanyFeatures
}

func NewCompanyEventHandler(companyFeatures company.CompanyFeatures) *CompanyEventHandler {
	return &CompanyEventHandler{
		companyFeatures: companyFeatures,
	}
}

func (h *CompanyEventHandler) RegisterEvents(nats *nats.Conn) error {
	return nil
}
