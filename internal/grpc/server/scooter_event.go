package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "scootin/internal/grpc"
	"scootin/internal/models"
	"scootin/internal/repository"
	"time"
)

type ScooterEventServer struct {
	pb.UnimplementedScooterEventServiceServer

	scooterRepository      repository.ScooterRepositoryInterface
	scooterEventRepository repository.ScooterEventRepositoryInterface
}

func NewScooterEventServer(
	scooterRepository repository.ScooterRepositoryInterface,
	scooterEventRepository repository.ScooterEventRepositoryInterface,
) *ScooterEventServer {
	return &ScooterEventServer{
		scooterRepository:      scooterRepository,
		scooterEventRepository: scooterEventRepository,
	}
}

func (s *ScooterEventServer) Create(ctx context.Context, event *pb.CreateScooterEvent) (*pb.ScooterEvent, error) {
	scooter, err := s.scooterRepository.FindOneByIdentifier(ctx, event.ScooterIdentifier)
	if err != nil {
		panic(err)
	}
	if scooter == nil {
		return nil, errors.New("scooter not found")
	}

	switch event.Event {
	case pb.Event_UPDATE, pb.Event_TRIP_START, pb.Event_TRIP_END:
		break
	default:
		return nil, errors.New("invalid event passed")
	}

	scooterEvent := models.ScooterEvent{
		Identifier:        uuid.NewString(),
		ScooterIdentifier: event.ScooterIdentifier,
		Event:             event.Event.String(),
		Timestamp:         time.Now(),
		Latitude:          event.Latitude,
		Longitude:         event.Longitude,
	}

	scooter.LastConfirmedLatitude = event.Latitude
	scooter.LastConfirmedLongitude = event.Longitude
	err = s.scooterRepository.Update(ctx, scooter)
	if err != nil {
		return nil, errors.New("scooter coordinates update failed")
	}

	err = s.scooterEventRepository.Insert(ctx, &scooterEvent)
	if err != nil {
		return nil, errors.New("insert scooter event failed")
	}

	response := pb.ScooterEvent{
		Identifier:        scooterEvent.Identifier,
		ScooterIdentifier: scooterEvent.ScooterIdentifier,
		Event:             event.Event,
		Latitude:          scooterEvent.Latitude,
		Longitude:         scooterEvent.Longitude,
	}

	return &response, nil
}

func (s *ScooterEventServer) Get(ctx context.Context, event *pb.GetScooterEvent) (*pb.ScooterEvent, error) {
	return nil, nil
}
