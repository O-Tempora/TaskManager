package sqlstore

import (
	"dip/internal/models"
)

type InviteRep struct {
	store *Store
}

func (r *InviteRep) GetAll(user_id int) ([]models.InviteShow, error) {
	res := make([]models.InviteShow, 0)
	buf := models.InviteShow{}
	var sender_id int

	rows, err := r.store.db.Query(`select * from invites i where i.receiver_id = $1`, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&buf.Id, &buf.ReceiverId, &sender_id, &buf.WorkspaceId, &buf.SentAt)
		if err != nil {
			return nil, err
		}
		sender, err := r.store.Person().GetNameById(sender_id)
		if err != nil {
			return nil, err
		}
		buf.Sender = sender
		if err := r.store.db.QueryRow(`select w.name, w.created_at, w.description from workspaces w where w.id = $1`,
			buf.WorkspaceId).Scan(&buf.Workspace, &buf.CreatedAt, &buf.Description); err != nil {
			return nil, err
		}
		res = append(res, buf)
	}
	return res, nil
}

func (r *InviteRep) Delete(invite_id int) error {
	res, err := r.store.db.Exec(`delete from invites where id = $1`, invite_id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func (r *InviteRep) Create(inv *models.Invite) error {
	var id int

	if err := r.store.db.QueryRow(`insert into invites
		(receiver_id, sender_id, workspace_id, sent_at)
		values ($1, $2, $3, $4) on conflict do nothing returning id`,
		inv.ReceiverId, inv.SenderId, inv.WsId, inv.SentAt,
	).Scan(&id); err != nil {
		return err
	}
	return nil
}
