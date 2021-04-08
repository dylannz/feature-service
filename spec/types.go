package spec

func NewFeaturesResponse() *FeaturesResponse {
	return &FeaturesResponse{}
}

func (r *FeaturesResponse) AddStatus(featureName string, enabled bool) *FeaturesResponse {
	if r.Features == nil {
		r.Features = &map[string]FeatureStatus{}
	}

	(*r.Features)[featureName] = FeatureStatus{
		Enabled: &enabled,
	}

	return r
}
