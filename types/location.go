package types

// Location is location schema
type Location struct {
	Id        uint `json:"-" korm:"pk"` // primary key
	UserId    uint `json:"user_id" korm:"fk:users.id:cascade:noaction"` // foreign key -> users.id
	Latitude  float64
	Longitude float64
}
