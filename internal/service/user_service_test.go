package service

import (
	"errors"
	"testing"

	"login/api"
	"login/internal/data"
	"login/internal/data/model"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	user      *model.User
	findErr   error
	createErr error
	created   *model.User
}

func (r *fakeUserRepository) FindUserByUsername(username string) (*model.User, error) {
	return r.user, r.findErr
}

func (r *fakeUserRepository) CreateUser(user *model.User) error {
	r.created = user
	return r.createErr
}

func TestRegisterServiceCreatesUserWithHashedPassword(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := NewUserService(repo)

	req := api.RegisterRequest{
		Username: "alice",
		Password: "secret123",
		Email:    "alice@example.com",
	}

	if err := svc.RegisterService(req); err != nil {
		t.Fatalf("RegisterService() error = %v", err)
	}
	if repo.created == nil {
		t.Fatal("RegisterService() did not create user")
	}
	if repo.created.Password == req.Password {
		t.Fatal("RegisterService() stored plaintext password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.created.Password), []byte(req.Password)); err != nil {
		t.Fatalf("stored password hash does not match password: %v", err)
	}
}

func TestRegisterServiceReturnsUserExists(t *testing.T) {
	repo := &fakeUserRepository{user: &model.User{ID: 1, Username: "alice"}}
	svc := NewUserService(repo)

	err := svc.RegisterService(api.RegisterRequest{Username: "alice", Password: "secret123"})
	if !errors.Is(err, ErrUserExists) {
		t.Fatalf("RegisterService() error = %v, want %v", err, ErrUserExists)
	}
}

func TestRegisterServiceMapsDuplicateKeyToUserExists(t *testing.T) {
	repo := &fakeUserRepository{createErr: data.ErrDuplicateKey}
	svc := NewUserService(repo)

	err := svc.RegisterService(api.RegisterRequest{Username: "alice", Password: "secret123"})
	if !errors.Is(err, ErrUserExists) {
		t.Fatalf("RegisterService() error = %v, want %v", err, ErrUserExists)
	}
}

func TestLoginServiceReturnsInvalidCredentialsForMissingUser(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := NewUserService(repo)

	_, err := svc.LoginService(api.LoginRequest{Username: "alice", Password: "secret123"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("LoginService() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestLoginServiceReturnsUserResponse(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &fakeUserRepository{
		user: &model.User{
			ID:       7,
			Username: "alice",
			Password: string(passwordHash),
			Email:    "alice@example.com",
		},
	}
	svc := NewUserService(repo)

	resp, err := svc.LoginService(api.LoginRequest{Username: "alice", Password: "secret123"})
	if err != nil {
		t.Fatalf("LoginService() error = %v", err)
	}
	if resp.ID != 7 || resp.Username != "alice" || resp.Email != "alice@example.com" {
		t.Fatalf("LoginService() response = %+v", resp)
	}
}
