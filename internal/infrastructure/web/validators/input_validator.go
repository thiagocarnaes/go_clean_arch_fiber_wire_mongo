package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InputValidator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewInputValidator() *InputValidator {
	return &InputValidator{
		validator: validator.New(),
	}
}

// ValidateStruct valida uma struct e retorna erro formatado se houver problemas
func (v *InputValidator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

// FormatValidationError formata os erros de validação em mensagens legíveis
func (v *InputValidator) FormatValidationError(err error) string {
	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("Field '%s' is required", err.Field()))
		case "min":
			messages = append(messages, fmt.Sprintf("Field '%s' must be at least %s characters", err.Field(), err.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("Field '%s' must be at most %s characters", err.Field(), err.Param()))
		case "email":
			messages = append(messages, fmt.Sprintf("Field '%s' must be a valid email", err.Field()))
		case "len":
			messages = append(messages, fmt.Sprintf("Field '%s' must be exactly %s characters", err.Field(), err.Param()))
		case "numeric":
			messages = append(messages, fmt.Sprintf("Field '%s' must be numeric", err.Field()))
		case "alpha":
			messages = append(messages, fmt.Sprintf("Field '%s' must contain only letters", err.Field()))
		case "alphanum":
			messages = append(messages, fmt.Sprintf("Field '%s' must contain only letters and numbers", err.Field()))
		default:
			messages = append(messages, fmt.Sprintf("Field '%s' is invalid", err.Field()))
		}
	}
	return strings.Join(messages, ", ")
}

// ParseAndValidate faz o parse do body e valida em uma única operação
func (v *InputValidator) ParseAndValidate(c *fiber.Ctx, s interface{}) error {
	if err := c.BodyParser(s); err != nil {
		return &ValidationError{Message: "Invalid JSON format"}
	}
	
	if err := v.ValidateStruct(s); err != nil {
		return &ValidationError{Message: v.FormatValidationError(err)}
	}
	
	return nil
}
