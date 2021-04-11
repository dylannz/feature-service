package httpsvc_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/dylannz/feature-service/httpsvc"
	mock_httpsvc "github.com/dylannz/feature-service/httpsvc/mock"
	"github.com/dylannz/feature-service/spec"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("httpsvc", func() {
	Describe("/features/status", func() {
		It("returns the list of enabled features", func() {
			logger := logrus.WithField("httpsvc", "test")
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()
			svc := mock_httpsvc.NewMockService(ctrl)

			svc.EXPECT().
				FeaturesStatus(gomock.Any(), gomock.Any(), "").
				Return(
					spec.NewFeaturesResponse().
						AddStatus("stripe_billing", true, nil).
						AddStatus("profile_page_v2", true, nil),
					nil,
				)

			server := httptest.NewServer(NewHTTPHandler(logger, svc))
			client := server.Client()

			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/features/status",
				strings.NewReader(`
				{
					"vars": {
						"customer_id": "5671"
					}
				}
				`),
			)
			Expect(err).NotTo(HaveOccurred())
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			b, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(MatchJSON(`
				{
					"features": {
						"stripe_billing":{
							"enabled":true
						},
						"profile_page_v2":{
							"enabled":true
						}
					}
				}
			`))
		})
	})
	Describe("/features/status/{feature}", func() {
		It("returns the list of enabled features", func() {
			logger := logrus.WithField("httpsvc", "test")
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()
			svc := mock_httpsvc.NewMockService(ctrl)

			svc.EXPECT().
				FeaturesStatus(gomock.Any(), gomock.Any(), "stripe_billing").
				Return(
					spec.NewFeaturesResponse().
						AddStatus("stripe_billing", true, nil),
					nil,
				)

			server := httptest.NewServer(NewHTTPHandler(logger, svc))
			client := server.Client()

			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/features/status/stripe_billing",
				strings.NewReader(`
				{
					"vars": {
						"customer_id": "5671"
					}
				}
				`),
			)
			Expect(err).NotTo(HaveOccurred())
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			b, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(MatchJSON(`
				{
					"features": {
						"stripe_billing":{
							"enabled":true
						}
					}
				}
			`))
		})
	})
})
