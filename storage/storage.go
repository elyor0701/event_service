package storage

import (
	"database/sql"
	"event/genproto"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConnectSql struct {
	db *sqlx.DB
}

func (c *ConnectSql) Push(req genproto.Event) (genproto.Event, error) {
	var newsId string
	err := c.db.QueryRow(`
		INSERT INTO event_time(id, time, event, status)
		VALUES ($1, $2, $3, $4) returning id`,
		req.GetId(),
		req.GetTime(),
		req.GetEvent(),
		req.GetStatus(),
	).Scan(&newsId)

	if err != nil {
		return genproto.Event{}, err
	}

	if newsId != req.Id {
		return genproto.Event{}, fmt.Errorf("not match id")
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

		err := rows.Scan(&event.Time, &event.Event, &event.Id, &event.Status)

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
		req.GetTime(),
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

		err := rows.Scan(&event.Time, &event.Event, &event.Id, &event.Status)

		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}

func (c *ConnectSql) GetByID(req genproto.Id) (genproto.Event, error) {
	var event genproto.Event
	var (
		time   sql.NullString
		eventt sql.NullString
		id     sql.NullString
		status sql.NullBool
	)
	err := c.db.QueryRow( // how to get one row
		`SELECT * FROM event_time
		WHERE id=$1
		`,
		req.GetId(),
	).Scan(&time, &eventt, &id, &status)

	if err != nil {
		return genproto.Event{}, err
	}

	event.Time = time.String
	event.Event = eventt.String
	event.Id = id.String
	event.Status = status.Bool

	return event, nil
}

func (c *ConnectSql) UpdateEvent(req genproto.Event) (genproto.Event, error) {
	var event genproto.Event
	err := c.db.QueryRow(
		`UPDATE event_time
		SET id = $1, time = $2, event = $3, status = $4
		WHERE id = $1 returning id, time, event, status`,
		req.GetId(),
		req.GetTime(),
		req.GetEvent(),
		req.GetStatus(),
	).Scan(&event.Id, &event.Time, &event.Event, &event.Status)

	if err != nil {
		return genproto.Event{}, err
	}

	if event.GetTime() != req.GetTime() || event.GetEvent() != req.GetEvent() || event.GetStatus() != req.Status {
		return genproto.Event{}, fmt.Errorf("mismatch new and pushed data")
	}

	return event, nil
}

func (c *ConnectSql) DeleteEvent(req genproto.Id) error {
	_, err := c.db.Exec(
		`DELETE FROM event_time
		WHERE id = $1`,
		req.GetId(),
	)

	if err != nil {
		return fmt.Errorf("cant delete from sql")
	}

	return nil
}
