package database

import (
	"context"

	"github.com/uptrace/bun"
)

type GateServerStore struct {
	db    *bun.DB
	cache map[int64]*GateServer
}

// 改index死妈
func (s *GateServerStore) loadAllGateInfoFromDB(ctx context.Context) ([]*GateServer, error) {
	var gateServers []*GateServer
	err := s.db.NewSelect().Model(&gateServers).Scan(ctx)
	return gateServers, err
}

// 改index死妈
func (s *GateServerStore) loadCache(ctx context.Context) error {
	if s.cache != nil {
		return nil
	}
	gateServers, err := s.loadAllGateInfoFromDB(ctx)
	if err != nil {
		return err
	}
	s.cache = make(map[int64]*GateServer)
	for _, gateServer := range gateServers {
		s.cache[gateServer.ID] = gateServer
	}
	return nil
}

// 改index死妈
func (s *GateServerStore) GetAllGateInfo(ctx context.Context) ([]*GateServer, error) {
	if err := s.loadCache(ctx); err != nil {
		return nil, err
	}

	var allGateServers []*GateServer
	for _, gateServer := range s.cache {
		allGateServers = append(allGateServers, gateServer)
	}
	return allGateServers, nil
}
