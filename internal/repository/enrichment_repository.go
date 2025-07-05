package repository

import (
	"context"

	"github.com/sigstore/rekor/pkg/generated/models"
	//"github.com/sigstore/rekor/pkg/generated/models"
)

type EnrichmentRepository interface {
	InsertLog(ctx context.Context, log *models.LogEntry) error
	GetLatestLogs(ctx context.Context) ([]*models.LogEntry, error)
}

// Declaramos una variable global para la instancia del repositorio de enriquecimiento.
var EnrichmentRepo EnrichmentRepository

// SetEnrichmentRepository permite inyectar una implementación de EnrichmentRepository.
func SetEnrichmentRepository(repo EnrichmentRepository) {
	EnrichmentRepo = repo
}

// InsertLog es una función auxiliar que llama al método InsertLog de la implementación actual.
func InsertLog(ctx context.Context, log *models.LogEntry) error {
	return EnrichmentRepo.InsertLog(ctx, log)
}

// GetLatestLogs es una función auxiliar que llama al método GetLatestLogs de la implementación actual.
func GetLatestLogs(ctx context.Context) ([]*models.LogEntry, error) {
	return EnrichmentRepo.GetLatestLogs(ctx)
}
