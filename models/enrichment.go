package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnrichmentData representa la información de enriquecimiento geográfico.
type EnrichmentData struct {
	Country   string `json:"country" bson:"country"`
	Region    string `json:"region" bson:"region"`
	Subregion string `json:"subregion" bson:"subregion"`
}

// UserIdentity representa la identidad del usuario.
type UserIdentity struct {
	Type        string `json:"type" bson:"type"`
	PrincipalID string `json:"principalId" bson:"principalId"`
	Arn         string `json:"arn" bson:"arn"`
	AccessKeyID string `json:"accessKeyId" bson:"accessKeyId"`
	AccountID   string `json:"accountId" bson:"accountId"`
	UserName    string `json:"userName" bson:"userName"`
}

// InstanceItem representa un elemento de instancia.
type InstanceItem struct {
	InstanceID string `json:"instanceId" bson:"instanceId"`
}

// InstancesSet representa un conjunto de instancias.
type InstancesSet struct {
	Items []InstanceItem `json:"items" bson:"items"`
}

// RequestParameters representa los parámetros de la solicitud.
type RequestParameters struct {
	InstancesSet InstancesSet `json:"instancesSet" bson:"instancesSet"`
}

// CurrentState representa el estado actual de una instancia.
type CurrentState struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

// PreviousState representa el estado previo de una instancia.
type PreviousState struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

// ResponseInstanceItem representa un elemento de instancia en la respuesta.
type ResponseInstanceItem struct {
	InstanceID    string        `json:"instanceId" bson:"instanceId"`
	CurrentState  CurrentState  `json:"currentState" bson:"currentState"`
	PreviousState PreviousState `json:"previousState" bson:"previousState"`
}

// ResponseInstancesSet representa un conjunto de instancias en la respuesta.
type ResponseInstancesSet struct {
	Items []ResponseInstanceItem `json:"items" bson:"items"`
}

// ResponseElements representa los elementos de la respuesta.
type ResponseElements struct {
	InstancesSet ResponseInstancesSet `json:"instancesSet" bson:"instancesSet"`
}

// Event es la estructura original que define el formato de entrada de los eventos.
// Hemos modificado los tipos anidados para usar los tipos nombrados definidos arriba.
type Event struct {
	Records []struct {
		EventVersion      string            `json:"eventVersion"`
		UserIdentity      UserIdentity      `json:"userIdentity"` // Usamos el tipo nombrado
		EventTime         time.Time         `json:"eventTime"`
		EventSource       string            `json:"eventSource"`
		EventName         string            `json:"eventName"`
		AwsRegion         string            `json:"awsRegion"`
		SourceIPAddress   string            `json:"sourceIPAddress"`
		UserAgent         string            `json:"userAgent"`
		RequestParameters RequestParameters `json:"requestParameters"` // Usamos el tipo nombrado
		ResponseElements  ResponseElements  `json:"responseElements"`  // Usamos el tipo nombrado
		Enrichment        EnrichmentData    `json:"enrichment"`        // Usamos el tipo nombrado
	} `json:"Records"`
}

// EnrichedEventRecord representa un único registro de evento después de ser enriquecido,
// listo para ser insertado en la base de datos.
type EnrichedEventRecord struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // ID de MongoDB
	EventVersion      string             `json:"eventVersion" bson:"eventVersion"`
	UserIdentity      UserIdentity       `json:"userIdentity" bson:"userIdentity"`
	EventTime         time.Time          `json:"eventTime" bson:"eventTime"`
	EventSource       string             `json:"eventSource" bson:"eventSource"`
	EventName         string             `json:"eventName" bson:"eventName"`
	AwsRegion         string             `json:"awsRegion" bson:"awsRegion"`
	SourceIPAddress   string             `json:"sourceIPAddress" bson:"sourceIPAddress"`
	UserAgent         string             `json:"userAgent" bson:"userAgent"`
	RequestParameters RequestParameters  `json:"requestParameters" bson:"requestParameters"`
	ResponseElements  ResponseElements   `json:"responseElements" bson:"responseElements"`
	Enrichment        EnrichmentData     `json:"enrichment" bson:"enrichment"` // La información de enriquecimiento
}
