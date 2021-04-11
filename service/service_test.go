package service_test

import (
	"context"

	"github.com/dylannz/feature-service/cfg"
	. "github.com/dylannz/feature-service/service"
	"github.com/dylannz/feature-service/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("service", func() {
	cfgProfile := func() cfg.Config {
		return cfg.Config{
			Version: "1.0",
			Features: map[string]cfg.Feature{
				"profile_page_v2": {
					Rules: cfg.Rules{
						Enable: []cfg.EnableRule{
							{
								Fields: []string{"email", "customer_id"},
								Weight: 10,
							},
						},
					},
				},
			},
		}
	}

	cfgStripeInclude := func() cfg.Config {
		return cfg.Config{
			Version: "1.0",
			Features: map[string]cfg.Feature{
				"stripe_billing": {
					Rules: cfg.Rules{
						Enable: []cfg.EnableRule{
							{
								Field: "customer_id",
								Values: cfg.MatchValues{
									Eq: []string{"123"},
								},
							},
						},
					},
				},
			},
		}
	}

	cfgStripeExclude := func() cfg.Config {
		return cfg.Config{
			Version: "1.0",
			Features: map[string]cfg.Feature{
				"stripe_billing": {
					Rules: cfg.Rules{
						Disable: []cfg.DisableRule{
							{
								Field: "customer_id",
								Values: cfg.MatchValues{
									Eq: []string{"321"},
								},
							},
						},
					},
				},
			},
		}
	}

	cfgStripeIncludeAndExclude := func() cfg.Config {
		return cfg.Config{
			Version: "1.0",
			Features: map[string]cfg.Feature{
				"stripe_billing": {
					Rules: cfg.Rules{
						Enable: []cfg.EnableRule{
							{
								Field: "customer_id",
								Values: cfg.MatchValues{
									Eq: []string{"123"},
								},
							},
						},
						Disable: []cfg.DisableRule{
							{
								Field: "customer_id",
								Values: cfg.MatchValues{
									Eq: []string{"123"},
								},
							},
						},
					},
				},
			},
		}
	}

	cfgStripeWeight := func() cfg.Config {
		return cfg.Config{
			Version: "1.0",
			Features: map[string]cfg.Feature{
				"stripe_billing": {
					Rules: cfg.Rules{
						Enable: []cfg.EnableRule{
							{
								Field:  "customer_id",
								Weight: 50,
							},
						},
					},
				},
			},
		}
	}

	cfgStripe := func() cfg.Config {
		c := cfgStripeInclude()
		c.Append(cfgStripeExclude())
		c.Append(cfgStripeWeight())
		return c
	}

	cfgCombined := func() cfg.Config {
		c := cfgProfile()
		c.Append(cfgStripe())
		return c
	}

	cfgSetVars := func() cfg.Config {
		return cfg.Config{
			Version: "1.0",
			Features: map[string]cfg.Feature{
				"profile_page_v2": {
					Rules: cfg.Rules{
						Enable: []cfg.EnableRule{
							{
								Fields: []string{"customer_id"},
								Values: cfg.MatchValues{Eq: []string{"123", "321"}},
							},
						},
						SetVars: []cfg.SetVarRule{
							{
								Fields: []string{"email", "customer_id"},
								Values: cfg.MatchValues{Eq: []string{"123"}},
								Set: map[string]interface{}{
									"int_key":    12345,
									"string_key": "my_string_value",
								},
							},
						},
					},
				},
			},
		}
	}

	newFeaturesRequest := func(vars map[string]interface{}) spec.FeaturesRequest {
		return spec.FeaturesRequest{
			Vars: &vars,
		}
	}

	DescribeTable(
		"FeaturesStatus",
		func(
			config cfg.Config,
			req spec.FeaturesRequest,
			featureName string,
			expectedResponse *spec.FeaturesResponse,
			expectedErrContains string,
		) {
			logger := logrus.WithField("service", "test")
			svc := NewService(logger, config)
			res, err := svc.FeaturesStatus(context.Background(), req, featureName)
			if expectedErrContains == "" {
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(expectedResponse))
			} else {
				Expect(err.Error()).To(ContainSubstring(expectedErrContains))
			}
		},
		Entry(
			"with no feature name, it returns all enabled features",
			cfgCombined(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "2"}),
			"",
			spec.NewFeaturesResponse().
				AddStatus("stripe_billing", true, nil).
				AddStatus("profile_page_v2", true, nil),
			"",
		),
		Entry(
			"when var is explicitly allowed and var is included in the request",
			cfgStripeInclude(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "123"}),
			"",
			spec.NewFeaturesResponse().AddStatus("stripe_billing", true, nil),
			"",
		),
		Entry(
			"when var is explicitly allowed and var is not included in the request",
			cfgStripeInclude(),
			newFeaturesRequest(map[string]interface{}{}),
			"",
			spec.NewFeaturesResponse(),
			"",
		),
		Entry(
			"when var is explicitly excluded and var is included in the request",
			cfgStripeInclude(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "321"}),
			"",
			spec.NewFeaturesResponse(),
			"",
		),
		Entry(
			"when var is explicitly excluded and var is not included in the request",
			cfgStripeInclude(),
			newFeaturesRequest(map[string]interface{}{}),
			"",
			spec.NewFeaturesResponse(),
			"",
		),
		Entry(
			"when the same var is explicitly excluded and included it says the feature is disabled",
			cfgStripeIncludeAndExclude(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "123"}),
			"",
			spec.NewFeaturesResponse(),
			"",
		),
		Entry(
			"when specific feature is requested it just returns that feature",
			cfgCombined(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "1"}),
			"stripe_billing",
			spec.NewFeaturesResponse().AddStatus("stripe_billing", true, nil),
			"",
		),
		Entry(
			"when invalid feature is requested it returns an error",
			cfg.Config{},
			newFeaturesRequest(map[string]interface{}{"customer_id": "1"}),
			"stripe_billing",
			nil,
			"unknown feature: 'stripe_billing'",
		),
		Entry(
			"vars are returned when they have been configured",
			cfgSetVars(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "123"}),
			"profile_page_v2",
			spec.NewFeaturesResponse().
				AddStatus(
					"profile_page_v2",
					true,
					map[string]interface{}{
						"int_key":    12345,
						"string_key": "my_string_value",
					},
				),
			"",
		),
		Entry(
			"vars are not returned when they don't meet the set_vars rules",
			cfgSetVars(),
			newFeaturesRequest(map[string]interface{}{"customer_id": "321"}),
			"profile_page_v2",
			spec.NewFeaturesResponse().AddStatus("profile_page_v2", true, nil),
			"",
		),
	)
})
