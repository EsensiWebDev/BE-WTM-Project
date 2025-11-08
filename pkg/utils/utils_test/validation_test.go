package utils_test

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"testing"

	"github.com/stretchr/testify/assert"
	"wtm-backend/pkg/utils"
)

func TestNotEmptyAfterTrim(t *testing.T) {
	rule := utils.NotEmptyAfterTrim("username")

	t.Run("valid string", func(t *testing.T) {
		err := rule.Validate("john_doe")
		assert.NoError(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		err := rule.Validate("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username cannot be blank")
	})

	t.Run("whitespace only", func(t *testing.T) {
		err := rule.Validate("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username cannot be blank")
	})
}

func TestParseValidationErrors(t *testing.T) {
	t.Run("should return error map when validation.Errors is present", func(t *testing.T) {
		err := validation.Errors{
			"username": errors.New("cannot be blank"),
			"email":    errors.New("must be a valid email"),
		}

		result := utils.ParseValidationErrors(err)

		expected := map[string]string{
			"username": "cannot be blank",
			"email":    "must be a valid email",
		}

		assert.Equal(t, expected, result)
	})

	t.Run("should return nil when error is not validation.Errors", func(t *testing.T) {
		err := errors.New("some other error")

		result := utils.ParseValidationErrors(err)

		assert.Nil(t, result)
	})

	t.Run("should return nil when error is nil", func(t *testing.T) {
		result := utils.ParseValidationErrors(nil)

		assert.Nil(t, result)
	})
}
