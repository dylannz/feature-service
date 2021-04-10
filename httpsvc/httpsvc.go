package httpsvc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dylannz/feature-service/reqcontext"
	"github.com/dylannz/feature-service/spec"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type HTTPService struct {
	logger  logrus.FieldLogger
	service Service
}

//go:generate mockgen -source=httpsvc.go -destination=mock/httpsvc.go
type Service interface {
	FeaturesStatus(ctx context.Context, req spec.FeaturesRequest, feature string) (*spec.FeaturesResponse, error)
}

func NewHTTPHandler(logger logrus.FieldLogger, service Service) http.Handler {
	svc := HTTPService{
		logger:  logger,
		service: service,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.MethodFunc(http.MethodGet, "/health", svc.handleHealth)

	return spec.HandlerFromMux(svc, r)
}

func (s HTTPService) handleHealth(w http.ResponseWriter, r *http.Request) {
	// just return 200
}

func (s HTTPService) PostFeaturesStatus(w http.ResponseWriter, r *http.Request) {
	s.PostFeaturesStatusFeature(w, r, "")
}

func (s HTTPService) PostFeaturesStatusFeature(w http.ResponseWriter, r *http.Request, feature string) {
	var req spec.FeaturesRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "decode request body"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ctx = reqcontext.ContextWithRequestID(ctx, r.Header.Get("x-request-id"))
	res, err := s.service.FeaturesStatus(ctx, req, feature)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "service"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(res)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "encode response body"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}
