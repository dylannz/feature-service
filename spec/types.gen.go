// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package spec

// FeatureStatus defines model for FeatureStatus.
type FeatureStatus struct {
	Enabled *bool                   `json:"enabled,omitempty"`
	Vars    *map[string]interface{} `json:"vars,omitempty"`
}

// FeaturesRequest defines model for FeaturesRequest.
type FeaturesRequest struct {
	Vars *map[string]interface{} `json:"vars,omitempty"`
}

// FeaturesResponse defines model for FeaturesResponse.
type FeaturesResponse struct {
	Features *map[string]FeatureStatus `json:"features,omitempty"`
}

// PostFeaturesStatusJSONBody defines parameters for PostFeaturesStatus.
type PostFeaturesStatusJSONBody FeaturesRequest

// PostFeaturesStatusFeatureJSONBody defines parameters for PostFeaturesStatusFeature.
type PostFeaturesStatusFeatureJSONBody FeaturesRequest

// PostFeaturesStatusJSONRequestBody defines body for PostFeaturesStatus for application/json ContentType.
type PostFeaturesStatusJSONRequestBody PostFeaturesStatusJSONBody

// PostFeaturesStatusFeatureJSONRequestBody defines body for PostFeaturesStatusFeature for application/json ContentType.
type PostFeaturesStatusFeatureJSONRequestBody PostFeaturesStatusFeatureJSONBody
