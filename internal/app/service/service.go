package service

import (
	"fmt"

	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
)

type Database interface {
	Test(metric storage.Metric) error
}

type Service struct {
	db Database
}

func New(db Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) SaveData(metric storage.Metric) error {
	const op = "service.SaveData"

	if err := s.db.Test(metric); err != nil {
		return fmt.Errorf("%s: %s", op, err.Error())
	}
	return nil
}
