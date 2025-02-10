package http

import (
	"net/http"

	"github.com/CamilleLange/XM-company/internal/features/company"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginparamsmapper "gitlab.com/Zandraz/gin-params-mapper"
	"go.uber.org/zap"
)

// CompagnyHandler for the company feature.
type CompagnyHandler struct {
	companyFeatures company.CompagnyFeatures
}

// NewCompagnyHandler is a factory method for the CompagnyHandler type.
func NewCompagnyHandler(companyFeatures company.CompagnyFeatures) *CompagnyHandler {
	return &CompagnyHandler{
		companyFeatures: companyFeatures,
	}
}

// RegisterRoutes for the CompagnyHangler, provide a ready *gin.Engine.
func (h *CompagnyHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/company", h.Post)
	router.GET("/company/:company_uuid", h.Get)
	router.GET("/company", h.GetAll)
	router.PUT("/company/:company_uuid", h.Put)
	router.DELETE("/company/:company_uuid", h.Delete)
}

// Post parse the HTTP request in order to create the company.
func (h *CompagnyHandler) Post(c *gin.Context) {
	var company company.CompagnyCreateDTO
	if err := c.ShouldBindJSON(&company); err != nil {
		log.Error("CompagnyHandler.Post fail", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request: can't parse request body")
		return
	}

	if err := ValidateInstance.Struct(&company); err != nil {
		log.Error("CompagnyHandler.Post fail :", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body, data contains unsupported format.")
		return
	}

	companyUUID, err := h.companyFeatures.Create(company)
	if err != nil {
		log.Error("CompagnyHandler.Post fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusCreated, companyUUID)
}

// Get parse the HTTP request in order to respond the requested Compagny.
func (h *CompagnyHandler) Get(c *gin.Context) {
	var companyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("company_uuid", c, &companyUUID); err != nil {
		log.Error("CompagnyHandler.Get fail to get path param company_uuid:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	company, err := h.companyFeatures.ReadByID(companyUUID)
	if err != nil {
		log.Error("CompagnyHandler.Get fail to get company by ID:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, company)
}

// GetAll parse the HTTP request in order to respond with all Compagnies.
func (h *CompagnyHandler) GetAll(c *gin.Context) {
	company, err := h.companyFeatures.ReadAll()
	if err != nil {
		log.Error("CompagnyHandler.GetAll fail",
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

// Put parse the HTTP request in order to update the requested Compagny.
func (h *CompagnyHandler) Put(c *gin.Context) {
	var companyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("company_uuid", c, &companyUUID); err != nil {
		log.Error("CompagnyHandler.Put fail to get path param company_uuid:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	var company company.CompagnyUpdateDTO
	if err := c.ShouldBindJSON(&company); err != nil {
		log.Error("CompagnyHandler.Put fail", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request: can't parse request body")
		return
	}

	if err := ValidateInstance.Struct(company); err != nil {
		log.Error("CompagnyHandler.Put fail :", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body, data contains unsupported format.")
		return
	}

	if err := h.companyFeatures.Update(companyUUID, company); err != nil {
		log.Error("CompagnyHandler.Put fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, "Compagny updated.")
}

// Delete parse the HTTP request in order to delete the requested Compagny.
func (h *CompagnyHandler) Delete(c *gin.Context) {
	var companyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("company_uuid", c, &companyUUID); err != nil {
		log.Error("CompagnyHandler.Delete fail to get path param company_uuid: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	err := h.companyFeatures.Delete(companyUUID)
	if err != nil {
		log.Error("CompagnyHandler.Delete fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusNoContent, "Compagny deleted.")
}
