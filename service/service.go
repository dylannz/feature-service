package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"sort"

	"github.com/dylannz/feature-service/cfg"
	"github.com/dylannz/feature-service/reqcontext"
	"github.com/dylannz/feature-service/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger logrus.FieldLogger
	config cfg.Config

	featureList []string
}

func NewService(logger logrus.FieldLogger, config cfg.Config) *Service {
	// TODO: build cache here
	svc := &Service{
		logger: logger,
		config: config,

		featureList: make([]string, 0, len(config.Features)),
	}

	for feature := range config.Features {
		svc.featureList = append(svc.featureList, feature)
	}
	sort.StringSlice(svc.featureList).Sort()

	return svc
}

func (s Service) FeaturesStatus(ctx context.Context, req spec.FeaturesRequest, featureName string) (*spec.FeaturesResponse, error) {
	if featureName != "" {
		return s.featureStatus(ctx, req, featureName)
	}

	res := spec.NewFeaturesResponse()
	for _, fn := range s.featureList {
		r, err := s.featureStatus(ctx, req, fn)
		if err != nil {
			return res, err
		}

		if r.Features != nil {
			for k, v := range *r.Features {
				res.AddStatus(k, *v.Enabled)
			}
		}
	}

	return res, nil
}

func (s Service) featureStatus(ctx context.Context, req spec.FeaturesRequest, featureName string) (*spec.FeaturesResponse, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"request_id": reqcontext.RequestIDFromContext(ctx),
	})
	res := spec.NewFeaturesResponse()

	feature, ok := s.config.Features[featureName]
	if !ok {
		return res, errors.Errorf("unknown feature: '%s'", featureName)
	}

	vars := map[string]string{}
	if req.Vars != nil {
		for k, v := range *req.Vars {
			switch t := v.(type) {
			case string:
				vars[k] = t
			default:
				vars[k] = fmt.Sprint(v)
			}
		}
	}

	// first we deal with all exclude rules
	if foreachField(feature.Rules, func(field string, rule cfg.Rule) bool {
		t := varInSlice(rule.Exclude, field, vars)
		logger.Debugf("check: field '%s' in %#v matches exclude rules %#v: %t", field, req.Vars, rule.Exclude, t)
		return t
	}) {
		logger.Debug("match: matched exclude rule")
		return res, nil
	}

	// now we deal with explicit includes
	if foreachField(feature.Rules, func(field string, rule cfg.Rule) bool {
		t := varInSlice(rule.Include, field, vars)
		logger.Debugf("check: field '%s' in %#v matches include rules %#v: %t", field, req.Vars, rule.Exclude, t)
		return t
	}) {
		logger.Debug("match: matched include rule")
		res.AddStatus(featureName, true)
		return res, nil
	}

	if ruleWeight(logger, feature.Rules, vars) {
		logger.Debug("match: matched weight rule")
		res.AddStatus(featureName, true)
		return res, nil
	}

	return res, nil
}

func ruleFields(rule cfg.Rule) []string {
	if len(rule.Fields) == 0 {
		return []string{rule.Field}
	}
	return rule.Fields
}

func ruleWeight(logger logrus.FieldLogger, rules []cfg.Rule, vars map[string]string) bool {
	var rule cfg.Rule
	logger.Debugf("looking in %d rules for a rule with weight > 0...", len(rules))
	// find the first rule with weight > 0, ignore all others
	for _, r := range rules {
		if r.Weight > 0 {
			rule = r
			logger.Debugf("found a rule with weight > 0: %#v", rule)
			break
		}
	}

	// if no rules had weight > 0, we exit early
	if rule.Weight == 0 {
		logger.Debug("found no rules with weight > 0")
		return false
	}

	// first build a string containing all the key/value pairs
	b := bytes.Buffer{}
	fields := ruleFields(rule)
	logger.Debugf("using keys/values from fields: %#v", fields)
	for _, field := range fields {
		b.WriteString(field)
		b.WriteString("=")
		if s, ok := vars[field]; ok {
			b.WriteString(s)
		}
		b.WriteString(";")
	}

	// Hash as md5, convert the first half of the hash to a number,
	// then modulo 100 and see if it's less than the defined weight.
	// And there we have it - a deterministic way to calculate whether
	// a feature should be enabled based on an an arbitrary list of
	// key/value pairs.
	h := md5.Sum(b.Bytes())
	var n uint64
	for i := 0; i < 8; i++ {
		n <<= 8
		n |= uint64(uint8(h[i]))
	}
	c := int(n % 100)
	t := c < rule.Weight
	logger.Debugf("check: hash result < weight (%d < %d): %t", c, rule.Weight, t)
	return t
}

func foreachField(rules []cfg.Rule, fn func(string, cfg.Rule) bool) bool {
	for _, rule := range rules {
		for _, field := range ruleFields(rule) {
			if fn(field, rule) {
				return true
			}
		}
	}
	return false
}

func varInSlice(a []string, field string, vars map[string]string) bool {
	s, ok := vars[field]
	if ok {
		for _, v := range a {
			if v == s {
				return true
			}
		}
	}
	return false
}
