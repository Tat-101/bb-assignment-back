package rest_test

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/internal"
	"github.com/tat-101/bb-assignment-back/tools/seed/seed"
)

func TestMain(m *testing.M) {
	go func() {
		r := internal.SetupServer()
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	seed.SeedAdminUser()

	code := m.Run()

	os.Exit(code)
}

func TestGetVersion(t *testing.T) {
	client := resty.New()

	resp, err := client.R().Get("http://localhost:8080/version")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.JSONEq(t, `{"version":"v0"}`, string(resp.Body()))
}

func TestLogin_Fail(t *testing.T) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"email": "invalid@bb.com", "password": "123456"}`).
		Post("http://localhost:8080/auth/login")
	body := string(resp.Body())

	assert.Contains(t, body, "invalid credentials")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
}

func TestLogin_Success(t *testing.T) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"email": "admin@bb.com", "password": "123456"}`).
		Post("http://localhost:8080/auth/login")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGJiLmNvbSIsImV4cCI6MTcyNDA4NTUyNywiaWF0IjoxNzIzOTk5MTI3fQ.KEn2SR8iFKT3PE9352LekKuBQqk3MyaCpt6bUj-BmCs"}
	body := string(resp.Body())
	assert.Contains(t, body, "token")

	var result map[string]string
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	token := result["token"]
	os.Setenv("TOKEN", token)
}

func TestGetUsers(t *testing.T) {
	client := resty.New()

	token := os.Getenv("TOKEN")
	// fmt.Println("Token:", token)

	resp, err := client.R().SetHeader("Authorization", token).Get("http://localhost:8080/users")
	body := string(resp.Body())

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	var users []domain.User
	err = json.Unmarshal([]byte(body), &users)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(users), 1, "The length of users should be at least 1")

	found := false
	for _, user := range users {
		if user.Email == "admin@bb.com" {
			os.Setenv("ADMIN_ID", strconv.Itoa(int(user.ID)))
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find email 'admin@bb.com'")
}

func TestGetUsersWithoutToken(t *testing.T) {
	client := resty.New()

	resp, err := client.R().Get("http://localhost:8080/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
}

func TestGetUserByID(t *testing.T) {
	client := resty.New()

	token := os.Getenv("TOKEN")
	// fmt.Println("Token:", token)

	resp, err := client.R().SetHeader("Authorization", token).Get("http://localhost:8080/users/" + os.Getenv("ADMIN_ID"))
	body := string(resp.Body())

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var user domain.User
	err = json.Unmarshal([]byte(body), &user)
	assert.NoError(t, err)

	assert.Equal(t, "admin@bb.com", user.Email)
}

func TestCreateUser_Success(t *testing.T) {
	client := resty.New()

	token := os.Getenv("TOKEN")

	// Send a POST request to create a new user
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(`{"email": "test_no1@example.com", "password": "123456", "name": "test"}`).
		Post("http://localhost:8080/users")

	// Check the response status and content
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())

	// Verify that the user exists in the list of users
	usersResp, err := client.R().SetHeader("Authorization", token).Get("http://localhost:8080/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, usersResp.StatusCode())

	var users []domain.User
	err = json.Unmarshal([]byte(usersResp.Body()), &users)
	assert.NoError(t, err)

	found := false
	for _, user := range users {
		if user.Email == "test_no1@example.com" {
			os.Setenv("TEST_ID", strconv.Itoa(int(user.ID)))
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find email 'test_no1@example.com'")
}

func TestUpdateUser(t *testing.T) {
	client := resty.New()

	token := os.Getenv("TOKEN")
	userID := os.Getenv("TEST_ID")

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(`{"name": "updated_name"}`).
		Put("http://localhost:8080/users/" + userID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Verify that the user is updated
	resp, err = client.R().
		SetHeader("Authorization", token).
		Get("http://localhost:8080/users/" + userID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var user domain.User
	err = json.Unmarshal([]byte(resp.Body()), &user)
	assert.NoError(t, err)

	assert.Equal(t, "updated_name", user.Name)
}

func TestDeleteUser(t *testing.T) {
	client := resty.New()

	token := os.Getenv("TOKEN")
	userID := os.Getenv("TEST_ID")

	resp, err := client.R().
		SetHeader("Authorization", token).
		Delete("http://localhost:8080/users/" + userID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())

	// Verify that the user is deleted
	resp, err = client.R().
		SetHeader("Authorization", token).
		Get("http://localhost:8080/users/" + userID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
}
