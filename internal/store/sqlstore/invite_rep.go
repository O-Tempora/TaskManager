package sqlstore

import (
	"dip/internal/models"
	"time"
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

func (r *InviteRep) Decline(invite_id int) error {
	if err := r.Delete(invite_id); err != nil {
		return err
	}
	return nil
}

func (r *InviteRep) Accept(invite_id, ws_id, usr_id int) (*models.WorkspaceJoined, error) {
	role_id, err := r.store.Role().GetIdByName("User")
	if err != nil {
		return nil, err
	}
	res, err := r.store.db.Exec(`insert into person_workspace
		(person_id, workspace_id, role_id)
		values($1, $2, $3) on conflict do nothing`,
		usr_id, ws_id, role_id,
	)
	if err != nil {
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if err := r.Delete(invite_id); err != nil {
		return nil, err
	}

	w := &models.WorkspaceJoined{}
	if err := r.store.db.QueryRow(`select pw.workspace_id, w.name, w.description, w.created_at, w.isactive, w.closed_at, ur.name
		from person_workspace as pw 
		join persons as p on p.id = pw.person_id 
		join user_role as ur on ur.id = pw.role_id 
		join workspaces as w on w.id = pw.workspace_id 
		where w.id = $1 `,
		ws_id).Scan(&w.Id, &w.Name, &w.Description, &w.CreatedAt, &w.IsActive, &w.ClosedAt, &w.Role); err != nil {
		return nil, err
	}

	return w, nil
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

func (r *InviteRep) Send(email string, ws_id, user_id int) error {
	p, err := r.store.Person().GetByEmail(email)
	if err != nil {
		return err
	}
	err = r.Create(&models.Invite{
		Id:         p.Id,
		SenderId:   user_id,
		ReceiverId: p.Id,
		WsId:       ws_id,
		SentAt:     time.Now(),
	})
	return err
}
