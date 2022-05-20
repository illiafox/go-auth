package repository

import (
	"auth-example/utils/config"
	"auth-example/utils/templates"
	"go.uber.org/zap"
)

type Model struct {
	Log    *zap.Logger
	Rep    Repository
	TS     *templates.Templates
	Config config.Host
}

func NewModel(log *zap.Logger, rep Repository, ts *templates.Templates, conf config.Host) Model {
	return Model{log, rep, ts, conf}
}
