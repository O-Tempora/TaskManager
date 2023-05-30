package sqlstore

import (
	"dip/internal/models"
)

type CommentRep struct {
	store *Store
}

func (r *CommentRep) GetByTask(id int) ([]models.CommentShow, error) {
	comments := make([]models.CommentShow, 0)
	c := models.CommentShow{}

	rows, err := r.store.db.Query(`select m.id, m.content, m.created_at, m.person_id from commentaries m where m.task_id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&c.Id, &c.Content, &c.CreatedAt, &c.PersonId)
		if err != nil {
			return nil, err
		}
		name, err := r.store.Person().GetNameById(c.PersonId)
		if err != nil {
			return nil, err
		}
		c.Person = name
		comments = append(comments, c)
	}
	return comments, nil
}
func (r *CommentRep) Create(c models.Comment) (*models.CommentShow, error) {
	comment := &models.CommentShow{
		Content:   c.Content,
		PersonId:  c.PersonId,
		CreatedAt: c.CreatedAt,
	}
	if err := r.store.db.QueryRow(`insert into commentaries
		(content, created_at, person_id, task_id)
		values ($1, $2, $3, $4) returning id`,
		c.Content, c.CreatedAt, c.PersonId, c.TaskId,
	).Scan(&comment.Id); err != nil {
		return nil, err
	}
	name, err := r.store.Person().GetNameById(comment.PersonId)
	if err != nil {
		return nil, err
	}
	comment.Person = name
	return comment, nil
}
func (r *CommentRep) Delete(id int) error {
	res, err := r.store.db.Exec(`delete from commentaries where id = $1`, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}
