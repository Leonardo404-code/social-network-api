package repository

import (
	"api/src/models"
	"database/sql"
)

type Publication struct {
	db *sql.DB
}

func NewPublicationRepository(db *sql.DB) *Publication {
	return &Publication{db}
}

func (repository Publication) CreatePost(publication models.Publication) (uint64, error) {
	statement, erro := repository.db.Prepare(
		"insert into publications (title, content, autor_id) VALUES (?, ?, ?)",
	)

	if erro != nil {
		return 0, erro
	}

	defer statement.Close()

	result, erro := statement.Exec(
		publication.Titulo, publication.Conteudo, publication.AutorID,
	)

	if erro != nil {
		return 0, erro
	}

	lastInsertID, erro := result.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(lastInsertID), nil
}

func (repository Publication) SearchPostByID(publicationID uint64) (models.Publication, error) {
	line, erro := repository.db.Query(`
	SELECT p.*, u.nick FROM publications p INNER JOIN users u ON u.id = p.autor_id WHERE
	p.id = ?`, publicationID)

	if erro != nil {
		return models.Publication{}, erro
	}

	defer line.Close()

	var publication models.Publication

	if line.Next() {
		if erro = line.Scan(
			&publication.ID,
			&publication.Titulo,
			&publication.Conteudo,
			&publication.AutorID,
			&publication.Likes,
			&publication.Created_at,
			&publication.AutorNick,
		); erro != nil {
			return models.Publication{}, erro
		}
	}

	return publication, nil
}

func (repository Publication) SearchPost(userID uint64) ([]models.Publication, error) {
	lines, erro := repository.db.Query(`
	select distinct p.*, u.nick from publications p 
	inner join users u on u.id = p.autor_id 
	inner join followers s on p.autor_id = s.user_id 
	where u.id = ? or s.follower_id = ?
	order by 1 desc
	`, userID, userID,
	)

	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var publications []models.Publication

	if lines.Next() {
		var publication models.Publication

		if erro = lines.Scan(
			&publication.ID,
			&publication.Titulo,
			&publication.Conteudo,
			&publication.AutorID,
			&publication.Likes,
			&publication.Created_at,
			&publication.AutorNick,
		); erro != nil {
			return nil, erro
		}

		publications = append(publications, publication)
	}

	return publications, nil
}

func (repository Publication) UpdatePost(publicationID uint64, publication models.Publication) error {
	statement, erro := repository.db.Prepare(`
		update publications set title = ?, content = ? where id = ?
	`)

	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(publication.Titulo, publication.Conteudo, publicationID); erro != nil {
		return erro
	}

	return nil
}

func (repository Publication) DeletePost(publicationID uint64) error {
	statement, erro := repository.db.Prepare("DELETE FROM publications WHERE id = ?")

	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(publicationID); erro != nil {
		return erro
	}

	return nil
}

func (repository Publication) SearchUserPublications(usuarioID uint64) ([]models.Publication, error) {
	linhas, erro := repository.db.Query(`
		select p.*, u.nick from publications p
		join users u on u.id = p.autor_id
		where p.autor_id = ?`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []models.Publication

	for linhas.Next() {
		var publication models.Publication

		if erro = linhas.Scan(
			&publication.ID,
			&publication.Titulo,
			&publication.Conteudo,
			&publication.AutorID,
			&publication.Likes,
			&publication.Created_at,
			&publication.AutorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publication)
	}

	return publicacoes, nil
}

func (repository Publication) LikePost(publicationID uint64) error {
	statement, erro := repository.db.Prepare(
		`UPDATE publications SET likes = likes + 1 WHERE id = ?`,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicationID); erro != nil {
		return erro
	}

	return nil
}

func (repository Publication) UnlikePost(publicationID uint64) error {
	statement, erro := repository.db.Prepare(`
	UPDATE publications SET likes = CASE WHEN likes > 0 THEN likes -1 ELSE 0 END
	WHERE id = ?
	`)

	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(publicationID); erro != nil {
		return erro
	}

	return nil
}
