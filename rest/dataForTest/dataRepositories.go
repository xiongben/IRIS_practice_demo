package dataForTest

import "sync"

type Query func(movie Movie) bool

type MovieRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (found bool)
	SelectMany(query Query, limit int) (results []Movie)
	InsertOrUpdate(movie Movie) (updatedMovie Movie, err error)
	Delete(query Query, limit int) (deleted bool)
}

func NewMovieRepository(source map[int64]Movie) MovieRepository {
	return &movieMemoryRepository{source: source}
}

type movieMemoryRepository struct {
	source map[int64]Movie
	mu     sync.RWMutex
}

const (
	ReadOnlyMode = iota
	ReadWriteMode
)
