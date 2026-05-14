package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/enterprise-erp/core/internal/domain"
	"github.com/enterprise-erp/core/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo   *postgres.UserRepository
	jwtSecret  string
	jwtExpTime time.Duration
}

func NewAuthUseCase(userRepo *postgres.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		jwtExpTime: time.Hour * 24, // 1 day
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (u *AuthUseCase) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUseCase) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(u.jwtExpTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
