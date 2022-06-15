package repository

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

// User Representa um repositório de usuário
type User struct {
	db *sql.DB
}

// Create a user repository
func NewRepositoryUser(db *sql.DB) *User {
	return &User{db}
}

// CreateUser create a user into the database
func (repository User) CreateUser(user models.User) (uint64, error) {
	statement, error := repository.db.Prepare(
		"INSERT INTO users (name, nick, email, password) VALUES(?, ?, ?, ?)",
	)

	if error != nil {
		return 0, error
	}

	defer statement.Close()

	result, error := statement.Exec(user.Name, user.Nick, user.Email, user.Password)

	if error != nil {
		return 0, error
	}

	ultimoIDInserido, erro := result.LastInsertId()

	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// SearchUsers is a function that searches users by their name or nickname
func (repositorio User) SearchUsers(nomeOuNick string) ([]models.User, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)

	linhas, erro := repositorio.db.Query(`SELECT id, name, nick, email, created_at FROM 
	users WHERE name LIKE ? or nick LIKE ?`, nomeOuNick, nomeOuNick)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var users []models.User

	for linhas.Next() {
		var user models.User

		if erro = linhas.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Created_at,
		); erro != nil {
			return nil, erro
		}

		users = append(users, user)
	}

	return users, nil
}

// SearchUserByID is a function that searches users by their ID
func (repositorio User) SearchUserByID(ID uint64) (models.User, error) {
	lines, erro := repositorio.db.Query(
		"SELECT id, name, nick, email, created_at FROM users WHERE id = ?", ID,
	)

	if erro != nil {
		return models.User{}, erro
	}

	defer lines.Close()

	var user models.User

	if lines.Next() {
		if erro = lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Created_at,
		); erro != nil {
			return models.User{}, erro
		}
	}

	return user, nil
}

// UpdateUser is a function for update the data from a user in database
func (repositorio User) UpdateUser(ID uint64, user models.User) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE users SET name = ?, nick = ?, email = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(user.Name, user.Nick, user.Email, ID); erro != nil {
		return erro
	}

	return nil
}

// DeleteUser is a function for delete a user in database
func (repositorio User) DeleteUser(ID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM users WHERE id = ?")
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// SearchEmail Search an user for email and return a id and password  with hash
func (repositorio User) SearchEmail(email string) (models.User, error) {
	line, erro := repositorio.db.Query(
		"SELECT id, password FROM users WHERE email = ?", email,
	)

	if erro != nil {
		return models.User{}, erro
	}

	defer line.Close()

	var user models.User

	if line.Next() {
		if erro = line.Scan(&user.ID, &user.Password); erro != nil {
			return models.User{}, erro
		}
	}

	return user, nil
}

// Follow add a follower to a user in database
func (repositorio User) Follow(userID, followerID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"INSERT IGNORE INTO followers (user_id, follower_id) values (?, ?)",
	)

	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(userID, followerID); erro != nil {
		return erro
	}

	return nil
}

// StopFollowing remove a follower to a user in database
func (repositorio User) StopFollowing(userID, followerID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM followers WHERE user_id = ? and follower_id = ?",
	)
	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(userID, followerID); erro != nil {
		return erro
	}

	return nil
}

// SearchFollowers search for user followers in database
func (repositorio User) SearchFollowers(userID uint64) ([]models.User, error) {
	lines, erro := repositorio.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.created_at FROM users u inner join
		followers s on u.id = s.follower_id WHERE s.user_id = ?`, userID)

	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var users []models.User
	for lines.Next() {
		var user models.User

		if erro = lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Created_at,
		); erro != nil {
			return nil, erro
		}
		users = append(users, user)
	}

	return users, nil
}

// SearchPassword search for a user password in database
func (repositorio User) SearchPassword(userID uint64) (string, error) {
	line, erro := repositorio.db.Query("SELECT password FROM users WHERE id = ?", userID)

	if erro != nil {
		return "", erro
	}

	defer line.Close()

	var user models.User

	if line.Next() {
		if erro = line.Scan(&user.Password); erro != nil {
			return "", erro
		}
	}

	return user.Password, nil
}

// UpdatePassword update user password in database
func (repositorio User) UpdatePassword(userID uint64, password string) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE users SET password = ? WHERE id = ?
	`)
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(password, userID); erro != nil {
		return erro
	}

	return nil
}
