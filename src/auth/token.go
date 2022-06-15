package auth

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateToken return a token assigned with permissions of user
func CreateToken(usuarioID uint64) (string, error) {
	permissions := jwt.MapClaims{} // Map of permisions of user
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix() // Time of expiration
	permissions["usuarioId"] = usuarioID                      // Add the token to UserID
	// Create a token with permisions
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString(config.SecretKey)
}

// TokenValidation Verify if token in request is valid
func TokenValidation(r *http.Request) error {
	tokenString := extractToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token inválido")
}

func ExtractUserID(r *http.Request) (uint64, error) {
	tokenString := extractToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if erro != nil {
		return 0, erro
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, erro := strconv.ParseUint(
			fmt.Sprintf("%.0f", permissions["usuarioId"]), 10, 64,
		)
		if erro != nil {
			return 0, erro
		}

		return userID, nil
	}
	return 0, errors.New("token inválido")
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Metodo de assinatura inesperado! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
