package config

type Data struct {
	Artists   []Artist
	Locations []Location
	Dates     []Date
	Relations []Relation
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationYear int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationIndex struct {
	Index []Location `json:"index"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DateIndex struct {
	Index []Date `json:"index"`
}

type Relation struct {
	ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type RelationIndex struct {
	Index []Relation `json:"index"`
}
