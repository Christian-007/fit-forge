package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_applog "github.com/Christian-007/fit-forge/internal/pkg/applog/mocks"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Logger Middleware", func() {
	var (
		ctrl             *gomock.Controller
		mockLogger       *mock_applog.MockLogger
		loggerMiddleware func(http.Handler) http.Handler
		next             http.HandlerFunc
		request          *http.Request
		recorder         *httptest.ResponseRecorder
		nextCalled       bool
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockLogger = mock_applog.NewMockLogger(ctrl)
		loggerMiddleware = middlewares.NewLogRequest(mockLogger)

		nextCalled = false
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
		})

		request = httptest.NewRequest("GET", "/test-uri", nil)
		request.RemoteAddr = "127.0.0.1:3000"
		request.Proto = "HTTP/1.1"

		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should log the request and call the next handler", func() {
		mockLogger.EXPECT().Info(
			"received request",
			"ip", "127.0.0.1:3000",
			"proto", "HTTP/1.1",
			"method", "GET",
			"uri", "/test-uri",
		)

		handler := loggerMiddleware(next)
		handler.ServeHTTP(recorder, request)

		Expect(nextCalled).To(BeTrue())
	})

})

func TestLoggerMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoggerMiddleware Suite")
}
