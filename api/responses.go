package api

type RtmStartResponseData struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Url   string `json:"url"`
	Self  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
}

type UsersInfoResponseData struct {
	Ok    bool      `json:"ok"`
	Error string    `json:"error"`
	User  SlackUser `json:"user"`
}

type SlackUserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RealName  string `json:"real_name"`
	Email     string `json:"email"`
	Image24   string `json:"image_24"`
}

type SlackUser struct {
	Id      string           `json:"id"`
	Name    string           `json:"name"`
	IsAdmin bool             `json:"is_admin"`
	IsOwner bool             `json:"is_owner"`
	Profile SlackUserProfile `json:"profile"`
}
