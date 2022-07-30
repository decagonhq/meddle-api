package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) handleGetMedDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
			medication, err := s.MedicationService.GetMedicationDetail(user.ID)
			if err != nil {
				log.Printf("get medications error : %v\n", err)
				response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
				return
			}
			response.JSON(c, "retrieved medications successfully", http.StatusOK, gin.H{"medication": medication}, nil)
			return
	}
}
