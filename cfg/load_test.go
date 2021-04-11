package cfg_test

import (
	. "github.com/dylannz/feature-service/cfg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cfg", func() {
	Describe("LoadYAMLDir", func() {
		It("loads all the yml files from a given directory", func() {
			cfg, err := LoadYAMLDir("./fixtures/dir")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).To(Equal(Config{
				Version: "1.0",
				Features: map[string]Feature{
					"profile_page_v2": {
						Rules: Rules{
							Enable: []EnableRule{
								{
									Fields: []string{"email", "customer_id"},
									Weight: 10,
								},
							},
							SetVars: []SetVarRule{
								{
									Fields: []string{"customer_id"},
									Weight: 50,
									Set: map[string]interface{}{
										"int_key":    1337,
										"string_key": "my_string_value",
									},
								},
							},
						},
					},
					"stripe_billing": {
						Rules: Rules{
							Enable: []EnableRule{
								{
									Field:  "customer_id",
									Values: MatchValues{Eq: []string{"123", "456"}},
									Weight: 50,
								},
								{
									Field:  "customer_id",
									Values: MatchValues{Eq: []string{"111"}},
								},
							},
							Disable: []DisableRule{
								{
									Field:  "customer_id",
									Values: MatchValues{Eq: []string{"234", "567"}},
								},
								{
									Field:  "customer_id",
									Values: MatchValues{Eq: []string{"222"}},
								},
							},
						},
					},
				},
			}))
		})
	})
})
