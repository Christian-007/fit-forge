package middlewares_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	mock_services "github.com/Christian-007/fit-forge/internal/app/auth/services/mocks"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Jwt Auth Middleware", func() {
	var (
		ctrl             *gomock.Controller
		mockAuthService       *mock_services.MockAuthService
		jwtAuthMiddleware func(http.Handler) http.Handler
		next             http.HandlerFunc
		request          *http.Request
		recorder         *httptest.ResponseRecorder
	)

	BeforeEach(func ()  {
		ctrl = gomock.NewController(GinkgoT())
		mockAuthService = mock_services.NewMockAuthService(ctrl)
		jwtAuthMiddleware = middlewares.JwtAuth(mockAuthService)

		request = httptest.NewRequest("GET", "/test-uri", nil)
		recorder = httptest.NewRecorder()
	})

	AfterEach(func ()  {
		request.Header.Set("Authorization", "")
		ctrl.Finish()
	})

	It("should return http status 200 if the auth token is valid", func ()  {
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessTokenUuid, _ := requestctx.AccessTokenUuid(r.Context())	
			Expect(accessTokenUuid).To(Equal("a2b3c4"))
		})
		request.Header.Set("Authorization", "Bearer 123")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{
			UserID: 1,
			Uuid: "a2b3c4",
		}, nil)

		handler := jwtAuthMiddleware(next)
		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusOK))
	})

	It("should return http status 401 unauthorized if the Authorization header is nil", func ()  {
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Empty mock body function
		})

		handler := jwtAuthMiddleware(next)
		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Unauthorized"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return http status 401 unauthorized if the Authorization header value is only the token", func ()  {
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Empty mock body function
		})
		request.Header.Set("Authorization", "123")

		handler := jwtAuthMiddleware(next)
		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Unauthorized"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return http status 401 unauthorized if the token is not a valid JWT", func ()  {
		request.Header.Set("Authorization", "Bearer 123")
		mockError := errors.New("Token is not a valid JWT")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{}, mockError)
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Empty mock body function
		})

		handler := jwtAuthMiddleware(next)
		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Token is invalid"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})
})
