package unit

import (
	"testing"

	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	validation.Init()

	t.Run("Valid CreateItemRequest", func(t *testing.T) {
		req := models.CreateItemRequest{
			Name:  "Test Item",
			Email: "test@example.com",
		}

		errs := validation.ValidateStruct(req)
		assert.Empty(t, errs)
	})

	t.Run("Invalid CreateItemRequest - short name", func(t *testing.T) {
		req := models.CreateItemRequest{
			Name:  "", // Empty name
			Email: "invalid-email",
		}

		errs := validation.ValidateStruct(req)
		assert.Len(t, errs, 2)
		assert.Equal(t, "name", errs[0].Field)
		assert.Equal(t, "email", errs[1].Field)
	})

	t.Run("Valid User", func(t *testing.T) {
		user := models.User{
			Username: "john_doe123",
			Password: "StrongPass123!",
			Email:    "john@example.com",
		}

		errs := validation.ValidateStruct(user)
		assert.Empty(t, errs)
	})

	t.Run("Invalid User - weak password", func(t *testing.T) {
		user := models.User{
			Username: "john_doe123",
			Password: "weak", // Too weak
			Email:    "invalid-email",
		}

		errs := validation.ValidateStruct(user)
		assert.Len(t, errs, 2)
	})
}

func TestCustomValidations(t *testing.T) {
	validation.Init()

	testCases := []struct {
		name      string
		username  string
		password  string
		shouldErr bool
	}{
		{"Valid credentials", "user123", "StrongPass1!", false},
		{"Username too short", "ab", "StrongPass1!", true},
		{"Username with spaces", "user name", "StrongPass1!", true},
		{"Weak password no upper", "user123", "weakpass1!", true},
		{"Weak password no number", "user123", "WeakPass!!", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := models.User{
				Username: tc.username,
				Password: tc.password,
				Email:    "test@example.com",
			}

			errs := validation.ValidateStruct(user)
			if tc.shouldErr {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}