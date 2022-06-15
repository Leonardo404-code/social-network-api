package models

import (
	security "api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario cria modelo de usuário da aplicação
type User struct {
	ID         uint64    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Nick       string    `json:"nick,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
}

func (user *User) Prepare(etapa string) error {

	if erro := user.validation(etapa); erro != nil {
		return erro
	}

	if erro := user.format(etapa); erro != nil {
		return erro
	}

	return nil
}

func (user *User) validation(etapa string) error {
	if user.Name == "" {
		return errors.New("o nome é obrigatório e não pode estar em branco")
	}

	if user.Nick == "" {
		return errors.New("o nick é obrigatório e não pode estar em branco")
	}

	if user.Email == "" {
		return errors.New("o e-mail é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("o e-mail inserido é inválido")
	}

	if etapa == "cadastro" && user.Password == "" {
		return errors.New("o senha é obrigatório e não pode estar em branco")
	}

	return nil
}

func (user *User) format(etapa string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if etapa == "cadastro" {
		passwordHash, erro := security.Hash(user.Password)
		if erro != nil {
			return erro
		}

		user.Password = string(passwordHash)
	}

	return nil
}
