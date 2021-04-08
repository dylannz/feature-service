package main

import (
	"net/http"

	"github.com/Netflix/go-env"
	"github.com/dylannz/feature-service/cfg"
	"github.com/dylannz/feature-service/httpsvc"
	"github.com/dylannz/feature-service/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate ./generate.sh

type Env struct {
	LogLevel  string `env:"LOG_LEVEL"`
	ConfigDir string `env:"CONFIG_DIR"`
	HTTPAddr  string `env:"HTTP_ADDR"`
}

func initEnv() Env {
	e := Env{
		LogLevel:  "info",
		ConfigDir: "./config",
		HTTPAddr:  "127.0.0.1:3000",
	}

	_, err := env.UnmarshalFromEnviron(&e)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "parse environment variables"))
	}

	lvl, err := logrus.ParseLevel(e.LogLevel)
	if err != nil {
		logrus.Fatal(err, "parse log level")
	}
	logrus.SetLevel(lvl)
	return e
}

func main() {
	e := initEnv()
	logger := logrus.WithField("service", "feature-service")

	config, err := cfg.LoadYAMLDir(e.ConfigDir)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := service.NewService(logger, config)
	h := httpsvc.NewHTTPHandler(logger, svc)
	logger.Fatal(http.ListenAndServe(e.HTTPAddr, h))
}
