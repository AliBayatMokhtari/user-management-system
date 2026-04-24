package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"ums/model"
	"ums/service"

	"github.com/go-chi/chi/v5"
)

type APIResponse struct {
	Error      *string `json:"error"`
	Data       any     `json:"data"`
	Successful bool    `json:"successful"`
}

func respondSuccess(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := APIResponse{
		Error:      nil,
		Data:       data,
		Successful: true,
	}

	json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	msg := message
	resp := APIResponse{
		Error:      &msg,
		Data:       nil,
		Successful: false,
	}

	json.NewEncoder(w).Encode(resp)
}

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")

		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		respondError(w, http.StatusBadRequest, "name and email are required")

		return
	}

	user, err := h.svc.RegisterUser(r.Context(), req.Name, req.Email)

	if err != nil {
		if err.Error() == "email already exists" {
			respondError(w, http.StatusConflict, err.Error())

			return
		}

		respondError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	respondSuccess(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)

	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")

		return
	}

	user, err := h.svc.GetUser(r.Context(), uint(id))

	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	if user == nil {
		respondError(w, http.StatusNotFound, "user not found")

		return
	}

	respondSuccess(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())

	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	if users == nil {
		users = []model.User{}
	}

	respondSuccess(w, http.StatusOK, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)

	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")

		return
	}

	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")

		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		respondError(w, http.StatusBadRequest, "name and email are required")

		return
	}

	user, err := h.svc.UpdateUser(r.Context(), uint(id), req.Name, req.Email)

	if err != nil {
		if err.Error() == "user not found" {
			respondError(w, http.StatusNotFound, err.Error())

			return
		}

		respondError(w, http.StatusInternalServerError, "internal server error")

		return

	}

	respondSuccess(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)

	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")

		return
	}

	err = h.svc.DeleteUser(r.Context(), uint(id))

	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	respondSuccess(w, http.StatusNoContent, nil)
}
