package reservation

// Repository defines the data access operations for reservations
// All methods are explicitly named with the Reservation entity
type Repository interface {
	ReadReservationList() ([]Reservation, error)
	ReadReservation(id uint) (*Reservation, error)
	CreateReservation(reservation *Reservation) error
	UpdateReservation(reservation *Reservation) error
	DeleteReservation(id uint) error
}
