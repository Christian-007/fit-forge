package services_test

import (
	"errors"
	"time"

	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	mock_repositories "github.com/Christian-007/fit-forge/internal/app/users/repositories/mocks"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/jackc/pgx/v5"
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

	Describe("Get One", func() {
		mockUserModel := domains.UserModel{
			Id:        1,
			Name:      "John Doe",
			Email:     "johndoe@gmail.com",
			Password:  []byte{100},
			CreatedAt: time.Date(2024, 02, 01, 1, 1, 1, 0, time.UTC),
		}

		When("there is no error", func() {
			It("should return one user", func() {
				mockUserId := 1
				mockUserResponse := dto.UserResponse{
					Id:    mockUserId,
					Name:  "John Doe",
					Email: "johndoe@gmail.com",
				}

				mockUserRepository.EXPECT().GetOne(mockUserId).Return(mockUserModel, nil)
				result, err := userService.GetOne(mockUserId)

				Expect(result).To(Equal(mockUserResponse))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("there is an unexpected error", func() {
			It("should return an empty user and the error", func() {
				mockUserId := 1
				mockEmptyUserModel := domains.UserModel{}
				mockEmptyUserResponse := dto.UserResponse{}
				mockError := errors.New("Some error")

				mockUserRepository.EXPECT().GetOne(mockUserId).Return(mockEmptyUserModel, mockError)
				result, err := userService.GetOne(mockUserId)

				Expect(result).To(Equal(mockEmptyUserResponse))
				Expect(err).To(MatchError(mockError))
			})
		})

		When("there is no row returned", func() {
			It("should return an empty user and the user not found error", func() {
				mockUserId := 1
				mockEmptyUserModel := domains.UserModel{}
				mockEmptyUserResponse := dto.UserResponse{}

				mockUserRepository.EXPECT().GetOne(mockUserId).Return(mockEmptyUserModel, pgx.ErrNoRows)
				result, err := userService.GetOne(mockUserId)

				Expect(result).To(Equal(mockEmptyUserResponse))
				Expect(err).To(MatchError(apperrors.ErrUserNotFound))
			})
		})
	})

	Describe("Get One By Email", func() {
		It("should return a user not found error if the user does not exist", func() {
			mockEmail := "email@exmaple.com"
			mockEmptyUserModel := domains.UserModel{}
			mockEmptyUserResponse := dto.GetUserByEmailResponse{}
			mockUserRepository.EXPECT().GetOneByEmail(mockEmail).Return(mockEmptyUserModel, pgx.ErrNoRows)

			user, err := userService.GetOneByEmail(mockEmail)

			Expect(user).To(Equal(mockEmptyUserResponse))
			Expect(err).To(MatchError(apperrors.ErrUserNotFound))
		})

		It("should return the unexpected error if there is an unexpected error", func() {
			mockEmail := "email@exmaple.com"
			mockEmptyUserModel := domains.UserModel{}
			mockEmptyUserResponse := dto.GetUserByEmailResponse{}
			mockError := errors.New("an unexpected error")
			mockUserRepository.EXPECT().GetOneByEmail(mockEmail).Return(mockEmptyUserModel, mockError)

			user, err := userService.GetOneByEmail(mockEmail)

			Expect(user).To(Equal(mockEmptyUserResponse))
			Expect(err).To(MatchError(mockError))
		})

		It("should return the correct user if the user is found", func() {
			mockEmail := "email@exmaple.com"
			mockUserModel := domains.UserModel{
				Id:        1,
				Name:      "John",
				Email:     mockEmail,
				Password:  []byte("asd"),
				Role:      2,
				CreatedAt: time.Date(2024, 02, 01, 1, 1, 1, 0, time.UTC),
			}
			mockUserResponse := dto.GetUserByEmailResponse{
				Id:       1,
				Name:     "John",
				Email:    mockEmail,
				Password: []byte("asd"),
				Role:     2,
			}
			mockUserRepository.EXPECT().GetOneByEmail(mockEmail).Return(mockUserModel, nil)

			user, err := userService.GetOneByEmail(mockEmail)

			Expect(user).To(Equal(mockUserResponse))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
