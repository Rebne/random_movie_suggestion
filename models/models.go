package models

type ID struct {
	MovieID string `json:"movieID"`
	Index   int    `json:"index"`
	Title   string `json:"title"`
}

type IDdata struct {
	Length int  `json:"length"`
	IDs    []ID `json:"ids"`
}

type MovieData struct {
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Plot    string `json:"Plot"`
	Runtime string `json:"Runtime"`
	Poster  string `json:"Poster"`
	Genre   string `json:"Genre"`
}

func (data *IDdata) ReIndexMovieIDs() {
	for i := 0; i < data.Length; i++ {
		data.IDs[i].Index = i
	}
}
