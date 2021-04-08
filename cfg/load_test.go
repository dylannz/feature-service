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
						Rules: []Rule{
							{
								Fields: []string{"email", "customer_id"},
								Weight: 10,
							},
						},
					},
					"stripe_billing": {
						Rules: []Rule{
							{
								Field:   "customer_id",
								Include: []string{"123", "456"},
								Exclude: []string{"234", "567"},
								Weight:  50,
							},
							{
								Field:   "customer_id",
								Include: []string{"111"},
								Exclude: []string{"222"},
							},
						},
					},
				},
			}))
		})
	})
})
