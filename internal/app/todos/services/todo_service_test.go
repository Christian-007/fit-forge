package services_test

import (
	"errors"

	"github.com/Christian-007/fit-forge/internal/app/todos/domains"
	"github.com/Christian-007/fit-forge/internal/app/todos/dto"
	mock_repositories "github.com/Christian-007/fit-forge/internal/app/todos/repositories/mocks"
	"github.com/Christian-007/fit-forge/internal/app/todos/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Todo Service", func() {
	var (
		ctrl               *gomock.Controller
		mockTodoRepository *mock_repositories.MockTodoRepository
		todoService        services.TodoService
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockTodoRepository = mock_repositories.NewMockTodoRepository(ctrl)
		todoService = services.NewTodoService(services.TodoServiceOptions{
			TodoRepository: mockTodoRepository,
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Get All", func() {
		It("should return an error when TodoRepository.GetAll() returns an error", func ()  {
			mockGetAllTodosDto := []dto.GetAllTodosResponse{}
			mockTodoModel := []domains.TodoModel{}
			mockTodoRepository.EXPECT().GetAll().Return(mockTodoModel, errors.New("some error"))

			todos, err := todoService.GetAll()

			Expect(todos).To(Equal(mockGetAllTodosDto))
			Expect(err).To(MatchError("some error"))
		})

		It("should return a list of todos when TodoRepository.GetAll() returns a success", func ()  {
			mockGetAllTodosDto := []dto.GetAllTodosResponse{
				{Id: 1, Title: "Todo 1", IsCompleted: false, UserId: 1},
				{Id: 2, Title: "Todo 2", IsCompleted: false, UserId: 1},
			}
			mockTodoModel := []domains.TodoModel{
				{Id: 1, Title: "Todo 1", IsCompleted: false, UserId: 1},
				{Id: 2, Title: "Todo 2", IsCompleted: false, UserId: 1},
			}
			mockTodoRepository.EXPECT().GetAll().Return(mockTodoModel, nil)

			todos, err := todoService.GetAll()

			Expect(todos).To(Equal(mockGetAllTodosDto))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get All By User Id", func() {
		It("should return an error when TodoRepository.GetAllByUserId() returns an error", func ()  {
			mockTodoResponse := []dto.TodoResponse{}
			mockTodoModel := []domains.TodoModel{}
			mockError := errors.New("some error")
			mockTodoRepository.EXPECT().GetAllByUserId(1).Return(mockTodoModel, mockError)

			todos, err := todoService.GetAllByUserId(1)

			Expect(todos).To(Equal(mockTodoResponse))
			Expect(err).To(MatchError(mockError))
		})

		It("should return a list of todos when TodoRepository.GetAllByUserId() returns a success", func ()  {
			mockTodoResponse := []dto.TodoResponse{
				{Id: 1, Title: "Todo 1", IsCompleted: false},
				{Id: 2, Title: "Todo 2", IsCompleted: false},
			}
			mockTodoModel := []domains.TodoModel{
				{Id: 1, Title: "Todo 1", IsCompleted: false},
				{Id: 2, Title: "Todo 2", IsCompleted: false},
			}
			mockTodoRepository.EXPECT().GetAllByUserId(1).Return(mockTodoModel, nil)

			todos, err := todoService.GetAllByUserId(1)

			Expect(todos).To(Equal(mockTodoResponse))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get One By User Id", func() {
		It("should return an error when TodoRepository.GetOneByUserId() returns a todo not found error", func ()  {
			mockTodoResponse := dto.TodoResponse{}
			mockTodoModel := domains.TodoModel{}
			mockTodoRepository.EXPECT().GetOneByUserId(1, 2).Return(mockTodoModel, pgx.ErrNoRows)

			todo, err := todoService.GetOneByUserId(1,2)

			Expect(todo).To(Equal(mockTodoResponse))
			Expect(err).To(MatchError(apperrors.ErrTodoNotFound))
		})

		It("should return an error when TodoRepository.GetOneByUserId() returns an unexpected error", func ()  {
			mockTodoResponse := dto.TodoResponse{}
			mockTodoModel := domains.TodoModel{}
			mockError := errors.New("an unxpected error")
			mockTodoRepository.EXPECT().GetOneByUserId(1, 2).Return(mockTodoModel, mockError)

			todo, err := todoService.GetOneByUserId(1,2)

			Expect(todo).To(Equal(mockTodoResponse))
			Expect(err).To(MatchError(mockError))
		})

		It("should return a correct todo when TodoRepository.GetAllByUserId() returns a success", func ()  {
			mockTodoResponse := dto.TodoResponse{Id: 1, Title: "Todo 1", IsCompleted: false}
			mockTodoModel := domains.TodoModel{Id: 1, Title: "Todo 1", IsCompleted: false}
			mockTodoRepository.EXPECT().GetOneByUserId(1,2).Return(mockTodoModel, nil)

			todo, err := todoService.GetOneByUserId(1,2)

			Expect(todo).To(Equal(mockTodoResponse))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Create", func() {
		It("should return an error when TodoRepository.Create() returns an unexpected error", func ()  {
			mockCreateTodoRequest := dto.CreateTodoRequest{
				Title: "A new Todo",
			}
			mockTodoResponse := dto.TodoResponse{}
			mockTodoModel := domains.TodoModel{Title: mockCreateTodoRequest.Title}
			mockError := errors.New("an unxpected error")
			mockTodoRepository.EXPECT().Create(1, mockTodoModel).Return(domains.TodoModel{}, mockError)

			todo, err := todoService.Create(1,mockCreateTodoRequest)

			Expect(todo).To(Equal(mockTodoResponse))
			Expect(err).To(MatchError(mockError))
		})

		It("should return the created todo when TodoRepository.Create() returns a success", func ()  {
			mockCreateTodoRequest := dto.CreateTodoRequest{
				Title: "A new Todo",
			}
			mockTodoResponse := dto.TodoResponse{Id: 1, Title: mockCreateTodoRequest.Title, IsCompleted: false}
			mockTodoModel := domains.TodoModel{Id: 1, Title: mockCreateTodoRequest.Title, IsCompleted: false, UserId: 1}
			mockTodoRepository.EXPECT().Create(1, domains.TodoModel{Title: mockCreateTodoRequest.Title}).Return(mockTodoModel, nil)

			todo, err := todoService.Create(1, mockCreateTodoRequest)

			Expect(todo).To(Equal(mockTodoResponse))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Delete", func ()  {
		It("should return an error when TodoRepository.Delete() returns an unexpected error", func ()  {
			mockError := errors.New("an unexpected error")
			mockTodoRepository.EXPECT().Delete(1,1).Return(mockError)

			err := todoService.Delete(1,1)

			Expect(err).To(MatchError(mockError))
		})

		It("should return no error if the delete operation is a success", func ()  {
			mockTodoRepository.EXPECT().Delete(1,1).Return(nil)

			err := todoService.Delete(1,1)

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
