package http

import (
	"net/http"

	"github.com/CamilleLange/XM-company/internal/features/company"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginparamsmapper "gitlab.com/Zandraz/gin-params-mapper"
	"go.uber.org/zap"
)

// CompanyHandler for the company feature.
type CompanyHandler struct {
	companyFeatures company.CompanyFeatures
}

// NewCompanyHandler is a factory method for the CompanyHandler type.
func NewCompanyHandler(companyFeatures company.CompanyFeatures) *CompanyHandler {
	return &CompanyHandler{
		companyFeatures: companyFeatures,
	}
}

// RegisterRoutes for the CompanyHangler, provide a ready *gin.Engine.
func (h *CompanyHandler) RegisterRoutes(router *gin.Engine) {
	// public endpoints.
	router.GET("/company", h.GetAll)
	router.GET("/company/:company_uuid", h.Get)

	// Protected endpoints.
	router.Group("/company").
		Use(BasicJWTMiddleware()).
		POST("", h.Post).
		PUT("/:company_uuid", h.Put).
		DELETE("/:company_uuid", h.Delete)
}

// Post parse the HTTP request in order to create the company.
func (h *CompanyHandler) Post(c *gin.Context) {
	var company company.CompanyCreateDTO
	if err := c.ShouldBindJSON(&company); err != nil {
		log.Error("CompanyHandler.Post fail", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request: can't parse request body")
		return
	}

	if err := ValidateInstance.Struct(&company); err != nil {
		log.Error("CompanyHandler.Post fail :", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body, data contains unsupported format.")
		return
	}

	companyUUID, err := h.companyFeatures.Create(company)
	if err != nil {
		log.Error("CompanyHandler.Post fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusCreated, companyUUID)
}

// Get parse the HTTP request in order to respond the requested Company.
func (h *CompanyHandler) Get(c *gin.Context) {
	var companyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("company_uuid", c, &companyUUID); err != nil {
		log.Error("CompanyHandler.Get fail to get path param company_uuid:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	company, err := h.companyFeatures.ReadByID(companyUUID)
	if err != nil {
		log.Error("CompanyHandler.Get fail to get company by ID:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, company)
}

// GetAll parse the HTTP request in order to respond with all Compagnies.
func (h *CompanyHandler) GetAll(c *gin.Context) {
	company, err := h.companyFeatures.ReadAll()
	if err != nil {
		log.Error("CompanyHandler.GetAll fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	if len(company) == 0 {
		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusOK, company)
	}
}

// Put parse the HTTP request in order to update the requested Company.
func (h *CompanyHandler) Put(c *gin.Context) {
	var companyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("company_uuid", c, &companyUUID); err != nil {
		log.Error("CompanyHandler.Put fail to get path param company_uuid:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	var company company.CompanyUpdateDTO
	if err := c.ShouldBindJSON(&company); err != nil {
		log.Error("CompanyHandler.Put fail", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request: can't parse request body")
		return
	}

	if err := ValidateInstance.Struct(company); err != nil {
		log.Error("CompanyHandler.Put fail :", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body, data contains unsupported format.")
		return
	}

	if err := h.companyFeatures.Update(companyUUID, company); err != nil {
		log.Error("CompanyHandler.Put fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Delete parse the HTTP request in order to delete the requested Company.
func (h *CompanyHandler) Delete(c *gin.Context) {
	var companyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("company_uuid", c, &companyUUID); err != nil {
		log.Error("CompanyHandler.Delete fail to get path param company_uuid: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	err := h.companyFeatures.Delete(companyUUID)
	if err != nil {
		log.Error("CompanyHandler.Delete fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
