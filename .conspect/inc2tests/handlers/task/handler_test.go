package task_test

import (
	"encoding/json"
	"github.com/AtCliffUnderline/url-shortener/.conspect/inc2tests/handlers/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserViewHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		user        task.User
	}
	tests := []struct {
		name    string
		request string
		users   map[string]task.User
		want    want
	}{
		{
			name: "simple test #1",
			users: map[string]task.User{
				"id1": {
					ID:        "id1",
					FirstName: "Misha",
					LastName:  "Popov",
				},
			},
			want: want{
				contentType: "application/json",
				statusCode:  200,
				user:        task.User{},
			},
			request: "/users?user_id=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(task.UserViewHandler(tt.users))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			var user task.User
			err = json.Unmarshal(userResult, &user)
			require.NoError(t, err)

			assert.Equal(t, tt.want.user, user)
		})
	}
}
