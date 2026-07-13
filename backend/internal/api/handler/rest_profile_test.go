package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

func TestGetPublicProfile(t *testing.T) {
	env := setupTestEnv(t)
	defer env.MiniRedis.Close()
	defer env.DB.Close()

	publicUsername := "public_user"
	_, err := env.UserRepo.CreateUser(ctx, service.CreateUserParams{
		Username:      publicUsername,
		PasswordHash:  "pass",
		FirstName:     "Public",
		LastName:      "User",
		Website:       "",
		IsPublic:      true,
		AgbAcceptedAt: time.Now(),
	})
	assert.NoError(t, err)

	privateUsername := "private_user"
	_, err = env.UserRepo.CreateUser(ctx, service.CreateUserParams{
		Username:      privateUsername,
		PasswordHash:  "pass",
		FirstName:     "Private",
		LastName:      "User",
		Website:       "",
		IsPublic:      false,
		AgbAcceptedAt: time.Now(),
	})
	assert.NoError(t, err)

	t.Run("Get Public Profile - Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+publicUsername, nil)
		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp user.GetPublicProfileResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, publicUsername, resp.GetProfile().GetUsername())
	})

	t.Run("Get Private Profile - Forbidden/NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+privateUsername, nil)
		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
		var prob apperrors.ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &prob)
		assert.NoError(t, err)
		assert.Equal(t, apperrors.TypeNotFound, prob.Type)
		assert.Equal(t, apperrors.ErrPublicProfileNotFoundOrPrivate.Error(), prob.Detail)
	})

	t.Run("Get Non-Existent Profile - NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/non_existent", nil)
		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
		var prob apperrors.ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &prob)
		assert.NoError(t, err)
		assert.Equal(t, apperrors.TypeNotFound, prob.Type)
		assert.Equal(t, apperrors.ErrPublicProfileNotFoundOrPrivate.Error(), prob.Detail)
	})
}

func TestUpdateUser_Privacy(t *testing.T) {
	env := setupTestEnv(t)
	defer env.MiniRedis.Close()
	defer env.DB.Close()

	username := "privacy_tester"
	createdUser, err := env.UserRepo.CreateUser(ctx, service.CreateUserParams{
		Username:      username,
		PasswordHash:  "pass",
		FirstName:     "Privacy",
		LastName:      "Tester",
		Website:       "",
		IsPublic:      false,
		AgbAcceptedAt: time.Now(),
	})
	assert.NoError(t, err)

	token, _ := service.GenerateToken(createdUser, testSecret)

	reqBody := `{"first_name": "Privacy", "last_name": "Tester", "website": "", "is_public": true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/profile", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	env.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	updatedUser, err := env.UserRepo.GetUser(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.True(t, updatedUser.IsPublic, "User should be public after update")

	reqBody = `{"first_name": "Privacy", "last_name": "Tester", "website": "", "is_public": false}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/profile", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	env.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	updatedUser, err = env.UserRepo.GetUser(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.False(t, updatedUser.IsPublic, "User should be private after second update")
}
