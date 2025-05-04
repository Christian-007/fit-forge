package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/todos/dto"
	"github.com/Christian-007/fit-forge/internal/app/todos/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/pkg/topics"
	"github.com/Christian-007/fit-forge/internal/utils"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type TodoHandler struct {
	TodoHandlerOptions
}

type TodoHandlerOptions struct {
	TodoService services.TodoService
	Logger      applog.Logger
	Publisher   message.Publisher
}

func NewTodoHandler(options TodoHandlerOptions) TodoHandler {
	return TodoHandler{
		options,
	}
}

func (t TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	getAllTodosResponse, err := t.TodoService.GetAll()
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	res := apphttp.CollectionRes[dto.GetAllTodosResponse]{
		Results: getAllTodosResponse,
	}
	utils.SendResponse(w, http.StatusOK, res)
}

func (t TodoHandler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := requestctx.UserId(r.Context())
	if !ok {
		t.Logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	todoResponse, err := t.TodoService.GetAllByUserId(userId)
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	res := apphttp.CollectionRes[dto.TodoResponse]{
		Results: todoResponse,
	}
	utils.SendResponse(w, http.StatusOK, res)
}

func (t TodoHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Todo ID is invalid"})
		return
	}

	userId, ok := requestctx.UserId(r.Context())
	if !ok {
		t.Logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	todoResponse, err := t.TodoService.GetOneByUserId(userId, todoId)
	if err != nil {
		t.Logger.Error(err.Error())

		if err == apperrors.ErrTodoNotFound {
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Todo not found"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, todoResponse)
}

func (t TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createTodoRequest dto.CreateTodoRequest
	err := json.NewDecoder(r.Body).Decode(&createTodoRequest)
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	if err = createTodoRequest.Validate(); err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: err.Error()})
		return
	}

	userId, ok := requestctx.UserId(r.Context())
	if !ok {
		t.Logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	todoResponse, err := t.TodoService.CreateWithPoints(r.Context(), userId, createTodoRequest)
	if err != nil {
		t.Logger.Error("Unable to create a todo",
			slog.String("error", err.Error()),
		)
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, todoResponse)
}

func (t TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Todo ID is invalid"})
		return
	}

	userId, ok := requestctx.UserId(r.Context())
	if !ok {
		t.Logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	err = t.TodoService.Delete(todoId, userId)
	if err != nil {
		if err == apperrors.ErrUserOrTodoNotFound {
			t.Logger.Error(err.Error())
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO: add validations and publish complete todo if completed
func (t TodoHandler) Patch(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Todo ID is invalid"})
		return
	}

	var updateTodoReq dto.UpdateTodoRequest
	err = json.NewDecoder(r.Body).Decode(&updateTodoReq)
	if err != nil {
		t.Logger.Error("Error unmarshalling update todo request",
			slog.String("error", err.Error()),
		)
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	// Validate if the update request is empty
	if (updateTodoReq == dto.UpdateTodoRequest{}) {
		t.Logger.Error("Error validations",
			slog.String("error", "Update request cannot be empty"),
		)
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Bad Request"})
		return
	}

	if err = updateTodoReq.Validate(); err != nil {
		t.Logger.Error("Error validations",
			slog.String("error", err.Error()),
		)
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Bad Request"})
		return
	}

	err = t.TodoService.Update(r.Context(), todoId, updateTodoReq)
	if err != nil {
		t.Logger.Error("Error updating todo",
			slog.String("error", err.Error()),
		)
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	t.Logger.Info("Successfully completed a todo", slog.Int("todoId", todoId))

	// Publish "TodoCompleted" topic when completing a todo
	if updateTodoReq.IsCompletedTrue() {
		userId, ok := requestctx.UserId(r.Context())
		if !ok {
			t.Logger.Error(apperrors.ErrTypeAssertion.Error())
			utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
			return
		}

		msg := message.NewMessage(watermill.NewUUID(), []byte(strconv.Itoa(userId)))
		err = t.Publisher.Publish(topics.TodoCompleted, msg)
		// If there is an unexpected error, it's decided to not send any http error response
		if err != nil {
			t.Logger.Error("Fail to publish TodoCompleted", slog.String("error:", err.Error()))
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
