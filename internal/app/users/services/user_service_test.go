package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	mock_repositories "github.com/Christian-007/fit-forge/internal/app/users/repositories/mocks"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
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

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Get All", func() {
		mockUserModels := []domains.UserModel{{
			Id:        1,
			Name:      "John Doe",
			Email:     "johndoe@gmail.com",
			Password:  []byte{100},
			CreatedAt: time.Date(2024, 02, 01, 1, 1, 1, 0, time.UTC),
		}, {
			Id:        2,
			Name:      "Mark",
			Email:     "mark@gmail.com",
			Password:  []byte{100},
			CreatedAt: time.Date(2024, 02, 01, 1, 1, 1, 0, time.UTC),
		}}

		When("there is no error", func() {
			It("should return a list of users", func() {
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
		When("there is an error", func() {
			It("should return an empty user array with the error", func() {
				mockEmptyUserResponse := []dto.UserResponse{}

				mockUserRepository.EXPECT().GetAll().Return(mockUserModels, errors.New("Some error"))
				result, err := userService.GetAll()

				Expect(result).To(Equal(mockEmptyUserResponse))
				Expect(err).To(MatchError("Some error"))
			})
		})
	})
})

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")
}
