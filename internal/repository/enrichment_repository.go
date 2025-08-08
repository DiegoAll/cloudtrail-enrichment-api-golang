package repository

import (
	"cloudtrail-enrichment-api-golang/models"
	"context"
)

type EnrichmentRepository interface {
	InsertLog(ctx context.Context, event *models.EnrichedEventRecord) error
	GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error)
}

// Declare a global variable for the enrichment repository instance.
var EnrichmentRepo EnrichmentRepository

// SetEnrichmentRepository allows injecting an implementation of EnrichmentRepository.
func SetEnrichmentRepository(repo EnrichmentRepository) {
	EnrichmentRepo = repo
}

// InsertLog Calls the InsertLog method of the current implementation.
func InsertLog(ctx context.Context, log *models.EnrichedEventRecord) error {
	return EnrichmentRepo.InsertLog(ctx, log)
}

// GetLatestLogs Calls the GetLatestLogs method of the current implementation.
func GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	return EnrichmentRepo.GetLatestLogs(ctx)
}
