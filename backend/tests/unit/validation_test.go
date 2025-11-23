package unit

import (
	"testing"

	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	validation.Init()

	t.Run("Valid Web Project", func(t *testing.T) {
		project := models.CreateWebProjectRequest{
			Name:        "Valid Web Project Name Here",
			Description: "This is a valid description that meets the minimum length requirement of 20 characters.",
			Img:         "https://example.com/image.jpg",
			Price:       1500.00,
			TimeDevelop: 30,
		}

		errs := validation.ValidateStruct(project)
		assert.Empty(t, errs)
	})

	t.Run("Valid Staff Member", func(t *testing.T) {
		staff := models.CreateStaffRequest{
			Name:        "Valid Staff Member Full Name",
			Description: "This is a valid description that meets the minimum length requirement of 20 characters for staff members.",
			Img:         "https://example.com/staff.jpg",
			Role:        "Senior Developer",
		}

		errs := validation.ValidateStruct(staff)
		assert.Empty(t, errs)
	})

	t.Run("Invalid Staff Member - Short Name", func(t *testing.T) {
		staff := models.CreateStaffRequest{
			Name:        "Short",
			Description: "Valid description that meets requirements",
			Img:         "https://example.com/staff.jpg",
			Role:        "Developer",
		}

		errs := validation.ValidateStruct(staff)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "name", errs[0].Field)
	})

	t.Run("Invalid Staff Member - Short Description", func(t *testing.T) {
		staff := models.CreateStaffRequest{
			Name:        "Valid Staff Member Full Name",
			Description: "Too short",
			Img:         "https://example.com/staff.jpg",
			Role:        "Developer",
		}

		errs := validation.ValidateStruct(staff)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "description", errs[0].Field)
	})

	t.Run("Invalid Staff Member - Invalid URL", func(t *testing.T) {
		staff := models.CreateStaffRequest{
			Name:        "Valid Staff Member Full Name",
			Description: "Valid description that meets requirements",
			Img:         "invalid-url",
			Role:        "Developer",
		}

		errs := validation.ValidateStruct(staff)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "img", errs[0].Field)
	})

	t.Run("Invalid Staff Member - Empty Role", func(t *testing.T) {
		staff := models.CreateStaffRequest{
			Name:        "Valid Staff Member Full Name",
			Description: "Valid description that meets requirements",
			Img:         "https://example.com/staff.jpg",
			Role:        "",
		}

		errs := validation.ValidateStruct(staff)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "role", errs[0].Field)
	})

	t.Run("Invalid Web Project - Short Name", func(t *testing.T) {
		project := models.CreateWebProjectRequest{
			Name:        "Short",
			Description: "Valid description that meets requirements",
			Img:         "https://example.com/image.jpg",
			Price:       1500.00,
			TimeDevelop: 30,
		}

		errs := validation.ValidateStruct(project)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "name", errs[0].Field)
	})

	t.Run("Invalid Web Project - Invalid URL", func(t *testing.T) {
		project := models.CreateWebProjectRequest{
			Name:        "Valid Web Project Name Here",
			Description: "Valid description that meets requirements",
			Img:         "invalid-url",
			Price:       1500.00,
			TimeDevelop: 30,
		}

		errs := validation.ValidateStruct(project)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "img", errs[0].Field)
	})

	t.Run("Invalid Web Project - Negative Price", func(t *testing.T) {
		project := models.CreateWebProjectRequest{
			Name:        "Valid Web Project Name Here",
			Description: "Valid description that meets requirements",
			Img:         "https://example.com/image.jpg",
			Price:       -100.00,
			TimeDevelop: 30,
		}

		errs := validation.ValidateStruct(project)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "price", errs[0].Field)
	})

	t.Run("Invalid Web Project - Time Too Long", func(t *testing.T) {
		project := models.CreateWebProjectRequest{
			Name:        "Valid Web Project Name Here",
			Description: "Valid description that meets requirements",
			Img:         "https://example.com/image.jpg",
			Price:       1500.00,
			TimeDevelop: 2000,
		}

		errs := validation.ValidateStruct(project)
		assert.NotEmpty(t, errs)
		assert.Equal(t, "timedevelop", errs[0].Field)
	})
}