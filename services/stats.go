package services

import (
	"book-tracker/models"
	"book-tracker/store"
	"context"
	"fmt"
)

type StatsService interface {
	GetStats(ctx context.Context) (models.Stats, error)
}

type statsService struct {
	store store.StatsStore
}

func NewStatsService(store store.StatsStore) StatsService {
	return &statsService{store: store}
}

func (s *statsService) GetStats(ctx context.Context) (models.Stats, error) {
	// NOTE: Could already be handled at store level - lets see if we clean-up
	totalRead, readingProgress, popularAuthor, err := s.store.GetStats(ctx)

	if err != nil {
		return models.Stats{}, fmt.Errorf("get stats: %w", err)
	}

	return models.Stats{
		TotalRead:       totalRead,
		ReadingProgress: readingProgress,
		PopularAuthor:   popularAuthor,
	}, nil
}
