package booking

import (
	"errors"
	"io"
	"testing"
)

type stubRooms struct {
	rooms map[string]Room
	fail  error
}

func (s *stubRooms) SaveRoom(r Room) error {
	if s.fail != nil {
		return s.fail
	}
	if s.rooms == nil {
		s.rooms = make(map[string]Room)
	}
	s.rooms[r.ID] = r
	return nil
}

func (s *stubRooms) GetRoom(id string) (Room, error) {
	if s.fail != nil {
		return Room{}, s.fail
	}
	r, ok := s.rooms[id]
	if !ok {
		return Room{}, ErrNotFound
	}
	return r, nil
}

type stubReservations struct {
	data     map[string]Reservation
	bookings map[string][]string

	failSave error
	failGet  error
	failList error
}

func (s *stubReservations) Save(r Reservation) error {
	if s.failSave != nil {
		return s.failSave
	}
	if s.data == nil {
		s.data = make(map[string]Reservation)
	}
	if s.bookings == nil {
		s.bookings = make(map[string][]string)
	}
	s.data[r.ID] = r

	s.bookings[r.RoomID] = append(s.bookings[r.RoomID], r.ID)
	return nil
}

func (s *stubReservations) Get(id string) (Reservation, error) {
	if s.failGet != nil {
		return Reservation{}, s.failGet
	}
	r, ok := s.data[id]
	if !ok {
		return Reservation{}, ErrNotFound
	}
	return r, nil
}

func (s *stubReservations) ListByRoom(roomID string) ([]Reservation, error) {
	if s.failList != nil {
		return nil, s.failList
	}
	ids := s.bookings[roomID]
	result := make([]Reservation, 0, len(ids))
	for _, id := range ids {
		result = append(result, s.data[id])
	}
	return result, nil
}

func TestCancelMissing(t *testing.T) {
	s := NewMemoryService()
	err := s.Cancel("nope")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestReservationSaveFails(t *testing.T) {
	rooms := &stubRooms{rooms: map[string]Room{"1": {ID: "1"}}}
	res := &stubReservations{failSave: io.ErrClosedPipe}
	
	s := NewService(rooms, res)
	
	_, err := s.Book("1", 10, 12)
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("expected io.ErrClosedPipe, got %v", err)
	}
}

func TestBookOverlap(t *testing.T) {
	s := NewMemoryService()
	_ = s.AddRoom("orion", "Orion")
	_, _ = s.Book("orion", 9, 11)
	
	_, err := s.Book("orion", 9, 12)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestBookAfterCancel(t *testing.T) {
	s := NewMemoryService()
	_ = s.AddRoom("orion", "Orion")
	a, _ := s.Book("orion", 9, 11)
	_ = s.Cancel(a.ID)
	
	_, err := s.Book("orion", 9, 11)
	if err != nil {
		t.Fatalf("expected success after cancel, got %v", err)
	}
}