package service

import (
	"errors"
	"time"

	"dot-be-go/internal/domain/entity"
	"dot-be-go/internal/domain/repository"
	"dot-be-go/pkg/hash"
	"dot-be-go/pkg/jwt"
)

// AuthRequest represents authentication request data
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name,omitempty" validate:"omitempty,min=3"`
}

// AuthResponse represents authentication response data
type AuthResponse struct {
	Token    string       `json:"token"`
	User     *entity.User `json:"user"`
	ExpireAt time.Time    `json:"expire_at"`
}

// AuthService handles authentication operations
type AuthService interface {
	Register(req *AuthRequest) (*AuthResponse, error)
	Login(req *AuthRequest) (*AuthResponse, error)
	GetUserByID(id uint) (*entity.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register creates a new user and returns auth response
func (s *authService) Register(req *AuthRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil || existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := hash.GenerateHash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := jwt.GenerateToken(user, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:    token,
		User:     user,
		ExpireAt: time.Now().Add(s.jwtExpiry),
	}, nil
}

// Login authenticates a user and returns auth response
func (s *authService) Login(req *AuthRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !hash.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := jwt.GenerateToken(user, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:    token,
		User:     user,
		ExpireAt: time.Now().Add(s.jwtExpiry),
	}, nil
}

// GetUserByID returns a user by ID
func (s *authService) GetUserByID(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
}
