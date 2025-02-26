package middleware

import (
	"campus-wall/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest(obj interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(obj); err != nil {
            c.Error(errors.NewBadRequestError("Invalid request parameters"))
            c.Abort()
            return
        }
        
        if err := validate.Struct(obj); err != nil {
            c.Error(errors.NewBadRequestError("Validation failed"))
            c.Abort()
            return
        }
        
        c.Next()
    }
}