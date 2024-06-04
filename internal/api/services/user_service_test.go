package services_test

import (
	"testing"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/dto"
	mock_repositories "github.com/Christian-007/fit-forge/internal/api/repositories/mocks"
	"github.com/Christian-007/fit-forge/internal/api/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User Service", func() {
	var (
		ctrl               *gomock.Controller
		mockUserRepository *mock_repositories.MockUserRepository
		userService        services.UserService
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockUserRepository = mock_repositories.NewMockUserRepository(ctrl)
		userService = services.NewUserService(services.UserServiceOptions{
			UserRepository: mockUserRepository,
		})
	})

	Describe("Get All", func() {
		When("there is no error", func() {
			It("should return a list of users", func() {
				mockUserModels := []domains.UserModel{{
					Id:    1,
					Name:  "John Doe",
					Email: "johndoe@gmail.com",
				}, {
					Id:    2,
					Name:  "Mark",
					Email: "mark@gmail.com",
				}}
				mockUserResponse := []dto.UserResponse{{
					Id:    1,
					Name:  "John Doe",
					Email: "johndoe@gmail.com",
				}, {
					Id:    2,
					Name:  "Mark",
					Email: "mark@gmail.com",
				}}

				mockUserRepository.EXPECT().GetAll().Return(mockUserModels, nil)
				result, err := userService.GetAll()

				Expect(result).To(Equal(mockUserResponse))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")
}
