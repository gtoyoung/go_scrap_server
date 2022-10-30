package data

type SportsNews struct {
	Link     string
	Thumnail string
	Title    string
}

type ResponseMsg struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Link  string `json:"link"`
}
