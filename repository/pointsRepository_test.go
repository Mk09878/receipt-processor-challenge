package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointRepository(t *testing.T) {
	// Ensure pointRepository is nil initially
	assert.Nil(t, pointRepository)

	// First call to GetPointRepository should initialize the repository
	repo1 := GetPointRepository()
	assert.NotNil(t, pointRepository, "Expected pointRepository to be initialized after the first call")

	// Second call to GetPointRepository should return the same instance
	repo2 := GetPointRepository()
	assert.Equal(t, repo1, repo2, "Expected the same instance to be returned on subsequent calls")

	// Put some data into the repository
	repo2.Put("ID1", 100)
	repo2.Put("ID2", 200)

	// Test retrieving data from the repository
	points, found := repo2.Get("ID1")
	assert.True(t, found, "Expected to find points for ID1")
	assert.Equal(t, 100, points, "Expected points for ID1 to be 100")

	points, found = repo2.Get("ID2")
	assert.True(t, found, "Expected to find points for ID2")
	assert.Equal(t, 200, points, "Expected points for ID2 to be 200")

	// Test retrieving non-existent data
	points, found = repo2.Get("ID4")
	assert.False(t, found, "Expected not to find points for ID4")
	assert.Equal(t, 0, points, "Expected points for ID4 to be 0")
}
