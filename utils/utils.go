package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Scrape(imdbID string) (TitleData, error) {
	var data TitleData

	err := Titles.FindOne(Ctx, bson.M{"imdb_id": imdbID}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with the given IMDb ID")
		} else {
			log.Fatalf("Error finding document: %v", err)
		}
	} else {
		return data, nil
	}

	var c = colly.NewCollector()

	// IMDb ID
	data.IMDbID = imdbID

	// Title Name
	c.OnHTML("span.hero__primary-text", func(e *colly.HTMLElement) {
		data.Title = e.Text
	})

	// Type
	c.OnHTML("h1[data-testid='hero__pageTitle']", func(e *colly.HTMLElement) {
		data.Type = e.DOM.Next().Find("li").First().Text()
	})

	// Short Overview
	c.OnHTML("div[data-testid='interests']", func(e *colly.HTMLElement) {
		data.Overview = e.DOM.NextFiltered("p").Find("span").Eq(1).Text()
	})

	// Poster Image
	c.OnHTML("img.ipc-image", func(e *colly.HTMLElement) {
		if e.Index == 0 {
			match := RegeExImageHash.FindStringSubmatch(e.Attr("src"))
			if len(match) > 1 {
				data.Poster = match[1]
			}
		}
	})

	// Number of Seasons
	c.OnHTML("select#browse-episodes-season", func(e *colly.HTMLElement) {
		seasons := e.Attr("aria-label")
		data.Seasons, err = strconv.Atoi(strings.TrimSuffix(seasons, " seasons"))
		if err != nil {
			log.Println("Error converting text to integer:", err)
			return
		}
	})

	// Number of Episodes
	c.OnHTML("span.ipc-title__subtext", func(e *colly.HTMLElement) {
		if e.Index == 0 {
			episodes, err := strconv.Atoi(strings.TrimSpace(e.Text))
			if err != nil {
				log.Println("Error converting text to integer:", err)
				return
			}
			data.Episodes = episodes
		}
	})

	// Score & Scored By
	c.OnHTML("div[data-testid='hero-rating-bar__aggregate-rating__score']", func(e *colly.HTMLElement) {
		if e.Index == 0 {
			data.Score = e.DOM.Find("span").Text()
			data.ScoredBy = e.DOM.Next().Next().Text()
		}
	})

	// Released Year
	c.OnHTML(fmt.Sprintf("a[href='/title/%s/releaseinfo?ref_=tt_ov_rdat']", imdbID), func(e *colly.HTMLElement) {
		data.Year = e.Text
	})

	// Rating
	c.OnHTML(fmt.Sprintf("a[href='/title/%s/parentalguide/certificates?ref_=tt_ov_pg']", imdbID), func(e *colly.HTMLElement) {
		data.Rating = strings.TrimSpace(e.Text)
	})

	// Directors
	foundDirectors := false
	c.OnHTML("span", func(e *colly.HTMLElement) {
		if !foundDirectors && (strings.Contains(e.Text, "Director") || strings.Contains(e.Text, "Directors")) {
			e.DOM.Next().Find("a").Each(func(_ int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists {
					match := RegExPersonID.FindStringSubmatch(href)
					if len(match) > 1 {
						data.Directors = append(data.Directors, ListItem{ID: match[1], Name: s.Text()})
					}
				}
			})
			foundDirectors = true
		}
	})

	// Creators
	c.OnHTML(fmt.Sprintf("a[href='/title/%s/fullcredits/?ref_=tt_cst_scc_sm#writer']", imdbID), func(e *colly.HTMLElement) {
		e.DOM.Next().Find("a").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				match := RegExPersonID.FindStringSubmatch(href)
				if len(match) > 1 {
					data.Creators = append(data.Creators, ListItem{ID: match[1], Name: s.Text()})
				}
			}
		})
	})

	// Genres
	c.OnHTML("div[data-testid='interests']", func(e *colly.HTMLElement) {
		e.DOM.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				match := RegExGenreID.FindStringSubmatch(href)
				if len(match) > 1 {
					data.Genres = append(data.Genres, ListItem{ID: match[1], Name: s.Text()})
				}
			}
		})
	})

	// Production Companies
	c.OnHTML("a:contains('Production companies')", func(e *colly.HTMLElement) {
		e.DOM.Next().Find("a").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				match := RegExCompanyID.FindStringSubmatch(href)
				if len(match) > 1 {
					data.ProductionCompanies = append(data.ProductionCompanies, ListItem{ID: match[1], Name: s.Text()})
				}
			}
		})
	})

	err = c.Visit("https://www.imdb.com/title/" + imdbID)
	if err != nil {
		return data, err
	}

	c.Wait()

	// expire after 7 days
	data.ExpireAt = time.Now().AddDate(0, 0, 7)

	if len(data.Directors) != 0 {
		data.Type = "Movie"
		data.Seasons = 0
		data.Episodes = 0
	}

	_, err = Titles.InsertOne(Ctx, data)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
		return data, err
	}

	return data, nil
}

func Search(query string) ([]TitleData, error) {
	var c = colly.NewCollector(
		colly.Async(true),
	)

	var wg sync.WaitGroup
	var results []TitleData

	c.OnHTML("a.ipc-title-link-wrapper", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		match := RegExIMDbID.FindStringSubmatch(href)

		if len(match) > 1 {
			wg.Add(1)
			go func(imdbID string) {
				defer wg.Done()
				data, err := Scrape(imdbID)
				if err != nil {
					log.Printf("Error scraping IMDb ID %s: %v", imdbID, err)
					return
				}
				results = append(results, data)
			}(match[1])
		}
	})

	err := c.Visit("https://www.imdb.com/search/title/?title=" + query + "&title_type=feature,tv_series,tv_miniseries")
	if err != nil {
		log.Fatal("Failed to visit the website:", err)
	}

	c.Wait()
	wg.Wait()

	return results, nil
}

func ScrapeCast(imdbID string) ([]Cast, error) {
	var c = colly.NewCollector(
		colly.Async(true),
	)

	var cast []Cast

	c.OnHTML("tr.odd, tr.even", func(e *colly.HTMLElement) {
		a := e.DOM.Find("td.primary_photo").Find("a")

		href, existsHref := a.Attr("href")
		name, _ := a.Find("img").Attr("title")

		if !existsHref {
			return
		}

		cast = append(cast, Cast{
			ID:        RegExPersonID.FindStringSubmatch(href)[1],
			Actor:     name,
			Character: e.DOM.Find("td.character").Find("a").First().Text(),
		})
	})

	err := c.Visit("https://www.imdb.com/title/" + imdbID + "/fullcredits")
	if err != nil {
		log.Fatal("Failed to visit the website:", err)
	}

	c.Wait()

	return cast, nil
}

func ScrapeEpisodes(imdbID string, season int) (Season, error) {
	var c = colly.NewCollector(
		colly.Async(true),
	)

	var episodes []Episode

	c.OnHTML("article.episode-item-wrapper", func(e *colly.HTMLElement) {
		match := RegeExImageHash.FindStringSubmatch(e.ChildAttr("img", "src"))
		if len(match) > 1 {
			score, err := strconv.ParseFloat(e.ChildText("span.ipc-rating-star--rating"), 64)
			if err != nil {
				log.Fatal(err)
				return
			}
			episode := Episode{
				Name:     e.ChildText("div.ipc-title__text"),
				Overview: e.ChildText("div.ipc-html-content-inner-div"),
				Image:    match[1],
				Aired:    e.DOM.Find("h4[data-testid='slate-list-card-title']").Next().Text(),
				Score:    score,
			}

			episodes = append(episodes, episode)
		}
	})

	err := c.Visit("https://www.imdb.com/title/" + imdbID + "/episodes/?season=" + strconv.Itoa(season))
	if err != nil {
		log.Fatal("Failed to visit the website:", err)
	}

	c.Wait()

	return Season{
		IMDbID:   imdbID,
		Season:   season,
		Episodes: episodes,
	}, nil
}
