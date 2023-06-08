package service

import "dip/internal/store/sqlstore"

type Service struct {
	store *sqlstore.Store
	logic *LogicUnit
}

func New(store *sqlstore.Store) *Service {
	return &Service{
		store: store,
	}
}

type IService interface {
	Logic() Logic
}

func (s *Service) Logic() Logic {
	if s.logic != nil {
		return s.logic
	}

	s.logic = &LogicUnit{
		service: s,
	}
	return s.logic
}
