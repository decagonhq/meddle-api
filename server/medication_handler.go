package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleCreateMedication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var medicationRequest models.MedicationRequest

		userI, exists := c.Get("user")
		if !exists {
			err := errors.New("user not signed in", http.StatusUnauthorized)
			response.JSON(c, err.Message, err.Status, nil, err)
			return
		}
		userId := userI.(*models.User).ID

		if err := c.ShouldBindJSON(&medicationRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		medicationRequest.UserID = userId
		createdMedication, err := s.MedicationService.CreateMedication(&medicationRequest)
		if err != nil {
			err.Respond(c)
		}

		response.JSON(c, "medication created successful", http.StatusCreated, createdMedication, nil)
	}
}
