package types

type ChatMessage struct {
	Type        string         `json:"type"`
	TextMessage string         `json:"text_message"`
	Buttons     ButtonsWrapper `json:"buttons"`
}

type ButtonsWrapper struct {
	Rows []ButtonRow `json:"rows"`
}

type ButtonRow struct {
	Buttons []Button `json:"buttons"`
}

type Button struct {
	Action   Action `json:"action"`
	IconName string `json:"icon_name"`
	Caption  string `json:"caption"`
}

type Action struct {
	OpenDirectLink string `json:"open_direct_link,omitempty"`
}

type Ad struct {
	Title string `json:"title"`
	Price int    `json:"price"`
	Image string `json:"image"`
	Token string `json:"token"`
}
