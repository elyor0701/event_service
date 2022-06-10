package storage

import (
	"event/genproto"

	"github.com/jmoiron/sqlx"
)

type ConnectSql struct {
	db *sqlx.DB
}

func (c *ConnectSql) Push(req genproto.Event) (genproto.Event, error) {
	_, err := c.db.Exec(`
		INSERT INTO event_time(time, event)
		VALUES ($1, $2) returning event`,
		req.Time,
		req.Event,
	)

	if err != nil {
		return genproto.Event{}, err
	}

	return req, nil
}

func (c *ConnectSql) Get() ([]*genproto.Event, error) {
	rows, err := c.db.Queryx(
		`SELECT * FROM event_time`,
	)

	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*genproto.Event

	for rows.Next() {
		var event genproto.Event

		err := rows.Scan(&event.Time, &event.Event)

		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}

func (c *ConnectSql) GetByTime(req genproto.Time) ([]*genproto.Event, error) {
	rows, err := c.db.Queryx(
		`SELECT * FROM event_time
		WHERE time=$1
		`,
		req.Time,
	)

	if err != nil {
		return nil, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*genproto.Event

	for rows.Next() {
		var event genproto.Event

		err := rows.Scan(&event.Time, &event.Event)

		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}
