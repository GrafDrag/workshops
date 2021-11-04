package handler

import (
	"calendar/internal/app"
	"calendar/internal/controller"
	"calendar/internal/model"
	"calendar/pb"
	"golang.org/x/net/context"
)

type EventHandler struct {
	Server *app.IServer
}

func (e *EventHandler) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	searchModel := model.SearchEvent{
		Title:    request.GetTitle(),
		Timezone: request.GetTimezone(),
		DateFrom: request.GetDateFrom(),
		DateTo:   request.GetDateTo(),
		TimeFrom: request.GetTimeFrom(),
		TimeTo:   request.GetTimeTo(),
	}

	c := controller.NewEventController(e.Server.Store)
	data, err := c.List(ctx, searchModel)

	if err != nil {
		return nil, err
	}

	res := &pb.ListResponse{}
	for _, event := range data {
		item := &pb.Event{
			Id:          int32(event.ID),
			Title:       event.Title,
			Description: event.Description,
			Time:        event.Time,
			Timezone:    event.Timezone,
			Duration:    event.Duration,
			Notes:       event.Notes,
		}
		res.Event = append(res.Event, item)
	}

	return res, nil
}

func (e *EventHandler) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	c := controller.NewEventController(e.Server.Store)
	event := &model.Event{
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Time:        request.GetTime(),
		Timezone:    request.GetTimezone(),
		Duration:    request.GetDuration(),
		Notes:       request.GetNotes(),
	}

	if err := c.Create(ctx, event); err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Status: pb.CreateResponse_Successful,
	}, nil
}

func (e *EventHandler) GetById(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	c := controller.NewEventController(e.Server.Store)
	event, err := c.FindById(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		Event: &pb.Event{
			Id:          int32(event.ID),
			Title:       event.Title,
			Description: event.Description,
			Time:        event.Time,
			Timezone:    event.Timezone,
			Duration:    event.Duration,
			Notes:       event.Notes,
		},
	}, nil
}

func (e *EventHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	c := controller.NewEventController(e.Server.Store)
	event := &model.Event{
		ID:          int(request.GetId()),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Time:        request.GetTime(),
		Timezone:    request.GetTimezone(),
		Duration:    request.GetDuration(),
		Notes:       request.GetNotes(),
	}

	if err := c.Update(ctx, event); err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{
		Status: pb.UpdateResponse_Successful,
	}, nil
}

func (e *EventHandler) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	c := controller.NewEventController(e.Server.Store)

	if err := c.Delete(ctx, request.GetId()); err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{
		Status: pb.DeleteResponse_Successful,
	}, nil
}
