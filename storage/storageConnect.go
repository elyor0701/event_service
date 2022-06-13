package storage

import (
	"event/genproto"

	"github.com/jmoiron/sqlx"
)

func NewEventConnectSql(db *sqlx.DB) *storageEvent {
	return &storageEvent{
		db:            db,
		eventCommands: &ConnectSql{db: db},
	}
}

type SqlQuerysEvent interface {
	Push(genproto.Event) (genproto.Event, error)
	Get() ([]*genproto.Event, error)
	GetByTime(genproto.Time) ([]*genproto.Event, error)
	GetByID(genproto.Id) (genproto.Event, error)
	UpdateEvent(genproto.Event) (genproto.Event, error)
	DeleteEvent(genproto.Id) error
}

type ToDo interface {
	ToDo() SqlQuerysEvent
}

type storageEvent struct {
	db            *sqlx.DB
	eventCommands SqlQuerysEvent
}

func (s storageEvent) ToDo() SqlQuerysEvent {
	return s.eventCommands
}
