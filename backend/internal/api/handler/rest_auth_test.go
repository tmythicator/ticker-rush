package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

func TestCreateUser(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	reqBody := fmt.Sprintf(
		`{"username": "%s", "password": "password123", "first_name": "Test", "last_name": "User", "agb_accepted": true}`,
		testUsername,
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp user.CreateUserResponse

	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	userObj, _, err := userRepo.GetUserByUsername(ctx, resp.GetUser().GetUsername())
	assert.NoError(t, err)
	assert.Equal(t, testUsername, userObj.Username)
}

func TestCreateUser_AgbNotAccepted(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	reqBody := `{"username": "test_user_2", "password": "password123", "first_name": "Test", "last_name": "User", "agb_accepted": false}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	var prob apperrors.ProblemDetails
	err := json.Unmarshal(w.Body.Bytes(), &prob)
	assert.NoError(t, err)
	assert.Equal(t, apperrors.TypeValidation, prob.Type)
	assert.Equal(t, apperrors.ErrAGBNotAccepted.Error(), prob.Detail)
}

func TestLogin(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err := userRepo.CreateUser(ctx, service.CreateUserParams{
		Username:      testUsername,
		PasswordHash:  string(hashedPassword),
		FirstName:     "Test",
		LastName:      "User",
		Website:       "",
		IsPublic:      false,
		AgbAcceptedAt: time.Now(),
	})
	assert.NoError(t, err)

	reqBody := fmt.Sprintf(`{"username": "%s", "password": "password123"}`, testUsername)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	found := false

	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			found = true

			assert.True(t, cookie.HttpOnly, "Cookie should be HttpOnly")
			assert.Equal(t, "/", cookie.Path, "Cookie path should be /")
			assert.NotEmpty(t, cookie.Value, "Cookie value should not be empty")
		}
	}

	assert.True(t, found, "auth_token cookie should be present")

	var resp user.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.NotNil(t, resp.User, "Response should contains user")
	assert.Equal(t, testUsername, resp.User.Username)
}

func TestLogout(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	_, token, _ := setupJoinedUser(ctx, t, router, 10000.0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/sessions", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			found = true
			assert.True(t, cookie.MaxAge < 0 || cookie.Value == "")
		}
	}
	assert.True(t, found, "auth_token cookie should be returned in response to clear it")

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Logged out successfully", resp["message"])
}
