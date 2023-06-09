package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Wave-95/pgserver/internal/apiresponse"
	"github.com/Wave-95/pgserver/pkg/logger"
	"github.com/Wave-95/pgserver/pkg/validator"
	"github.com/go-chi/chi"
)

var (
	ErrInternalServer        = errors.New("Internal server error")
	ErrGetUserInvalidRequest = errors.New("Invalid get user request")
	ErrGetUserEncodeJSON     = errors.New("Error encoding user to JSON")
)

type GetUserRequest struct {
	UserID string `validate:"required,uuid4"`
}

func (r GetUserRequest) Validate(v validator.Validate) error {
	return v.Struct(r)
}

type GetUserResponse struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (api *API) handleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := logger.FromContext(ctx)

	// Validate get user request
	userID := chi.URLParam(r, "userID")
	input := GetUserRequest{UserID: userID}
	if err := input.Validate(api.validate); err != nil {
		apiresponse.RespondWithError(w, http.StatusBadRequest, ErrGetUserInvalidRequest)
		return
	}

	// Get user and handle errors
	user, err := api.service.GetUser(ctx, input.UserID)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			apiresponse.RespondWithError(w, http.StatusNotFound, ErrUserNotFound)
		default:
			apiresponse.RespondWithError(w, http.StatusInternalServerError, ErrInternalServer)
			l.Errorf("Issue getting user: %s", err)
		}
		return
	}

	// Write user response
	res := GetUserResponse{
		Id:        user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		apiresponse.RespondWithError(w, http.StatusInternalServerError, ErrGetUserEncodeJSON)
	}
}
