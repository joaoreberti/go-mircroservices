package memory

import (
	"context"

	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

// Repository defines a rating repository.
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{
		model.RecordTypeMovie: {
			"0": {
				{
					RecordID:   "tt0111161",
					RecordType: model.RecordTypeMovie,
					UserID:     "user1",
					Value:      5,
				},
				{
					RecordID:   "tt0111161",
					RecordType: model.RecordTypeMovie,
					UserID:     "user2",
					Value:      3,
				},
			},
		},
		model.RecordTypeActor: {
			"1": {
				{
					RecordID:   "nm0000151",
					RecordType: model.RecordTypeActor,
					UserID:     "user1",
					Value:      5,
				},
			},
		},
	}}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}
	r.data[recordType][recordID] =
		append(r.data[recordType][recordID], *rating)
	return nil
}
