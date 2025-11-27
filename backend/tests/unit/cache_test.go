package unit

import (
	"testing"
	"time"

	"ASMO-site-backend/internal/models"
	testutils "ASMO-site-backend/tests/testutils"

	"github.com/stretchr/testify/assert"
)

func TestRedisMock(t *testing.T) {
	mock := testutils.NewRedisMock()
	defer mock.Close()

	// Test Set and Get
	project := models.WebProjects{
		ID:          1,
		Name:        "Test Project",
		Description: "Test Description",
		Img:         "https://example.com/image.jpg",
		Price:       1000.00,
		TimeDevelop: 30,
	}

	err := mock.Set("test:project", project, 5*time.Minute)
	assert.NoError(t, err)

	var retrieved models.WebProjects
	err = mock.Get("test:project", &retrieved)
	assert.NoError(t, err)
	assert.Equal(t, project.Name, retrieved.Name)
	assert.Equal(t, project.Price, retrieved.Price)

	// Test Get non-existent key
	err = mock.Get("test:nonexistent", &retrieved)
	assert.Error(t, err)
	assert.Equal(t, testutils.ErrNotFound, err)

	// Test Delete
	err = mock.Delete("test:project")
	assert.NoError(t, err)

	err = mock.Get("test:project", &retrieved)
	assert.Error(t, err)
	assert.Equal(t, testutils.ErrNotFound, err)
}

func TestRedisMockWithSlice(t *testing.T) {
	mock := testutils.NewRedisMock()
	defer mock.Close()

	projects := []models.WebProjects{
		{
			ID:          1,
			Name:        "Project 1",
			Description: "Description 1",
			Img:         "https://example.com/1.jpg",
			Price:       1000.00,
			TimeDevelop: 30,
		},
		{
			ID:          2,
			Name:        "Project 2",
			Description: "Description 2",
			Img:         "https://example.com/2.jpg",
			Price:       2000.00,
			TimeDevelop: 45,
		},
	}

	err := mock.Set("web_projects:all", projects, 5*time.Minute)
	assert.NoError(t, err)

	var retrieved []models.WebProjects
	err = mock.Get("web_projects:all", &retrieved)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, "Project 1", retrieved[0].Name)
	assert.Equal(t, "Project 2", retrieved[1].Name)
}