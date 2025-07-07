package services

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/repository"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EnrichmentService interface {
	EnrichEvent(ctx context.Context, event *models.EnrichmentData) (models.EnrichedEventRecord, error)
	// Top10QueryEvents(ctx context.Context, numEvents int) (*models.EnrichedEventRecord, error)
	Top10QueryEvents(ctx context.Context) ([]*models.EnrichedEventRecord, error)
}

type DefaultEnrichmentService struct {
	repo repository.EnrichmentRepository
}

func NewDefaultEnrichmentService(repo repository.EnrichmentRepository) *DefaultEnrichmentService {
	return &DefaultEnrichmentService{
		repo: repo,
	}
}

func (s *DefaultEnrichmentService) EnrichEvent(ctx context.Context, event *models.Event) (any, error) {

	var enrichedRecords []models.EnrichedEventRecord

	// Iterar sobre cada record en el evento de entrada
	for i, record := range event.Records {
		sourceIP := record.SourceIPAddress
		if sourceIP == "" {
			logger.ErrorLog.Printf("El campo 'sourceIPAddress' está vacío en el registro %d. Saltando enriquecimiento para este registro.", i)
			continue // Saltamos este registro si la IP está vacía
		}
		logger.InfoLog.Printf("IP extraída del registro %d: %s", i, sourceIP)

		country, err := GetCountryFromIP(sourceIP)
		if err != nil {
			logger.ErrorLog.Printf("Error al obtener el país para la IP %s: %v", event.Records[0].SourceIPAddress, err)
			return nil, fmt.Errorf("error al obtener el país: %w", err)
		}
		logger.InfoLog.Printf("País obtenido para la IP %s: %s", event.Records[0].SourceIPAddress, country)
		fmt.Println("Country:", country)

		region, err := GetRegionFromCountry(country)
		if err != nil {
			logger.ErrorLog.Printf("Error al obtener la región para el país %s: %v", country, err)
			return nil, fmt.Errorf("error al obtener la región: %w", err)
		}
		logger.InfoLog.Printf("País: %s, Región: %s", country, region)

		// Crear una nueva instancia de EnrichedEventRecord para la base de datos
		enrichedRecord := models.EnrichedEventRecord{
			EventVersion:      record.EventVersion,
			UserIdentity:      record.UserIdentity,
			EventTime:         record.EventTime,
			EventSource:       record.EventSource,
			EventName:         record.EventName,
			AwsRegion:         record.AwsRegion,
			SourceIPAddress:   record.SourceIPAddress,
			UserAgent:         record.UserAgent,
			RequestParameters: record.RequestParameters,
			ResponseElements:  record.ResponseElements,
			Enrichment: models.EnrichmentData{ // Asignar la información de enriquecimiento
				Country:   country,
				Region:    region,
				Subregion: record.Enrichment.Subregion, // Mantener la subregión si ya existe en la entrada
			},
		}

		if err := s.repo.InsertLog(ctx, &enrichedRecord); err != nil {
			logger.ErrorLog.Printf("Error en el servicio al insertar evento enriquecido: %v", err)
			return models.EnrichedEventRecord{}, fmt.Errorf("error al insertar evento enriquecido: %w", err)
		}

		enrichedRecords = append(enrichedRecords, enrichedRecord)

		logger.InfoLog.Print("Evento enriquecido insertado exitosamente, event")

	}

	return nil, nil
}

func (s *DefaultEnrichmentService) Top10QueryEvents(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	// events, err := s.repo.GetLatestLogs(ctx)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		logger.InfoLog.Printf("Producto con UUID %s no encontrado en el servicio.")
	// 		return nil, fmt.Errorf("producto no encontrado")
	// 	}
	// 	logger.ErrorLog.Printf("Error en el servicio al obtener producto por UUID")
	// 	return nil, fmt.Errorf("error al obtener producto")
	// }
	// logger.InfoLog.Printf("Producto obtenido exitosamente con UUID")
	return nil, nil
}

type IPInfo struct {
	Country string `json:"country"`
	Status  string `json:"status"`
	Message string `json:"message"` // Añadido para capturar mensajes de error de la API
}

type CountryInfo []struct {
	Region    string `json:"region"`
	Subregion string `json:"subregion"`
}

// Event Parser  from input handler
func extractIPFromEvent(event string) string {

	return ""
}

// 205.251.233.176
// Retrieves the country of an IP address using the ip-api.com API.
func GetCountryFromIP(ip string) (string, error) {

	// request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)
	// response, err := http.Get(request)
	// if err != nil {
	// 	return models.PokeApiPokemonResponse{}, err
	// }

	request := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(request)
	if err != nil {
		return "", fmt.Errorf("error al realizar la solicitud HTTP a ip-api.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta inesperada de ip-api.com: %s", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer el cuerpo de la respuesta de ip-api.com: %w", err)
	}

	var ipInfo IPInfo
	err = json.Unmarshal(bodyBytes, &ipInfo)
	if err != nil {
		return "", fmt.Errorf("error al decodificar la respuesta de ip-api.com: %w", err)
	}

	if ipInfo.Status != "success" {
		// Incluir el mensaje de la API si está disponible
		errMsg := fmt.Sprintf("la consulta a ip-api.com no fue exitosa para IP %s. Estado: %s", ip, ipInfo.Status)
		if ipInfo.Message != "" {
			errMsg = fmt.Sprintf("%s, Mensaje: %s", errMsg, ipInfo.Message)
		}
		return "", fmt.Errorf(errMsg)
	}

	return ipInfo.Country, nil
}

// Retrieves the geographical region of a country using the restcountries.com API.
func GetRegionFromCountry(country string) (string, error) {

	request := fmt.Sprintf("https://restcountries.com/v3.1/name/%s", country)
	resp, err := http.Get(request)
	if err != nil {
		return "", fmt.Errorf("error al realizar la solicitud HTTP a restcountries.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// restcountries.com devuelve un error 404 si el país no se encuentra
		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("país '%s' no encontrado por restcountries.com", country)
		}
		return "", fmt.Errorf("respuesta inesperada de restcountries.com: %s", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer el cuerpo de la respuesta de restcountries.com: %w", err)
	}

	var countryInfo CountryInfo
	err = json.Unmarshal(bodyBytes, &countryInfo)
	if err != nil {
		return "", fmt.Errorf("error al decodificar la respuesta de restcountries.com: %w", err)
	}

	if len(countryInfo) > 0 && countryInfo[0].Region != "" {
		return countryInfo[0].Region, nil
	}

	return "", fmt.Errorf("no se encontró la región para el país: %s o la respuesta está vacía", country)
}
