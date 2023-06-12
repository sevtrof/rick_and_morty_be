package model

type Character struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Species  string `json:"species"`
	Type     string `json:"type"`
	Gender   string `json:"gender"`
	Image    string `json:"image"`
	Url      string `json:"url"`
	Created  string `json:"created"`
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Origin struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"origin"`
	Episode []string `json:"episode"`
}

type APIResponse struct {
	Info    Info        `json:"info"`
	Results []Character `json:"results"`
}
