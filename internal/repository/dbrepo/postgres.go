package dbrepo

import (
	"errors"
	"time"

	"github.com/marcelofranco/webapp-go-demo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	result := m.DB.Create(&res)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(res.ID), nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	result := m.DB.Create(&r)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	var count int64
	result := m.DB.Model(&models.RoomRestriction{}).
		Where("room_id = ? AND ? < end_date AND ? > start_date", roomID, end, start).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	if count == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	result := m.DB.Model(&models.Room{}).Select("room.id, room.room_name").
		InnerJoins("room_restrictions").Where("? < room_restrictions.end_date and ? > room_restrictions.start_date", end, start).
		Find(&rooms)
	if result.Error != nil {
		return rooms, result.Error
	}
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room = models.Room{}
	result := m.DB.First(&room, id)
	if result.Error != nil {
		return room, result.Error
	}
	return room, nil
}

// GetUserByID returns a user by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	var user = models.User{}
	result := m.DB.First(&user, id)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	var user = models.User{}
	result := m.DB.First(&user, u.ID)
	if result.Error != nil {
		return result.Error
	}
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.Email = u.Email
	user.AccessLevel = u.AccessLevel

	result = m.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	var user models.User

	result := m.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return int(user.ID), "", result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return int(user.ID), user.Password, nil
}

// AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	result := m.DB.Model(&models.Reservation{}).
		Joins("left join rooms on (reservation.room_id = rooms.id)").
		Order("reservation.start_date asc").
		Find(&reservations)

	if result.Error != nil {
		return reservations, result.Error
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	result := m.DB.Model(&models.Reservation{}).
		Joins("left join rooms on (reservation.room_id = rooms.id)").
		Where("processed = 0").
		Order("reservation.start_date asc").
		Find(&reservations)

	if result.Error != nil {
		return reservations, result.Error
	}

	return reservations, nil
}

// GetReservationByID returns one reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	result := m.DB.Model(&models.Reservation{}).
		Joins("left join rooms on (reservation.room_id = rooms.id)").
		Where("reservation.id = ?", id).
		Find(&reservation)

	if result.Error != nil {
		return reservation, result.Error
	}
	return reservation, nil
}

// UpdateReservation updates a reservation in the database
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	var reservation models.Reservation
	result := m.DB.First(&reservation, u.ID)
	if result.Error != nil {
		return result.Error
	}

	reservation.FirstName = u.FirstName
	reservation.LastName = u.LastName
	reservation.Email = u.Email
	reservation.Phone = u.Phone

	result = m.DB.Save(&reservation)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteReservation deletes one reservation by id
func (m *postgresDBRepo) DeleteReservation(id int) error {
	result := m.DB.Delete(&models.Reservation{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	var reservation models.Reservation
	result := m.DB.First(&reservation, id)
	if result.Error != nil {
		return result.Error
	}

	reservation.Processed = processed

	result = m.DB.Save(&reservation)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room

	result := m.DB.Select("id, room_name, created_at, updated_at").Order("room_name").Find(&rooms)
	if result.Error != nil {
		return rooms, result.Error
	}

	return rooms, nil
}

// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction

	result := m.DB.Select("id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date").
		Where("? < end_date AND ? >= start_date AND room_id = ?", end, start, roomID).
		Find(&restrictions)

	if result.Error != nil {
		return restrictions, result.Error
	}

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	var rr = models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       startDate.AddDate(0, 0, 1),
		RoomID:        id,
		RestrictionID: 2,
	}

	result := m.DB.Create(&rr)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	result := m.DB.Delete(&models.RoomRestriction{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
