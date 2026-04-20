package handlers

import (
	"errors"
	"net/http"
	"strconv"

	taskdomain "example.com/taskservice/internal/domain/task"
	"example.com/taskservice/internal/usecase/recurrence"
	recurrenceusecase "example.com/taskservice/internal/usecase/recurrence"
	"github.com/gorilla/mux"
)

type RecurrenceHandler struct {
	usecase recurrence.Usecase
}

func NewRecurrenceHandler(usecase recurrence.Usecase) *RecurrenceHandler {
	return &RecurrenceHandler{usecase: usecase}
}

func (h *RecurrenceHandler) Create(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req recurrenceMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	input := recurrenceusecase.CreateUpdateInput{
		TaskID:              taskID,
		RecurrenceType:      req.RecurrenceType,
		RecurrenceModifiers: req.RecurrenceModifiers,
		EndDate:             req.EndDate,
		MaxOccurrences:      req.MaxOccurrences,
		Interval:            req.Interval,
		Days:                req.Days,
	}

	rule, err := h.usecase.Create(r.Context(), input, taskID)
	if err != nil {
		writeRecurrenceUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newRecurrenceRuleDTO(rule))
}

func (h *RecurrenceHandler) Get(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rule, err := h.usecase.GetByTaskID(r.Context(), taskID)
	if err != nil {
		writeRecurrenceUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newRecurrenceRuleDTO(rule))
}

func (h *RecurrenceHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req recurrenceMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rule, err := h.usecase.Update(r.Context(), taskID, recurrenceusecase.CreateUpdateInput{
		RecurrenceType:      req.RecurrenceType,
		RecurrenceModifiers: req.RecurrenceModifiers,
		EndDate:             req.EndDate,
		MaxOccurrences:      req.MaxOccurrences,
		Interval:            req.Interval,
		Days:                req.Days,
	})
	if err != nil {
		writeRecurrenceUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newRecurrenceRuleDTO(rule))
}

func (h *RecurrenceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.Delete(r.Context(), taskID); err != nil {
		writeRecurrenceUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeRecurrenceUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, recurrenceusecase.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, recurrenceusecase.ErrInvalidRecurrenceRule):
		writeError(w, http.StatusBadRequest, err)
	case errors.Is(err, recurrenceusecase.ErrTaskNotRecurring):
		writeError(w, http.StatusBadRequest, err)
	case errors.Is(err, recurrenceusecase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	case errors.Is(err, taskdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func getTaskIDFromRequest(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["task_id"]
	if rawID == "" {
		return 0, errors.New("missing task id")
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid task id")
	}

	if id <= 0 {
		return 0, errors.New("invalid task id")
	}

	return id, nil
}
