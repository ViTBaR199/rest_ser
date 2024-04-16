package models

import "time"

// определение модели пользователя и методы работы с ней
// структуры данных, которые используются в приложении

type User struct {
	ID               int       `json:"id"`
	Login            string    `json:"login"`
	Email            string    `json:"email"`
	Password         string    `json:"password,omitempty"` // предполагается, что пароль не должен сериализовываться
	RegistrationTime time.Time `json:"registration_time"`
	LastUpdateTime   time.Time `json:"last_update_time"`
}
