package middlewares_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	mock_services "github.com/Christian-007/fit-forge/internal/app/auth/services/mocks"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Strict Session Middleware", func() {
	var (
		ctrl             *gomock.Controller
		mockAuthService       *mock_services.MockAuthService
		strictSessionMiddleware func(http.Handler) http.Handler
		nextHandler             http.HandlerFunc
		handler http.Handler
		request          *http.Request
		recorder         *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockAuthService = mock_services.NewMockAuthService(ctrl)
		strictSessionMiddleware = middlewares.StrictSession(mockAuthService)

		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only a mock HTTP status code
			w.WriteHeader(http.StatusOK)
		})
		handler = strictSessionMiddleware(nextHandler)

		request = httptest.NewRequest("GET", "/test-users", nil)
		recorder = httptest.NewRecorder()
	})
	
	AfterEach(func() {
		ctrl.Finish()
	})

	It("should return http status 401 unauthorized if the Authorization header is nil", func ()  {
		request.Header.Set("Authorization", "")

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Unauthorized"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return http status 401 unauthorized if the Authorization header value does not have the right format", func ()  {
		wrongAuthBearerFormat := "123" // should have 'Bearer'
		request.Header.Set("Authorization", wrongAuthBearerFormat)

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Unauthorized"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return http status 401 if the token is already expired", func ()  {
		request.Header.Set("Authorization", "Bearer 123")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{}, apperrors.ErrExpiredToken)

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Token is expired"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return http status 401 unauthorized if the token has an invalid signature", func ()  {
		request.Header.Set("Authorization", "Bearer 123")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{}, apperrors.ErrInvalidSignature)

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Token is invalid"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return http status 500 if there is an unexpected error in validating the token", func ()  {
		request.Header.Set("Authorization", "Bearer 123")
		mockError := errors.New("An unexpected error")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{}, mockError)

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Internal Server Error"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
	})

	Context("when GetHashAuthDataFromCache returns an error", func ()  {
		BeforeEach(func ()  {
			mockAccessToken := "someAccessToken123"
			request.Header.Set("Authorization", "Bearer " + mockAccessToken)
			mockAuthService.EXPECT().ValidateToken(mockAccessToken).Return(&domains.Claims{
				Uuid: "mockAccessTokenUuid",	
			}, nil)
		})

		It("should return http status 401 if the key is not found on Redis", func() {
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, apperrors.ErrRedisValueNotInHash)

			handler.ServeHTTP(recorder, request)

			expected := apphttp.ErrorResponse{Message: "Unauthorized"}
			var result apphttp.ErrorResponse
			err := json.NewDecoder(recorder.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return http status 500 if the error is unexpected", func() {
			mockError := errors.New("unexpected error")
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, mockError)

			handler.ServeHTTP(recorder, request)

			expected := apphttp.ErrorResponse{Message: "Internal Server Error"}
			var result apphttp.ErrorResponse
			err := json.NewDecoder(recorder.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
