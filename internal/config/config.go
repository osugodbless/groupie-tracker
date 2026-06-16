package config

type Artist struct {
	ID            int                 `json:"id"`
	Image         string              `json:"image"`
	Name          string              `json:"name"`
	Members       []string            `json:"members"`
	CreationYear  int                 `json:"creationDate"`
	FirstAlbum    string              `json:"firstAlbum"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type Relation struct {
	ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type RelationIndex struct {
	Index []Relation `json:"index"`
}
