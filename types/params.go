package types

type ProfileParams struct {
	Age    string `query:"age"`
	Gender string `query:"gender"`
}

type SwipeParams struct {
	ProfileId  string `form:"profileId"`
	Preference string `form:"preference"`
}

type LoginParams struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}