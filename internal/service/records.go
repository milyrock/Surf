package service

import (
	"fmt"

	"github.com/milyrock/Surf/internal/models"
)

func (s *Service) CreateRecord(rec *models.Record) (*models.Record, error) {
	err := s.recordRepo.Create(rec)
	if err != nil {
		return nil, fmt.Errorf("couldn't create record: %w", err)
	}

	return rec, nil
}
