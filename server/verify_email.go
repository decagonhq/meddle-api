package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) handleVerifyEmail() gin.HandlerFunc {
	return func(context *gin.Context) {
		//Generate token
		//Send token as email to user with link to verify email

		//token := context.Param("userToken")

		//if token != user.ReturnedToken {
		//	log.Println("invalid token")
		//	response.JSON(context, "", http.StatusInternalServerError, nil, []string{"Invalid user token or ID"})
		//	return
		//}
		//
		//err = s.DB.SetUserToActive(ID)
		//if err != nil {
		//	log.Printf("Error: %v", err.Error())
		//	response.JSON(context, "", http.StatusInternalServerError, nil, []string{"Could not set user"})
		//	return
		//}
		//response.JSON(context, fmt.Sprintf("%s,your email has been verified successfully.", user.FirstName), http.StatusOK, nil, nil)
	}
}
