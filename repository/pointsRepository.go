package repository

// PointRepository represents a repository for storing receipt IDs and points.
type PointRepository struct {
	data map[string]int
}

var (
	// pointRepository holds the single instance of the PointRepository.
	pointRepository *PointRepository
)

// GetPointRepository returns the single instance of the PointRepository.
func GetPointRepository() *PointRepository {
	if pointRepository == nil {
		pointRepository = &PointRepository{
			data: make(map[string]int),
		}
	}
	return pointRepository
}

// Adds the ID and corresponding points to the data store.
func (repo *PointRepository) Put(id string, points int) {
	repo.data[id] = points
}

// Get retrieves the points associated with the given ID from the repository.
func (repo *PointRepository) Get(id string) (int, bool) {
	points, ok := repo.data[id]
	return points, ok
}
