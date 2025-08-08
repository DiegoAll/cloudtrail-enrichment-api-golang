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
	EnrichEvent(ctx context.Context, event *models.Event) ([]*models.EnrichedEventRecord, error)
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

// Implementation of EnrichEvent for DefaultEnrichmentService
func (s *DefaultEnrichmentService) EnrichEvent(ctx context.Context, event *models.Event) ([]*models.EnrichedEventRecord, error) {
	var enrichedRecords []*models.EnrichedEventRecord

	// Iterate over each record in the input event
	for i, record := range event.Records {
		sourceIP := record.SourceIPAddress
		if sourceIP == "" {
			logger.ErrorLog.Printf("The 'sourceIPAddress' field is empty in record %d. Skipping enrichment for this record.", i)
			continue // Skip this record if IP is empty
		}
		logger.InfoLog.Printf("IP extracted from record %d: %s", i, sourceIP)

		country, err := GetCountryFromIP(sourceIP)
		if err != nil {
			logger.ErrorLog.Printf("Error getting country for IP %s (record %d): %v", sourceIP, i, err)
			return nil, fmt.Errorf("error getting country for record %d: %w", i, err)
		}
		logger.InfoLog.Printf("Country obtained for IP %s (record %d): %s", sourceIP, i, country)

		region, err := GetRegionFromCountry(country)
		if err != nil {
			logger.ErrorLog.Printf("Error getting region for country %s (record %d): %v", country, i, err)
			return nil, fmt.Errorf("error getting region for record %d: %w", i, err)
		}
		logger.InfoLog.Printf("Country: %s, Region: %s (record %d)", country, region, i)

		// Create a new instance of EnrichedEventRecord for the database
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
			Enrichment: models.EnrichmentData{ // Assign enrichment info
				Country:   country,
				Region:    region,
				Subregion: "", // Subregion is not obtained in this example.
			},
		}

		if err := s.repo.InsertLog(ctx, &enrichedRecord); err != nil {
			logger.ErrorLog.Printf("Service error inserting enriched event (record %d): %v", i, err)
			return nil, fmt.Errorf("error inserting enriched event (record %d): %w", i, err)
		}

		enrichedRecords = append(enrichedRecords, &enrichedRecord)

		logger.InfoLog.Printf("Enriched event successfully inserted (record %d). SourceIP: %s", i, sourceIP)
	}

	return enrichedRecords, nil
}

type IPInfo struct {
	Country string `json:"country"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CountryInfo []struct {
	Region    string `json:"region"`
	Subregion string `json:"subregion"`
}

func (s *DefaultEnrichmentService) Top10QueryEvents(ctx context.Context) ([]*models.EnrichedEventRecord, error) {

	records, err := s.repo.GetLatestLogs(ctx)
	if err != nil {
		logger.ErrorLog.Printf("Service error retrieving last 10 events: %v", err)
		return nil, fmt.Errorf("error retrieving last 10 events from repository: %w", err)
	}
	logger.InfoLog.Println("Service: Last 10 events retrieved successfully.")
	return records, nil
}

// Retrieves the country of an IP address using the ip-api.com API.
func GetCountryFromIP(ip string) (string, error) {
	request := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(request)
	if err != nil {
		return "", fmt.Errorf("error performing HTTP request to ip-api.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response from ip-api.com: %s", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from ip-api.com: %w", err)
	}

	var ipInfo IPInfo
	err = json.Unmarshal(bodyBytes, &ipInfo)
	if err != nil {
		return "", fmt.Errorf("error decoding response from ip-api.com: %w", err)
	}

	if ipInfo.Status != "success" {
		// Include API message if available
		errMsg := fmt.Sprintf("ip-api.com query was not successful for IP %s. Status: %s", ip, ipInfo.Status)
		if ipInfo.Message != "" {
			errMsg = fmt.Sprintf("%s, Message: %s", errMsg, ipInfo.Message)
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
		return "", fmt.Errorf("error performing HTTP request to restcountries.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// restcountries.com returns 404 if the country is not found
		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("country '%s' not found by restcountries.com", country)
		}
		return "", fmt.Errorf("unexpected response from restcountries.com: %s", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from restcountries.com: %w", err)
	}

	var countryInfo CountryInfo
	err = json.Unmarshal(bodyBytes, &countryInfo)
	if err != nil {
		return "", fmt.Errorf("error decoding response from restcountries.com: %w", err)
	}

	if len(countryInfo) > 0 && countryInfo[0].Region != "" {
		return countryInfo[0].Region, nil
	}

	return "", fmt.Errorf("region not found for country: %s or response is empty", country)
}
