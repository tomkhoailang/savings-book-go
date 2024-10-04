package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate = validator.New()
func ReadRequest(ctx *gin.Context, request interface{}) error {
	if err:= ctx.ShouldBind(request); err != nil {
		return err
	}
	if err := validate.StructCtx(ctx.Request.Context(),request); err != nil {
		return err
	}
	return nil
}
func GetUserId(c *gin.Context) (string, error){
	id, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return "", errors.New("user ID not found")
	}
	userId, ok := id.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error in casting id to string"})
		return "", errors.New("Invalid user Id")
	}
	return userId, nil
}

func ReadRequestGeneric[TInput any](c *gin.Context, request *TInput) []ValidationError {
	if err := c.ShouldBind(&request); err != nil {
		return []ValidationError{{Message: err.Error()}}
	}
	if err := validate.StructCtx(c.Request.Context(),request); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			errors := make([]ValidationError, len(validationErrs))

			for i, ve := range validationErrs {
				errors[i] = ValidationError{
					Field:   ve.Field(),
					Message: ve.Tag() + " validation failed: " + ve.Param(),
				}
			}
			return errors
		}
		return []ValidationError{{Message: err.Error()}}
	}
	return nil
}

func HandleCreateRequest[TInput any, TOutput any, TEntity any](createFunc func(ctx context.Context, input *TInput, userId string) (*TEntity, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input TInput

		errors := ReadRequestGeneric[TInput](c, &input)
		if errors != nil {
			c.JSON(http.StatusBadRequest, errors)
			return
		}
		userId, err := GetUserId(c)
		if err != nil {
			return
		}
		entity, err := createFunc(c.Request.Context(), &input, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var output TOutput
		err = copier.Copy(&output, entity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, &output)
		return

	}
}



