package server

import (
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) handleUpdateMedicationHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		medicationHistoryID, errr := strconv.ParseUint(c.Param("medicationHistoryID"), 10, 32)
		if errr != nil {
			response.JSON(c, "invalid ID", http.StatusBadRequest, nil, errr)
			return
		}
		medicationHistoryRequest := struct {
			HasMedicationBeenTaken bool `json:"has_medication_been_taken"`
		}{}
		if err := decode(c, &medicationHistoryRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		err = s.MedicationHistoryService.UpdateMedicationHistory(medicationHistoryRequest.HasMedicationBeenTaken, uint(medicationHistoryID), user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication history updated successfully", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleGetAllMedicationHistoryByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		medicationHistories, err := s.MedicationHistoryService.GetAllMedicationHistoryByUser(user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication history retrieved successfully", http.StatusOK, medicationHistories, nil)
	}
}
