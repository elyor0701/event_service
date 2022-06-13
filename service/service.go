package service

import (
	"context"
	"database/sql"
	"event/genproto"
	"event/logger"
	"event/storage"

	"github.com/gofrs/uuid"
)

type EventService struct {
	csql storage.ToDo
	log  logger.Logger
}

func NewEventService(csql storage.ToDo, log logger.Logger) *EventService {
	return &EventService{
		csql: csql,
		log:  log,
	}
}

func (e *EventService) Push(c context.Context, req *genproto.Event) (*genproto.Event, error) {

	id, err := uuid.NewV4()

	if err != nil {
		e.log.Error("cant get id uuid", logger.Error(err))
		return nil, err
	}

	req.Id = id.String() // check for unique

	res, err := e.csql.ToDo().Push(*req)

	if err != nil {
		e.log.Error("PUSH", logger.Error(err))
		return nil, err
	}
	//fmt.Println("Push", res)
	e.log.Info("PUSHED new data : ", logger.String("Time", res.GetTime()), logger.String("Event", res.GetEvent()))
	return &res, nil
}

func (e *EventService) Get(c context.Context, req *genproto.Empty) (*genproto.Events, error) {
	res, err := e.csql.ToDo().Get()

	if err != nil {
		e.log.Error("PUSH", logger.Error(err))
		return nil, err
	}

	return &genproto.Events{Events: res, Count: int32(len(res))}, nil
}

func (e *EventService) GetByTime(c context.Context, req *genproto.Time) (*genproto.Events, error) {
	res, err := e.csql.ToDo().GetByTime(*req)

	if err != nil {
		e.log.Error("PUSH", logger.Error(err))
		return nil, err
	}

	return &genproto.Events{Events: res, Count: int32(len(res))}, nil
}

func (e *EventService) GetByID(c context.Context, req *genproto.Id) (*genproto.Event, error) {
	res, err := e.csql.ToDo().GetByID(*req)

	if err == sql.ErrNoRows {
		return &genproto.Event{}, nil
	}

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (e *EventService) UpdateEvent(c context.Context, req *genproto.Event) (*genproto.Event, error) {
	res, err := e.csql.ToDo().UpdateEvent(*req)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (e *EventService) DeleteEvent(c context.Context, req *genproto.Id) (*genproto.Empty, error) {
	err := e.csql.ToDo().DeleteEvent(*req)

	if err != nil {
		return &genproto.Empty{}, err
	}

	return &genproto.Empty{}, nil
}
