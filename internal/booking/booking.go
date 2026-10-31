package booking

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RoomStore interface {
	SaveRoom(r Room) error
	GetRoom(id string) (Room, error)
}

type ReservationStore interface {
	Save(Reservation) error
	Get(id string) (Reservation, error)
	ListByRoom(roomID string) ([]Reservation, error)
}

var (
	ErrNotFound  = errors.New("booking: not found")
	ErrConflict  = errors.New("booking: slot conflict")
	ErrInvalid   = errors.New("booking: invalid argument")
	ErrCancelled = errors.New("booking: already cancelled")
)

type Room struct {
	ID   string
	Name string
}

type Reservation struct {
	ID        string
	RoomID    string
	Start     int // inclusive
	End       int // exclusive
	Cancelled bool
}

type Service struct {
	rooms RoomStore
	reservations ReservationStore
	nextId int
}

func NewService(rooms RoomStore, reservations ReservationStore) *Service {
	return &Service{
		rooms:        rooms,
		reservations: reservations,
	}
}

func NewMemoryService() *Service {
	return NewService(
		&MemoryRoomStore{rooms: make(map[string]Room)},
		&MemoryReservationStore{
			reservation: make(map[string]Reservation),
			booking:     make(map[string][]string),
		},
	)
}

func (s *Service) AddRoom(id, name string) error {
	var r Room
	r.ID = strings.TrimSpace(id)
	r.Name = strings.TrimSpace(name)

	if r.ID == "" || r.Name == "" {
		return ErrInvalid
	}

	if  _, err := s.rooms.GetRoom(r.ID); err == nil {
		return fmt.Errorf("add room %s: %w", id, ErrInvalid)
	}

	return s.rooms.SaveRoom(r)
}

func (s *Service) Book(roomID string, start, end int) (Reservation, error) {
	if start < 0 || start >= end || end > 24 {
		return Reservation{}, ErrInvalid
	}

	_, err := s.rooms.GetRoom(roomID)
	if err != nil {
		return Reservation{}, fmt.Errorf("book %s: %w", roomID, err)
	}

	reservations, err := s.reservations.ListByRoom(roomID)
	if err != nil {
		return Reservation{}, fmt.Errorf("book %s list: %w", roomID, err)
	}

	for _, res := range reservations {
		if !res.Cancelled && (start < res.End && res.Start < end) {
			return Reservation{}, fmt.Errorf("book %s: %w", roomID, ErrConflict)
		}
	}

	s.nextId++
	id := strconv.Itoa(s.nextId)
	newRes := Reservation{
		ID:     id,
		RoomID: roomID,
		Start:  start,
		End:    end,
	}

	if err := s.reservations.Save(newRes); err != nil {
		return Reservation{}, fmt.Errorf("book %s save: %w", roomID, err)
	}

	return newRes, nil
}

func (s *Service) Cancel(reservationID string) error {
	res, err := s.reservations.Get(reservationID)
	if  err != nil {
		return fmt.Errorf("cancel %s: %w",reservationID, ErrNotFound) 
	}
	if res.Cancelled {
		return fmt.Errorf("cancel %s: %w", reservationID, ErrCancelled)
	}
	res.Cancelled = true
	if err := s.reservations.Save(res); err != nil {
	return fmt.Errorf("cancel %s save: %w", reservationID, err)
}
	return nil
}

func (s *Service) Get(reservationID string) (Reservation, error) {
	res, err := s.reservations.Get(reservationID)
	if  err != nil {
		return Reservation{}, fmt.Errorf("get %s: %w", reservationID, ErrNotFound) 
	}
	
	return res, nil
}

func (s *Service) ListRoom(roomID string) ([]Reservation, error) {
	_, err := s.rooms.GetRoom(roomID)
	if err != nil {
		return []Reservation{}, fmt.Errorf("list %s: %w", roomID, err) 
	}

	res, err := s.reservations.ListByRoom(roomID)
	if err != nil {
		return []Reservation{}, fmt.Errorf("list %s: %w", roomID, err) 
	}
	
	return res, nil
}// includes cancelled