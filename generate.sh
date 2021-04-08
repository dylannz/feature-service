#!/bin/bash

go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
"$GOPATH"/bin/oapi-codegen -config spec/spec.cfg.yaml spec/spec.yaml
"$GOPATH"/bin/oapi-codegen -config spec/types.cfg.yaml spec/spec.yaml