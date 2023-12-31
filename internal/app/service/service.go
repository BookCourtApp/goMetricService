package service

import (
	"fmt"
	"log/slog"

	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
)

type Database interface {
	Test(metric storage.Metric) error
}

type Cache interface {
	GetSession(userId string) (string, error)
}

type Service struct {
	db     Database
	cache  Cache
	logger *slog.Logger
}

func New(logger *slog.Logger, db Database, cache Cache) *Service {
	return &Service{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (s *Service) SaveData(metric storage.Metric) error {
	const op = "service.SaveData"

	logger := s.logger.With(
		slog.String("op", op),
	)

	logger.Info("Saving metric", slog.Any("data", metric))
	logger.Debug("Getting session of user")
	sessionId, err := s.cache.GetSession(metric.UserID)
	if err != nil {
		return fmt.Errorf("%s: can't get sessionID: %s", op, sessionId)
	}

	metric.SessionID = sessionId

	logger.Debug("Inserting metric to database")
	if err := s.db.Test(metric); err != nil {
		return fmt.Errorf("%s: can't insert in database: %s", op, err.Error())
	}
	return nil
}
