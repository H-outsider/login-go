package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"login/api"
	"login/internal/auth"
	"login/internal/data/model"
	"login/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	user *model.User
}

func (r *fakeUserRepository) FindUserByUsername(username string) (*model.User, error) {
	return r.user, nil
}

func (r *fakeUserRepository) CreateUser(user *model.User) error {
	if r.user != nil {
		return service.ErrUserExists
	}
	r.user = user
	return nil
}

func TestRegisterReturnsConflictWhenUserExists(t *testing.T) {
	h := newTestHandler(&fakeUserRepository{user: &model.User{ID: 1, Username: "alice"}})
	router := gin.New()
	router.POST("/register", h.Register)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"alice","password":"secret123"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusConflict)
	}

	var resp api.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal response error = %v", err)
	}
	if resp.Code != 409 || resp.Message != service.ErrUserExists.Error() {
		t.Fatalf("response = %+v", resp)
	}
}

func TestLoginReturnsUniformSuccessResponse(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	h := newTestHandler(&fakeUserRepository{
		user: &model.User{
			ID:       7,
			Username: "alice",
			Password: string(passwordHash),
			Email:    "alice@example.com",
		},
	})
	router := gin.New()
	router.POST("/login", h.Login)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"alice","password":"secret123"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var resp api.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal response error = %v", err)
	}
	if resp.Code != 0 || resp.Message != "登录成功" || resp.Data == nil {
		t.Fatalf("response = %+v", resp)
	}
}

func newTestHandler(repo service.UserRepository) *UserHandler {
	gin.SetMode(gin.TestMode)
	userService := service.NewUserService(repo)
	jwtManager := auth.NewJWTManager("test-secret", "test", time.Hour)
	return NewUserHandler(userService, jwtManager)
}
