package booking

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
	rooms map[string]Room
	reservations map[string]Reservation
	booking map[string][]string
	nextId int
}

func NewService() *Service {
	return &Service{
		rooms: make(map[string]Room),
		reservations: make(map[string]Reservation),
		booking: make(map[string][]string),
	}
}

func (s *Service) AddRoom(id, name string) error {
	var r Room
	r.ID = strings.TrimSpace(id)
	r.Name = strings.TrimSpace(name)

	if r.ID == "" || r.Name == "" {
		return ErrInvalid
	}

	if  _, exist := s.rooms[r.ID]; exist {
		return fmt.Errorf("add room %s: %w", id, ErrInvalid)
	}

	s.rooms[r.ID] = r
	return nil
}

func (s *Service) Book(roomID string, start, end int) (Reservation, error) {
	if _, exist := s.rooms[roomID]; !exist {
		return Reservation{}, ErrNotFound
	}
	if start < 0 || start >= end || end > 24 {
		return Reservation{}, ErrInvalid
	}

	for _, id := range s.booking[roomID] {
		res := s.reservations[id]
		if !res.Cancelled && (start < res.End && res.Start < end) {
			return Reservation{}, fmt.Errorf("book %s: %w", roomID, ErrConflict)
		}
	}
	s.nextId++
	id := strconv.Itoa(s.nextId) 
	res := Reservation{
		ID: id,
		RoomID: roomID,
		Start:  start,
		End:    end,
	}
	s.reservations[id] = res
	s.booking[roomID] = append(s.booking[roomID], id) 

	return res, nil
}

func (s *Service) Cancel(reservationID string) error {
	res, ok := s.reservations[reservationID]
	if  !ok {
		return fmt.Errorf("cancel %s: %w",reservationID, ErrNotFound) 
	}
	if res.Cancelled {
		return fmt.Errorf("cancel %s: %w", reservationID, ErrCancelled)
	}
	res.Cancelled = true
	s.reservations[reservationID] = res
	return nil
}

func (s *Service) Get(reservationID string) (Reservation, error) {
	res, ok := s.reservations[reservationID]
	if  !ok {
		return Reservation{}, fmt.Errorf("get %s: %w", reservationID, ErrNotFound) 
	}
	
	return res, nil
}

func (s *Service) ListRoom(roomID string) ([]Reservation, error) {
	_, ok := s.rooms[roomID]
	if !ok {
		return []Reservation{}, fmt.Errorf("list %s: %w", roomID, ErrNotFound) 
	}

	res := s.booking[roomID]
	result := make([]Reservation, 0, len(res))

	for _, v := range res {
		result = append(result, s.reservations[v])
	}
	return result, nil
}// includes cancelled