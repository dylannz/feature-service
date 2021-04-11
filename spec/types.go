package spec

func NewFeaturesResponse() *FeaturesResponse {
	return &FeaturesResponse{
		Features: &map[string]FeatureStatus{},
	}
}

func (r *FeaturesResponse) AddStatus(featureName string, enabled bool, vars map[string]interface{}) *FeaturesResponse {
	s := FeatureStatus{
		Enabled: &enabled,
	}

	if len(vars) > 0 {
		s.Vars = &vars
	}

	(*r.Features)[featureName] = s
	return r
}
