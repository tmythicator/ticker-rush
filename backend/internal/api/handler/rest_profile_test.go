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
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	publicUsername := "public_user"
	_, err := userRepo.CreateUser(ctx, service.CreateUserParams{
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
	_, err = userRepo.CreateUser(ctx, service.CreateUserParams{
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
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp user.GetPublicProfileResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, publicUsername, resp.GetProfile().GetUsername())
	})

	t.Run("Get Private Profile - Forbidden/NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+privateUsername, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
		var prob apperrors.ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &prob)
		assert.NoError(t, err)
		assert.Equal(t, apperrors.TypeNotFound, prob.Type)
		assert.Equal(t, "User not found or profile is private", prob.Detail)
	})

	t.Run("Get Non-Existent Profile - NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/non_existent", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
		var prob apperrors.ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &prob)
		assert.NoError(t, err)
		assert.Equal(t, apperrors.TypeNotFound, prob.Type)
		assert.Equal(t, "User not found or profile is private", prob.Detail)
	})
}

func TestUpdateUser_Privacy(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	username := "privacy_tester"
	createdUser, err := userRepo.CreateUser(ctx, service.CreateUserParams{
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
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	updatedUser, err := userRepo.GetUser(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.True(t, updatedUser.IsPublic, "User should be public after update")

	reqBody = `{"first_name": "Privacy", "last_name": "Tester", "website": "", "is_public": false}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/profile", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	updatedUser, err = userRepo.GetUser(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.False(t, updatedUser.IsPublic, "User should be private after second update")
}
