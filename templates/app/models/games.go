package models

type CasinoGame struct {
	ProviderID   int64  `json:"provider_id"`
	Category     string `json:"category"`
	ProviderName string `json:"provider_name"`
	GameID       string `json:"game_id"`
	GameName     string `json:"game_name"`
	Image        string `json:"image"`
	Description  string `json:"Description"`
	Demo         bool   `json:"demo"`
	Type         int8   `json:"type" enums:"1,2"`
}

type GameList struct {
	GameID    string `json:"gameId"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedOn int64  `json:"createdOn"`
}

type GameListResponse struct {
	Brand string     `json:"brand"`
	Games []GameList `json:"games"`
}

type GameUrlRequest struct {
	GameID     string `json:"game_id"`
	AccountID  int64  `json:"account_id"`
	Demo       int64  `json:"demo"`
	ReturnURL  string `json:"return_url"`
	DeviceType string `json:"device_type"`
}

type GameResponse struct {
	PlayURL string `json:"url"`
}

type Player struct {
	PlayerID      string `json:"player_id"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	Currency      string `json:"currency"`
	ExternalToken string `json:"external_token"`
}

type LaunchGame struct {
	CasinoID     string `json:"casino_id"`
	Brand        string `json:"brand"`
	ClientType   string `json:"client_type"`
	GameID       string `json:"game_id"`
	IP           string `json:"ip"`
	Jurisdiction string `json:"jurisdiction"`
	Locale       string `json:"locale"`
	ReturnURL    string `json:"return_url"`
	Demo         bool   `json:"demo"`
	Player       Player `json:"player"`
}
