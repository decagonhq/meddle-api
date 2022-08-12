package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) handleCreateMedication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var medicationRequest models.MedicationRequest
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		userId := user.ID
		if err := c.ShouldBindJSON(&medicationRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		medicationRequest.UserID = userId
		createdMedication, err := s.MedicationService.CreateMedication(&medicationRequest)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication created successful", http.StatusCreated, createdMedication, nil)
	}
}

func (s *Server) handleGetMedDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		id := c.Param("id")
		idUint, errr := strconv.ParseUint(id, 10, 32)
		if errr != nil {
			response.JSON(c, "error parsing id", http.StatusBadRequest, nil, errr)
			return
		}
		medication, err := s.MedicationService.GetMedicationDetail(uint(idUint), user.ID)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "retrieved medications successfully", http.StatusOK, gin.H{"medication": medication}, nil)
	}
}

func (s *Server) handleGetAllMedications() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		medications, err := s.MedicationService.GetAllMedications(user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medications retrieved successfully", http.StatusOK, medications, nil)
	}
}

func (s *Server) handleGetNextMedication() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}

		medication, err := s.MedicationService.GetNextMedications(user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication retrieved successfully", http.StatusOK, medication, nil)
	}
}
