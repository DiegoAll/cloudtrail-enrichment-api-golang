package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnrichmentData represents geographic enrichment information.
type EnrichmentData struct {
	Country   string `json:"country" bson:"country"`
	Region    string `json:"region" bson:"region"`
	Subregion string `json:"subregion" bson:"subregion"`
}

// UserIdentity represents the user's identity.
type UserIdentity struct {
	Type        string `json:"type" bson:"type"`
	PrincipalID string `json:"principalId" bson:"principalId"`
	Arn         string `json:"arn" bson:"arn"`
	AccessKeyID string `json:"accessKeyId" bson:"accessKeyId"`
	AccountID   string `json:"accountId" bson:"accountId"`
	UserName    string `json:"userName" bson:"userName"`
}

// InstanceItem represents an instance item.
type InstanceItem struct {
	InstanceID string `json:"instanceId" bson:"instanceId"`
}

// InstancesSet represents a set of instances.
type InstancesSet struct {
	Items []InstanceItem `json:"items" bson:"items"`
}

// RequestParameters represents the request parameters.
type RequestParameters struct {
	InstancesSet InstancesSet `json:"instancesSet" bson:"instancesSet"`
}

// CurrentState represents the current state of an instance.
type CurrentState struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

// PreviousState represents the previous state of an instance.
type PreviousState struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

// ResponseInstanceItem represents an instance item in the response.
type ResponseInstanceItem struct {
	InstanceID    string        `json:"instanceId" bson:"instanceId"`
	CurrentState  CurrentState  `json:"currentState" bson:"currentState"`
	PreviousState PreviousState `json:"previousState" bson:"previousState"`
}

// ResponseInstancesSet represents a set of instances in the response.
type ResponseInstancesSet struct {
	Items []ResponseInstanceItem `json:"items" bson:"items"`
}

// ResponseElements represents the response elements.
type ResponseElements struct {
	InstancesSet ResponseInstancesSet `json:"instancesSet" bson:"instancesSet"`
}

// Event is the original structure defining the input event format.
type Event struct {
	Records []struct {
		EventVersion      string            `json:"eventVersion"`
		UserIdentity      UserIdentity      `json:"userIdentity"`
		EventTime         time.Time         `json:"eventTime"`
		EventSource       string            `json:"eventSource"`
		EventName         string            `json:"eventName"`
		AwsRegion         string            `json:"awsRegion"`
		SourceIPAddress   string            `json:"sourceIPAddress"`
		UserAgent         string            `json:"userAgent"`
		RequestParameters RequestParameters `json:"requestParameters"`
		ResponseElements  ResponseElements  `json:"responseElements"`
		Enrichment        EnrichmentData    `json:"enrichment"`
	} `json:"Records"`
}

// EnrichedEventRecord represents a single event record after being enriched, ready to be inserted into the database.
type EnrichedEventRecord struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // MongoDB ID
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
	Enrichment        EnrichmentData     `json:"enrichment" bson:"enrichment"` // The enrichment information
}
