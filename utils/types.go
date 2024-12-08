package utils

import "regexp"

type ListItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TitleData struct {
	IMDbID              string     `json:"imdb_id"`
	Title               string     `json:"title"`
	Type                string     `json:"type"`
	Overview            string     `json:"overview"`
	Poster              string     `json:"poster"`
	Directors           []ListItem `json:"directors"`
	Creators            []ListItem `json:"creators"`
	Genres              []ListItem `json:"genres"`
	ProductionCompanies []ListItem `json:"production_companies"`
	Score               string     `json:"score"`
	ScoredBy            string     `json:"scored_by"`
	Seasons             string     `json:"seasons"`
	Episodes            int        `json:"episodes"`
	Year                string     `json:"year"`
	Rating              string     `json:"rating"`
}

type Cast struct {
	ID        string `json:"imdb_id"`
	Actor     string `json:"actor"`
	Character string `json:"character"`
}

var RegExIMDbID = regexp.MustCompile(`/title/(tt\d+)/`)
var RegExPersonID = regexp.MustCompile(`/name/(nm\d+)/`)
var RegExGenreID = regexp.MustCompile(`/interest/(in\d+)/`)
var RegExCompanyID = regexp.MustCompile(`/company/(co\d+)/`)
