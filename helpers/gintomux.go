package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Convert Gin Handler to Gorilla Mux Handler
func ConvertGinToMux(ginHandler gin.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a gin context manually
		c, _ := gin.CreateTestContext(w)
		c.Request = r

		// Call the original gin handler with the new context
		ginHandler(c)
	}
}
