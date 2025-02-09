package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"hr-app-diploma-backend/auth-service/config"
	"hr-app-diploma-backend/auth-service/models"
	"hr-app-diploma-backend/auth-service/repository"
	"hr-app-diploma-backend/auth-service/utils"
	"hr-app-diploma-backend/pkg/logger"
	"log/slog"
	"time"
)

type AuthService struct {
	repo      *repository.AuthRepository
	JWTSecret string
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo:      repo,
		JWTSecret: config.AppConfig.JWT.Secret,
	}
}

func (s *AuthService) RegisterUser(email, password, firstName, lastName string) error {
	logger.Log.Info("Starting user registration",
		slog.String("email", email),
		slog.String("Name", firstName))

	exists, err := s.repo.UserExistsByEmail(email)
	if err != nil {
		logger.Log.Error("Error checking user existence",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return err
	}
	if exists {
		logger.Log.Warn("Registration attempt for existing user",
			slog.String("email", email))
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Log.Error("Failed to hash password", slog.String("error", err.Error()))
		return err
	}

	user := &models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  hashedPassword,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		logger.Log.Error("Failed to create user",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("User registered successfully",
		slog.String("email", email),
		slog.String("name", firstName))

	return nil
}

func (s *AuthService) LoginUser(email, password string) (string, error) {
	logger.Log.Info("Attempting user login", slog.String("email", email))

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		logger.Log.Warn("Login failed: user not found", slog.String("email", email))
		return "", errors.New("invalid credentials")
	}
	logger.Log.Info("Comparing passwords",
		slog.String("input_password", password),
		slog.String("hashed_password", user.Password))

	if !utils.CheckPassword(password, user.Password) {
		logger.Log.Warn("Login failed: incorrect password", slog.String("email", email))
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		logger.Log.Error("Failed to generate JWT", slog.String("error", err.Error()))
		return "", err
	}

	logger.Log.Info("User logged in successfully", slog.String("email", email))

	return token, nil
}

func (s *AuthService) VerifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return 0, errors.New("token expired")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("missing user_id")
	}

	return int(userID), nil
}
