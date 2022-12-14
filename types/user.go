package types

type User struct {
	Id             uint    `json:"id" korm:"pk"`                           // primary key, auto increment id
	Email          string  `json:"email,omitempty" korm:"size:50;iunique"` // email is insensitive unique
	Password       string  `json:"password,omitempty" korm:"size:150"`
	Name           string  `json:"name" korm:"size:50;iunique"`                                 // name is insensitive unique
	Gender         string  `json:"gender" korm:"size:10;check: gender in ('male','female','')"` // add db check
	Age            uint    `json:"age" korm:"check: age >= 18 and age < 100"`                   // add db check
	DistanceFromMe float64 `json:"distanceFromMe,omitempty" korm:"-"`
}

type Swipe struct {
	Id         uint   `json:"id" korm:"pk"`
	UserId     uint   `json:"user_id" korm:"fk:users.id:cascade:noaction"` // fk:table.column:onDelete:onUpdate
	ProfileId  uint   `json:"profile_id" korm:"fk:users.id:cascade:noaction"`
	Preference string `json:"preference" korm:"size:10;check: preference in ('yes','no')"` // add db check
}
