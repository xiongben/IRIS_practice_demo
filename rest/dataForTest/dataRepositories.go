package dataForTest

import (
	"errors"
	"sync"
)

type Query func(movie Movie) bool

type MovieRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (movie Movie, found bool)
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

func (r *movieMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0
	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}
	for _, movie := range r.source {
		ok = query(movie)
		if ok {
			if action(movie) {
				loops++
				if actionLimit >= loops {
					break
				}
			}
		}
	}
	return
}

func (r *movieMemoryRepository) Select(query Query) (movie Movie, found bool) {
	found = r.Exec(query, func(m Movie) bool {
		movie = m
		return true
	}, 1, ReadOnlyMode)
	return
}

func (r *movieMemoryRepository) SelectMany(query Query, limit int) (results []Movie) {
	r.Exec(query, func(m Movie) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)
	return
}

func (r *movieMemoryRepository) InsertOrUpdate(movie Movie) (updatedMovie Movie, err error) {
	id := movie.ID
	if id == 0 {
		var lastID int64
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()
		id = lastID + 1
		movie.ID = id
		r.mu.Lock()
		r.source[id] = movie
		r.mu.Unlock()
		return movie, nil
	}
	current, exists := r.Select(func(m Movie) bool {
		return m.ID == id
	})
	if !exists {
		return Movie{}, errors.New("failed to update a nonexistence movie")
	}
	if movie.Poster != "" {
		current.Poster = movie.Poster
	}
	if movie.Genre != "" {
		current.Genre = movie.Genre
	}
	r.mu.Lock()
	r.source[id] = current
	r.mu.Unlock()
	return movie, nil
}

func (r *movieMemoryRepository) Delete(query Query, limit int) (deleted bool) {
	return r.Exec(query, func(m Movie) bool {
		delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)

}
