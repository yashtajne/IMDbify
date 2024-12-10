package utils

import (
	"regexp"
	"time"
)

type ListItem struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type TitleData struct {
	IMDbID              string     `bson:"imdb_id" json:"imdb_id"`
	Title               string     `bson:"title" json:"title"`
	Type                string     `bson:"type" json:"type"`
	Overview            string     `bson:"overview" json:"overview"`
	Poster              string     `bson:"poster" json:"poster"`
	Directors           []ListItem `bson:"directors" json:"directors"`
	Creators            []ListItem `bson:"creators" json:"creators"`
	Genres              []ListItem `bson:"genres" json:"genres"`
	ProductionCompanies []ListItem `bson:"production_companies" json:"production_companies"`
	Score               string     `bson:"score" json:"score"`
	ScoredBy            string     `bson:"scored_by" json:"scored_by"`
	Seasons             int        `bson:"seasons" json:"seasons"`
	Episodes            int        `bson:"episodes" json:"episodes"`
	Year                string     `bson:"year" json:"year"`
	Rating              string     `bson:"rating" json:"rating"`
	ExpireAt            time.Time  `bson:"expireAt" json:"expireAt"`
}

type Cast struct {
	ID        string `bson:"id" json:"id"`
	Actor     string `bson:"actor" json:"actor"`
	Character string `bson:"character" json:"character"`
}

type Season struct {
	IMDbID   string
	Season   int
	Episodes []Episode
}

type Episode struct {
	Name     string
	Overview string
	Image    string
	Aired    string
	Score    float64
}

var RegExIMDbID = regexp.MustCompile(`/title/(tt\d+)/`)
var RegExPersonID = regexp.MustCompile(`/name/(nm\d+)/`)
var RegExGenreID = regexp.MustCompile(`/interest/(in\d+)/`)
var RegExCompanyID = regexp.MustCompile(`/company/(co\d+)/`)
var RegeExImageHash = regexp.MustCompile(`M/([^\.]+)`)
