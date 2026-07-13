package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLeaderboard(t *testing.T) {
	env := setupTestEnv(t)
	defer env.MiniRedis.Close()
	defer env.DB.Close()

	_, token, _ := env.setupJoinedUser(t, 10000.0)

	err := env.LeaderboardService.UpdateLeaderboard(ctx)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/leaderboard", nil)
	env.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rawResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &rawResp)
	assert.NoError(t, err)

	entries, ok := rawResp["entries"].([]interface{})
	assert.True(t, ok, "should have entries slice")
	assert.NotEmpty(t, entries, "should not be empty")

	entry, ok := entries[0].(map[string]interface{})
	assert.True(t, ok)

	userMap, ok := entry["user"].(map[string]interface{})
	assert.True(t, ok, "should have user profile map")

	assert.Equal(t, "Classified", userMap["username"])
	assertPublicProfilePrivacy(t, userMap)

	// Update user to public
	reqBody := `{"first_name": "Test", "last_name": "User", "website": "", "is_public": true}`
	wPatch := httptest.NewRecorder()
	reqPatch, _ := http.NewRequest(http.MethodPatch, "/api/v1/profile", bytes.NewBufferString(reqBody))
	reqPatch.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	env.Router.ServeHTTP(wPatch, reqPatch)
	assert.Equal(t, http.StatusOK, wPatch.Code)

	// Trigger leaderboard recalculation
	err = env.LeaderboardService.UpdateLeaderboard(ctx)
	assert.NoError(t, err)

	// Request Leaderboard
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/leaderboard", nil)
	env.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &rawResp)
	assert.NoError(t, err)

	entries, ok = rawResp["entries"].([]interface{})
	assert.True(t, ok, "should have entries slice")
	assert.NotEmpty(t, entries, "should not be empty")

	entry, ok = entries[0].(map[string]interface{})
	assert.True(t, ok)

	userMap, ok = entry["user"].(map[string]interface{})
	assert.True(t, ok, "should have user profile map")

	assert.Equal(t, testUsername, userMap["username"])
	assertPublicProfilePrivacy(t, userMap)
}
