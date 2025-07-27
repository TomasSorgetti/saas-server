package entities

type User struct {
	ID           int
	Email        string
	Password     string
	Role         string
	FirstName    string
	LastName     string
	Phone        string
	Address      string
	Country      string
	WorkshopName string
	IsActive     bool
	Deleted      bool
	LastLogin    string
	Verified     bool
}