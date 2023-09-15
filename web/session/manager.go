package session

import (
	"github.com/google/uuid"
	"myProject/web"
)

type Manager struct {
	Store
	Propagator
	CtxSessionKey string
}

func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	newSessionID := uuid.New().String()
	session, err := m.Generate(ctx.Req.Context(), newSessionID)
	if err != nil {
		return nil, err
	}
	err = m.Inject(ctx.Response, newSessionID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.CacheSession == nil {
		ctx.CacheSession = make(map[string]any, 1)
	}
	if session, ok := ctx.CacheSession[m.CtxSessionKey]; ok {
		return session.(Session), nil
	}
	sessionId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	session, err := m.Get(ctx.Req.Context(), sessionId)
	if err != nil {
		return nil, err
	}
	ctx.CacheSession[m.CtxSessionKey] = session
	return session, nil
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Response, session.ID())
}

func (m *Manager) RefreshSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), session.ID())
}
