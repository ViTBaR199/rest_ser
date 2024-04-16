package services

// реализуется бизнес-логика приложения
// доп. логика (валидация)

import (
	"context"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repositories"
	"regexp"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// validateUserData выполняет базовую валидацию данных пользователя.
func validateUserData(user models.User) error {
	if user.Login == "" || user.Email == "" || user.Password == "" {
		return errors.New("login, email and password must be provided")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !isValidEmail(user.Email) {
		return errors.New("email is not valid")
	}
	return nil
}

// isValidEmail проверяет, является ли строка валидной почтой.
func isValidEmail(email string) bool {
	// Простая регулярка для валидации email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) error {
	// Валидация данных пользователя
	if err := validateUserData(user); err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Authenticate(login, password string) (int, string, bool, error) {
	return s.repo.AuthenticateUser(context.Background(), login, password)
}
