package validators

import (
	"fmt"
	"reflect"
	"strconv"
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

// ParseQueryAndValidate faz o parse dos query parameters e valida em uma única operação
func (v *InputValidator) ParseQueryAndValidate(c *fiber.Ctx, s interface{}) error {
	if err := c.QueryParser(s); err != nil {
		return &ValidationError{Message: "Invalid query parameters format"}
	}

	// Aplica valores default para campos vazios
	if err := v.applyDefaults(s); err != nil {
		return &ValidationError{Message: "Error applying default values"}
	}

	// Ajusta o campo page subtraindo 1 (conversão de página baseada em 1 para índice baseado em 0)
	if err := v.adjustPageField(s); err != nil {
		return &ValidationError{Message: "Error adjusting page field"}
	}

	if err := v.ValidateStruct(s); err != nil {
		return &ValidationError{Message: v.FormatValidationError(err)}
	}

	return nil
}

// applyDefaults aplica valores default para campos que estão com valor zero
func (v *InputValidator) applyDefaults(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil
	}

	elem := val.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		// Verifica se o campo tem uma tag default
		defaultValue := fieldType.Tag.Get("default")
		if defaultValue == "" {
			continue
		}

		// Aplica o default apenas se o campo estiver com valor zero
		if field.IsZero() {
			if err := v.setFieldValue(field, defaultValue); err != nil {
				return err
			}
		}

	}

	return nil
}

// adjustPageField ajusta o campo 'page' subtraindo 1 (conversão de página baseada em 1 para índice baseado em 0)
func (v *InputValidator) adjustPageField(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil
	}

	elem := val.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		// Verifica se o campo se chama "Page" e é do tipo int64
		if fieldType.Name == "Page" && field.Kind() == reflect.Int64 {
			if field.CanSet() && field.Int() > 0 {
				field.SetInt(field.Int() - 1)
			}
			break
		}
	}

	return nil
}

// setFieldValue define o valor de um campo baseado no seu tipo
func (v *InputValidator) setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	}

	return nil
}
