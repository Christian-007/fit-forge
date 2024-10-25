package middlewares_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	mock_services "github.com/Christian-007/fit-forge/internal/app/auth/services/mocks"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Logout Session Middleware", func() {
	var (
		ctrl             *gomock.Controller
		mockAuthService       *mock_services.MockAuthService
		logoutSessionMiddleware func(http.Handler) http.Handler
		nextHandler             http.HandlerFunc
		handler http.Handler
		request          *http.Request
		recorder         *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockAuthService = mock_services.NewMockAuthService(ctrl)
		logoutSessionMiddleware = middlewares.LogoutSession(mockAuthService)

		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only a mock HTTP status code
			w.WriteHeader(http.StatusOK)
		})
		handler = logoutSessionMiddleware(nextHandler)

		request = httptest.NewRequest("GET", "/logout", nil)
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

	It("should return http status 200 if the token is already expired", func ()  {
		request.Header.Set("Authorization", "Bearer 123")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{}, apperrors.ErrExpiredToken)

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Logout successful"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusOK))
	})

	It("should return http status 401 unauthorized if the token is not a valid JWT", func ()  {
		request.Header.Set("Authorization", "Bearer 123")
		mockError := errors.New("Token is not a valid JWT")
		mockAuthService.EXPECT().ValidateToken("123").Return(&domains.Claims{}, mockError)

		handler.ServeHTTP(recorder, request)

		expected := apphttp.ErrorResponse{Message: "Token is invalid"}
		var result apphttp.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	Context("when GetHashAuthDataFromCache returns an error", func ()  {
		BeforeEach(func ()  {
			mockAccessToken := "someAccessToken123"
			request.Header.Set("Authorization", "Bearer " + mockAccessToken)
			mockAuthService.EXPECT().ValidateToken(mockAccessToken).Return(&domains.Claims{
				Uuid: "mockAccessTokenUuid",	
			}, nil)
		})

		It("should call next handler with the correct context if the value exists in the cache using old data structure", func() {
			mockUserId := 123
			mockCtx := requestctx.WithUserId(context.Background(), mockUserId)
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, apperrors.ErrRedisValueNotInHash)
			mockAuthService.EXPECT().GetAuthDataFromCache("mockAccessTokenUuid").Return(mockUserId, nil)
			nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				By("having the correct user id")
				userId, _ := requestctx.UserId(r.Context())
				Expect(userId).To(Equal(mockUserId))
	
				By("having the correct user role")
				role, _ := requestctx.Role(r.Context())
				Expect(role).To(Equal(0))
			})
			handler = logoutSessionMiddleware(nextHandler)

			handler.ServeHTTP(recorder, request.WithContext(mockCtx))
		})

		It("should return http status 401 unauthorized if apperrors.ErrRedisValueNotInHash and apperrors.ErrRedisKeyNotFound from GetAuthDataFromCache returned", func() {
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, apperrors.ErrRedisValueNotInHash)
			mockAuthService.EXPECT().GetAuthDataFromCache("mockAccessTokenUuid").Return(-1, apperrors.ErrRedisKeyNotFound)
			
			handler.ServeHTTP(recorder, request)

			expected := apphttp.ErrorResponse{Message: "Unauthorized"}
			var result apphttp.ErrorResponse
			err := json.NewDecoder(recorder.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return http status 500 if apperrors.ErrRedisValueNotInHash and any other errors from GetAuthDataFromCache returned", func() {
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, apperrors.ErrRedisValueNotInHash)
			mockAuthService.EXPECT().GetAuthDataFromCache("mockAccessTokenUuid").Return(-1, errors.New("unexpected error occurred"))
			
			handler.ServeHTTP(recorder, request)

			expected := apphttp.ErrorResponse{Message: "Internal Server Error"}
			var result apphttp.ErrorResponse
			err := json.NewDecoder(recorder.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		})
		
		It("should return http status 401 unauthorized if apperrors.ErrRedisKeyNotFound is returned", func() {
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, apperrors.ErrRedisKeyNotFound)
			
			handler.ServeHTTP(recorder, request)

			expected := apphttp.ErrorResponse{Message: "Unauthorized"}
			var result apphttp.ErrorResponse
			err := json.NewDecoder(recorder.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return http status 500 for any other errors returned", func() {
			mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(domains.AuthData{}, errors.New("some other errors"))
			
			handler.ServeHTTP(recorder, request)

			expected := apphttp.ErrorResponse{Message: "Internal Server Error"}
			var result apphttp.ErrorResponse
			err := json.NewDecoder(recorder.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	It("should call next handler with the correct context if there is no error at all", func() {
		mockAuthData := domains.AuthData{
			UserId: 123,
			Role: 2,
		}
		mockAccessToken := "someAccessToken123"
		request.Header.Set("Authorization", "Bearer " + mockAccessToken)
		mockAuthService.EXPECT().ValidateToken(mockAccessToken).Return(&domains.Claims{
			Uuid: "mockAccessTokenUuid",	
		}, nil)
		mockCtx := requestctx.WithUserId(context.Background(), mockAuthData.UserId)
		mockCtx = requestctx.WithRole(mockCtx, mockAuthData.Role)
		mockAuthService.EXPECT().GetHashAuthDataFromCache("mockAccessTokenUuid").Return(mockAuthData, nil)
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			By("having the correct user id")
			userId, _ := requestctx.UserId(r.Context())
			Expect(userId).To(Equal(mockAuthData.UserId))

			By("having the correct user role")
			role, _ := requestctx.Role(r.Context())
			Expect(role).To(Equal(mockAuthData.Role))
		})
		handler = logoutSessionMiddleware(nextHandler)

		handler.ServeHTTP(recorder, request.WithContext(mockCtx))
	})
})
