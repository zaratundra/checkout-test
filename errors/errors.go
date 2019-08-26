package errors

import (
	"fmt"
)

type ProductNotFound struct {
	Code string
}

type PromotionNotFound struct {
	Code string
}

type PromotionInvalid struct {
	Code string
	Msg  string
}

type BasketNotFound struct {
	Id string
}

type PrimaryKeyError struct {
	Id string
}

type ValidationError struct {
	Errors []*ValidationErrorDescription
}

type ValidationErrorDescription struct {
	Field   string
	Message string
}

func NewProductNotFound(code string) *ProductNotFound {
	return &ProductNotFound{Code: code}
}

func NewPromotionNotFound(code string) *PromotionNotFound {
	return &PromotionNotFound{
		Code: code,
	}
}

func NewPromotionInvalid(code, message string) *PromotionInvalid {
	return &PromotionInvalid{
		Code: code,
		Msg:  message,
	}
}

func NewBasketNotFound(id string) *BasketNotFound {
	return &BasketNotFound{Id: id}
}

func NewPrimaryKeyError(id string) *PrimaryKeyError {
	return &PrimaryKeyError{Id: id}
}

func NewValidationError(errors []*ValidationErrorDescription) *ValidationError {
	return &ValidationError{
		Errors: errors,
	}
}

func NewValidationErrorDescription(field, message string) *ValidationErrorDescription {
	return &ValidationErrorDescription{
		Field:   field,
		Message: message,
	}
}

//TODO: localization for error messages
func (p *ProductNotFound) Error() string {
	return fmt.Sprintf("Product %v not found", p.Code)
}

func (p *PromotionNotFound) Error() string {
	return fmt.Sprintf("Promotion %v not found", p.Code)
}

func (b *BasketNotFound) Error() string {
	return fmt.Sprintf("Basket %v not found", b.Id)
}

func (p *PromotionInvalid) Error() string {
	return fmt.Sprintf("Promotion %v invalid: %v", p.Code, p.Msg)
}

func (p *PrimaryKeyError) Error() string {
	return fmt.Sprintf("Primary key already exists: %v", p.Id)
}

func (e *ValidationError) Error() string {
	return fmt.Sprint("There has been a validation error")
}
