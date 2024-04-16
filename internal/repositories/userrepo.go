package repositories

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"errors"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

// непосредственное взаимодействие с базой данных
// реализация CRUD (создание, чтение, обновление, удаление)

// описывает интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	AuthenticateUser(ctx context.Context, login, password string) (int, string, bool, error)
}

// реализует интерфейс UserRepository и предоставляет доступ к базе данных.
type userRepository struct {
	db *sql.DB
}

// создает новый экземпляр userRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// использует SQL-функцию create_new_user для создания нового пользователя.
func (r *userRepository) CreateUser(ctx context.Context, user models.User) error {
	encryptedPassword, err := encryptMessage("8j#k9P$3m^L5&@7*1q!6zG2blwMbT8!2", user.Password)
	if err != nil {
		return err
	}

	user.Password = encryptedPassword
	_, err = r.db.ExecContext(ctx, "SELECT create_new_user($1, $2, $3)", user.Login, user.Email, user.Password)
	return err
}

func (r *userRepository) AuthenticateUser(ctx context.Context, login, password string) (int, string, bool, error) {
	var userID int
	var userEmail string
	var isAuthenticated bool

	encryptedPassword, err := encryptMessage("8j#k9P$3m^L5&@7*1q!6zG2blwMbT8!2", password)
	if err != nil {
		return 0, "", false, err
	}

	// Вызов функции authenticate_user из PostgreSQL
	err = r.db.QueryRowContext(ctx, "SELECT * FROM authenticate_user($1, $2)", login, encryptedPassword).Scan(&userID, &userEmail, &isAuthenticated)
	if err != nil {
		return 0, "", false, err
	}
	return userID, userEmail, isAuthenticated, nil
}

func encryptMessage(key string, message string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("key must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	staticIV := []byte("1234567890123456")

	plaintext := []byte(message)
	// Добавляем дополнение к сообщению
	plaintext = applyPKCS7Padding(plaintext, block.BlockSize())

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, staticIV)
	mode.CryptBlocks(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// applyPKCS7Padding добавляет дополнение к блоку данных в соответствии с PKCS#7.
func applyPKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
