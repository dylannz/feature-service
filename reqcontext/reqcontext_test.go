package reqcontext_test

import (
	"context"
	"strings"

	. "github.com/dylannz/feature-service/reqcontext"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reqcontext", func() {
	Context("with a blank request id", func() {
		It("generates an internal request uuid", func() {
			ctx := context.Background()
			ctx = ContextWithRequestID(ctx, "")
			requestID := RequestIDFromContext(ctx)
			parts := strings.Split(requestID, ":")
			Expect(len(parts)).To(Equal(2))
			Expect(parts[0]).To(Equal("internal"))

			// ensure second part is a valid uuid
			_, err := uuid.Parse(parts[1])
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("with a blank request id", func() {
		It("stores and retrieves the user request id", func() {
			const requestID = "ad85c722-b6fa-4f2b-b022-a8dcd2f228af"
			ctx := context.Background()
			ctx = ContextWithRequestID(ctx, requestID)
			actual := RequestIDFromContext(ctx)
			Expect(actual).To(Equal("user:" + requestID))
		})
	})

	Describe("RequestIDFromContext", func() {
		Context("with no previously stored requestID", func() {
			It("returns an empty string", func() {
				actual := RequestIDFromContext(context.Background())
				Expect(actual).To(Equal(""))
			})
		})
	})
})
