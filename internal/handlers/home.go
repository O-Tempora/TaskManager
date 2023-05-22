package handlers

import (
	"dip/internal/models"
	"dip/internal/store"
	"strconv"
)

func GetHome(store store.Store, id int) (*models.HomePage, int, error) {
	// id, err := strconv.Atoi(ids)
	// if err != nil {
	// 	return nil, 422, err
	// }

	ws, err := store.Workspace().GetByUser(id)
	if err != nil {
		return nil, 404, err
	}

	return ws, 200, nil
}

func GetFullWorkspace(store store.Store, ws_id string) (*models.WorkspaceFull, int, error) {
	id, err := strconv.Atoi(ws_id)
	if err != nil {
		return nil, 422, err
	}

	fg := &models.FullGroup{}
	groups := make([]models.FullGroup, 0)

	tgs, err := store.TaskGroup().GetByWorkspaceId(id)
	if err != nil {
		return nil, 400, err
	}

	for _, v := range tgs {
		fg.Id = v.Id
		fg.Color = v.Color
		fg.Name = v.Name
		fg.Tasks, err = store.Task().GetAllByGroup(v.Id)
		if err != nil {
			return nil, 500, err
		}
		groups = append(groups, *fg)
	}

	wsf := &models.WorkspaceFull{
		Id:     id,
		Groups: groups,
	}
	return wsf, 200, nil
}
