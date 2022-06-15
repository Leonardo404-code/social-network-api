package models

import (
	"errors"
	"strings"
	"time"
)

// Publication represent a post of user
type Publication struct {
	ID         uint64    `json:"id,omitempty"`
	Titulo     string    `json:"titulo,omitempty"`
	Conteudo   string    `json:"conteudo,omitempty"`
	AutorID    uint64    `json:"autorId,omitempty"`
	AutorNick  string    `json:"autorNick,omitempty"`
	Likes      uint64    `json:"likes,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
}

// Prepare call the methods for validation all filds in publication
func (publication *Publication) Preparar() error {
	if erro := publication.validar(); erro != nil {
		return erro
	}

	publication.formatar()
	return nil
}

func (publication *Publication) validar() error {
	if publication.Titulo == "" {
		return errors.New("o título é obrigatório e não pode estar em branco")
	}

	if publication.Conteudo == "" {
		return errors.New("o conteúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

func (publication *Publication) formatar() {
	publication.Titulo = strings.TrimSpace(publication.Titulo)
	publication.Conteudo = strings.TrimSpace(publication.Conteudo)
}
