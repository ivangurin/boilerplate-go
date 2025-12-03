package middleware

import "github.com/gin-gonic/gin"

func (m *middleware) Auth() gin.HandlerFunc {
	return m.authService.Validate()
}
