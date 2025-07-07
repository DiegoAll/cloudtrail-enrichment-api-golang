package repository

import (
	"cloudtrail-enrichment-api-golang/models"
	"context"
)

type EnrichmentRepository interface {
	InsertLog(ctx context.Context, event *models.EnrichedEventRecord) error
	GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error)
}

// Declaramos una variable global para la instancia del repositorio de enriquecimiento.
var EnrichmentRepo EnrichmentRepository

// SetEnrichmentRepository permite inyectar una implementación de EnrichmentRepository.
func SetEnrichmentRepository(repo EnrichmentRepository) {
	EnrichmentRepo = repo
}

// InsertLog es una función auxiliar que llama al método InsertLog de la implementación actual.
func InsertLog(ctx context.Context, log *models.EnrichedEventRecord) error {
	return EnrichmentRepo.InsertLog(ctx, log)
}

// GetLatestLogs es una función auxiliar que llama al método GetLatestLogs de la implementación actual.
// POR QUE EL ERROR ESTA ADENTRO?
func GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	return EnrichmentRepo.GetLatestLogs(ctx)
}
