package handlers

import (
	"dip/internal/models"
	"dip/internal/store"
	"strconv"
)

func GetHome(store store.Store, ids string) (*models.HomePage, int, error) {
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, 422, err
	}

	ws, err := store.Workspace().GetByUser(id)
	if err != nil {
		return nil, 404, err
	}

	return ws, 200, nil
}
