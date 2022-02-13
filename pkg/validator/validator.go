package validator

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var (
	validate *validator.Validate
	once     sync.Once
)

func New() {
	if validate == nil {
		once.Do(func() {
			if validate == nil {
				validate = validator.New()
			}
		})
	}
}

// Validate struct fields
func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func IsValidNetwork(network string) bool {
	if network == "BTC" || network == "LTC" || network == "DOGE" {
		return false
	}
	return true
}
