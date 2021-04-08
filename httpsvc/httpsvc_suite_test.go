package httpsvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHTTPSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTPSvc Suite")
}
