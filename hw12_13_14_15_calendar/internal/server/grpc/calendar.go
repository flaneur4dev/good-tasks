package grpcserver

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	es "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/mistakes"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/grpc/pb"
)

type application interface {
	Events(ctx context.Context, date time.Time, period string) ([]cs.Event, error)
	CreateEvent(ctx context.Context, ne cs.Event) error
	UpdateEvent(ctx context.Context, id string, ne cs.Event) error
	DeleteEvent(ctx context.Context, id string) error
	Check(ctx context.Context) error
}

type calendarServer struct {
	pb.UnimplementedCalendarServer
	app application
}

func newCalendarServer(app application) *calendarServer {
	return &calendarServer{app: app}
}

func (s *calendarServer) Events(ctx context.Context, req *pb.EventsRequest) (*pb.EventsResponse, error) {
	evs, err := s.app.Events(ctx, req.GetDate().AsTime(), req.GetPeriod())
	if err != nil {
		sc := codes.Internal
		if errors.Is(err, es.ErrPeriod) {
			sc = codes.InvalidArgument
		}

		return nil, status.Error(sc, err.Error())
	}

	pbevs := make([]*pb.Event, 0, len(evs))
	for _, e := range evs {
		pbevs = append(pbevs, &pb.Event{
			Id:              e.ID,
			Title:           e.Title,
			Description:     e.Description,
			OwnerId:         e.OwnerID,
			StartDate:       timestamppb.New(e.StartDate),
			FinishDate:      timestamppb.New(e.FinishDate),
			NotificationDay: timestamppb.New(e.NotificationDay),
		})
	}

	return &pb.EventsResponse{Events: pbevs}, nil
}

func (s *calendarServer) CreateEvent(ctx context.Context, req *pb.Event) (*pb.EventResponse, error) {
	e := cs.Event{
		ID:              req.GetId(),
		Title:           req.GetTitle(),
		Description:     req.GetDescription(),
		OwnerID:         req.GetOwnerId(),
		StartDate:       req.GetStartDate().AsTime(),
		FinishDate:      req.GetFinishDate().AsTime(),
		NotificationDay: req.GetNotificationDay().AsTime(),
	}

	err := s.app.CreateEvent(ctx, e)
	if err != nil {
		sc := codes.Internal
		if errors.Is(err, es.ErrDateBusy) {
			sc = codes.InvalidArgument
		}

		return nil, status.Error(sc, err.Error())
	}

	return &pb.EventResponse{Id: e.ID, Message: "created"}, nil
}

func (s *calendarServer) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.EventResponse, error) {
	e := cs.Event{
		ID:              req.GetId(),
		Title:           req.GetTitle(),
		Description:     req.GetDescription(),
		OwnerID:         req.GetOwnerId(),
		StartDate:       req.GetStartDate().AsTime(),
		FinishDate:      req.GetFinishDate().AsTime(),
		NotificationDay: req.GetNotificationDay().AsTime(),
	}

	err := s.app.UpdateEvent(ctx, e.ID, e)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EventResponse{Id: e.ID, Message: "updated"}, nil
}

func (s *calendarServer) DeleteEvent(ctx context.Context, req *pb.DeleteRequest) (*pb.EventResponse, error) {
	id := req.GetId()

	err := s.app.DeleteEvent(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EventResponse{Id: id, Message: "deleted"}, nil
}
