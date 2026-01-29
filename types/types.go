package types

type DiscordMember struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type MapMarker struct {
	User        DiscordMember `json:"user"`
	Coordinates [2]float64    `json:"coordinates"`
}
