package reqcontext_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReqContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReqContext Suite")
}
