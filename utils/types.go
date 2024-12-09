package utils

import (
	"regexp"
	"time"
)

type ListItem struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type TitleData struct {
	IMDbID              string     `bson:"imdb_id"`
	Title               string     `bson:"title"`
	Type                string     `bson:"type"`
	Overview            string     `bson:"overview"`
	Poster              string     `bson:"poster"`
	Directors           []ListItem `bson:"directors"`
	Creators            []ListItem `bson:"creators"`
	Genres              []ListItem `bson:"genres"`
	ProductionCompanies []ListItem `bson:"production_companies"`
	Score               string     `bson:"score"`
	ScoredBy            string     `bson:"scored_by"`
	Seasons             string     `bson:"seasons"`
	Episodes            int        `bson:"episodes"`
	Year                string     `bson:"year"`
	Rating              string     `bson:"rating"`
	ExpireAt            time.Time  `bson:"expireAt"`
}

type Cast struct {
	ID        string `bson:"imdb_id"`
	Actor     string `bson:"actor"`
	Character string `bson:"character"`
}

var RegExIMDbID = regexp.MustCompile(`/title/(tt\d+)/`)
var RegExPersonID = regexp.MustCompile(`/name/(nm\d+)/`)
var RegExGenreID = regexp.MustCompile(`/interest/(in\d+)/`)
var RegExCompanyID = regexp.MustCompile(`/company/(co\d+)/`)
