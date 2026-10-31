package booking

// type Service struct {
// 	rooms map[string]Room
// 	reservations map[string]Reservation
// 	booking map[string][]string
// 	nextId int
// }



type MemoryRoomStore struct {
	rooms map[string]Room
}

func (s *MemoryRoomStore) SaveRoom(r Room) error {
	s.rooms[r.ID] = r
	return nil
}

func (s *MemoryRoomStore) GetRoom(id string) (Room, error) {
	r, ok := s.rooms[id]
	if !ok {
		return Room{}, ErrNotFound
	}
	return r, nil
}

type MemoryReservationStore struct {
	reservation map[string]Reservation
	booking map[string][]string
}

func (s *MemoryReservationStore) Save(r Reservation) error {
	_, exists := s.reservation[r.ID]
	
	s.reservation[r.ID] = r
	
	if !exists {
		s.booking[r.RoomID] = append(s.booking[r.RoomID], r.ID)
	}
	return nil
}
	
func (s *MemoryReservationStore) Get(id string) (Reservation, error) {
	r, ok := s.reservation[id]
	if !ok {
		return Reservation{}, ErrNotFound
	}
	return r, nil
}

func (s *MemoryReservationStore) ListByRoom(roomID string) ([]Reservation, error) {
	ids := s.booking[roomID]
	
	result := make([]Reservation, 0, len(ids))
	for _, id := range ids {
		result = append(result, s.reservation[id])
	}
	return result, nil
}