package responses

type Agency struct {
	ID         string       `json:"aid"`
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Entity     string       `json:"entity"`
	UF         string       `json:"uf"`
	FlagURL    string       `json:"url"`
	Collecting []Collecting `json:"collecting"`
}

type Collecting struct {
	Timestamp   *int64   `json:"timestamp"`
	Description []string `json:"description"`
}
