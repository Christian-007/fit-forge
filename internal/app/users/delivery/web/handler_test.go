package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Christian-007/fit-forge/internal/app/users/delivery/web"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	mock_services "github.com/Christian-007/fit-forge/internal/app/users/services/mocks"
	mock_applog "github.com/Christian-007/fit-forge/internal/pkg/applog/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User Handler", func() {
	var (
		ctrl            *gomock.Controller
		mockUserService *mock_services.MockUserService
		mockLogger      *mock_applog.MockLogger
		userHandler     web.UserHandler
		recorder        *httptest.ResponseRecorder
		request         *http.Request
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockUserService = mock_services.NewMockUserService(ctrl)
		mockLogger = mock_applog.NewMockLogger(ctrl)
		userHandler = web.NewUserHandler(web.UserHandlerOptions{
			UserService: mockUserService,
			Logger:      mockLogger,
		})
		recorder = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/users", nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Get All", func() {
		mockGetAllUsersResponse := []dto.UserResponse{{
			Id:    1,
			Name:  "John Doe",
			Email: "johndoe@gmail.com",
		}, {
			Id:    2,
			Name:  "Mark",
			Email: "mark@gmail.com",
		}}
		When("there is no error from UserService", func() {
			It("should return 200 with a list of users", func() {
				mockUserService.EXPECT().GetAll().Return(mockGetAllUsersResponse, nil)

				userHandler.GetAll(recorder, request)

				mockLogger.EXPECT().Error(gomock.Any()).Times(0)
				Expect(recorder.Code).To(Equal(http.StatusOK))

				actual := recorder.Body.Bytes()
				expected := `
					{
						"results":[
							{"id":1,"name":"John Doe","email":"johndoe@gmail.com"},
							{"id":2,"name":"Mark","email":"mark@gmail.com"}
						]
					}`
				Expect(actual).To(MatchJSON(expected))
			})
		})
	})
})

func TestUserHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserHandler Suite")
}
