package http

import (
	"net/http"

	"github.com/CamilleLange/XM-compagny/internal/features/compagny"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginparamsmapper "gitlab.com/Zandraz/gin-params-mapper"
	"go.uber.org/zap"
)

type CompagnyHandler struct {
	compagnyFeatures compagny.CompagnyFeatures
}

func NewCompagnyHandler(compagnyFeatures compagny.CompagnyFeatures) *CompagnyHandler {
	return &CompagnyHandler{
		compagnyFeatures: compagnyFeatures,
	}
}

func (h *CompagnyHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/compagny", h.Post)
	router.GET("/compagny/:compagny_uuid", h.Get)
	router.GET("/compagny", h.GetAll)
	router.PUT("/compagny/:compagny_uuid", h.Put)
	router.DELETE("/compagny/:compagny_uuid", h.Delete)
}

func (h *CompagnyHandler) Post(c *gin.Context) {
	var compagny compagny.CompagnyCreateDTO
	if err := c.ShouldBindJSON(&compagny); err != nil {
		log.Error("CompagnyHandler.Post fail", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request: can't parse request body")
		return
	}

	log.Debug("compagny", zap.Any("c", compagny))

	if err := ValidateInstance.Struct(&compagny); err != nil {
		log.Error("CompagnyHandler.Post fail :", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body, data contains unsupported format.")
		return
	}

	compagnyUUID, err := h.compagnyFeatures.Create(compagny)
	if err != nil {
		log.Error("CompagnyHandler.Post fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusCreated, compagnyUUID)
}

func (h *CompagnyHandler) Get(c *gin.Context) {
	var compagnyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("compagny_uuid", c, &compagnyUUID); err != nil {
		log.Error("CompagnyHandler.Get fail to get path param compagny_uuid:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	compagny, err := h.compagnyFeatures.ReadByID(compagnyUUID)
	if err != nil {
		log.Error("CompagnyHandler.Get fail to get compagny by ID:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, compagny)
}

func (h *CompagnyHandler) GetAll(c *gin.Context) {
	compagny, err := h.compagnyFeatures.ReadAll()
	if err != nil {
		log.Error("CompagnyHandler.GetAll fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	if len(compagny) == 0 {
		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusOK, compagny)
	}
}

func (h *CompagnyHandler) Put(c *gin.Context) {
	var compagnyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("compagny_uuid", c, &compagnyUUID); err != nil {
		log.Error("CompagnyHandler.Put fail to get path param compagny_uuid:", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	var compagny compagny.CompagnyUpdateDTO
	if err := c.ShouldBindJSON(&compagny); err != nil {
		log.Error("CompagnyHandler.Put fail", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request: can't parse request body")
		return
	}

	if err := ValidateInstance.Struct(compagny); err != nil {
		log.Error("CompagnyHandler.Put fail :", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body, data contains unsupported format.")
		return
	}

	if err := h.compagnyFeatures.Update(compagnyUUID, compagny); err != nil {
		log.Error("CompagnyHandler.Put fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, "Compagny updated.")
}

func (h *CompagnyHandler) Delete(c *gin.Context) {
	var compagnyUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("compagny_uuid", c, &compagnyUUID); err != nil {
		log.Error("CompagnyHandler.Delete fail to get path param compagny_uuid: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	err := h.compagnyFeatures.Delete(compagnyUUID)
	if err != nil {
		log.Error("CompagnyHandler.Delete fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusNoContent, "Compagny deleted.")
}
