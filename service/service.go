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
				(*res.Features)[k] = v
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

	// first we deal with disable rules
	if foreachDisableField(feature.Rules.Disable, func(field string, rule cfg.DisableRule) bool {
		t := varInSlice(rule.Values.Eq, field, vars)
		logger.Debugf("check: field '%s' in %#v matches disable rules %#v: %t", field, req.Vars, rule.Values.Eq, t)
		return t
	}) {
		logger.Debug("match: matched values.eq rule")
		return res, nil
	}

	// now we deal with enable rules
	if foreachEnableField(feature.Rules.Enable, func(field string, rule cfg.EnableRule) bool {
		t := varInSlice(rule.Values.Eq, field, vars)
		logger.Debugf("check: field '%s' in %#v matches enable rules %#v: %t", field, req.Vars, rule.Values.Eq, t)
		return t
	}) {
		logger.Debug("match: matched values.eq rule")
		res.AddStatus(featureName, true, setVars(logger, feature.Rules.SetVars, vars))
		return res, nil
	}

	// now we deal with weight rules
	for _, rule := range feature.Rules.Enable {
		if ruleWeight(logger, ruleFields(rule.Field, rule.Fields), rule.Weight, vars) {
			logger.Debug("match: matched weight rule")
			res.AddStatus(featureName, true, setVars(logger, feature.Rules.SetVars, vars))
			return res, nil
		}
	}

	return res, nil
}

func setVars(logger logrus.FieldLogger, rules []cfg.SetVarRule, vars map[string]string) map[string]interface{} {
	if len(rules) == 0 {
		logger.Debug("no set_vars rules")
		return nil
	}

	logrus.Debugf("checking %d set_vars rules", len(rules))

	// now we deal with enable rules
	setVars := map[string]interface{}{}
	foreachSetVarField(rules, func(field string, rule cfg.SetVarRule) {
		t := varInSlice(rule.Values.Eq, field, vars)
		logger.Debugf("check: field '%s' in %#v matches set var rules %#v: %t", field, vars, rule.Values.Eq, t)
		if t {
			for k, v := range rule.Set {
				setVars[k] = v
			}
		}
	})

	for _, rule := range rules {
		if ruleWeight(logger, ruleFields(rule.Field, rule.Fields), rule.Weight, vars) {
			for k, v := range rule.Set {
				setVars[k] = v
			}
		}
	}

	logger.Debugf("finished checking %d set_vars rules", len(rules))

	return setVars
}

func ruleFields(field string, fields []string) []string {
	if field != "" {
		return []string{field}
	}
	return fields
}

func ruleWeight(logger logrus.FieldLogger, fields []string, weight int, vars map[string]string) bool {
	if weight < 0 || weight > 100 {
		logger.Debugf("weight (%d) outside range 0-100", weight)
		return false
	}

	// first build a string containing all the key/value pairs
	b := bytes.Buffer{}
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
	t := c < weight
	logger.Debugf("check: hash result < weight (%d < %d): %t", c, weight, t)
	return t
}

func foreachDisableField(rules []cfg.DisableRule, fn func(string, cfg.DisableRule) bool) bool {
	for _, rule := range rules {
		for _, field := range ruleFields(rule.Field, rule.Fields) {
			if fn(field, rule) {
				return true
			}
		}
	}
	return false
}

func foreachEnableField(rules []cfg.EnableRule, fn func(string, cfg.EnableRule) bool) bool {
	for _, rule := range rules {
		for _, field := range ruleFields(rule.Field, rule.Fields) {
			if fn(field, rule) {
				return true
			}
		}
	}
	return false
}

func foreachSetVarField(rules []cfg.SetVarRule, fn func(string, cfg.SetVarRule)) bool {
	for _, rule := range rules {
		for _, field := range ruleFields(rule.Field, rule.Fields) {
			fn(field, rule)
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
