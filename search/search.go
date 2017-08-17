package search

import (
	"errors"
	"fmt"
	"log"
	"search-analysis-API/datamodel"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	//set up errors
	ErrNoAPIKey       = errors.New("No API Key")
	ErrNoRadius       = errors.New("No Radius ")
	ErrNotInitialized = errors.New("Not Initialized")

	apikey string // "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	radius uint
	initOk bool
)

type Search struct {
	apikey string
	radius uint
}

func NewSearch(key string, radius uint) *Search {
	return &Search{key, radius}
}

//check apikey is not null
func Initialize(key string, rad uint) error {
	if key == "" {
		return ErrNoAPIKey
	}
	apikey = key
	radius = rad
	initOk = true

	return nil
}

func (s *Search) SetRadius(rad uint) error {
	s.radius = rad
	return nil
}

func (s *Search) Place(keyword string, LAT, LNG float64) ([]datamodel.Coffee, error) {
	location := &maps.LatLng{Lat: LAT, Lng: LNG}
	language := "zh-TW"

	c, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		log.Printf("PlaceSearch MapsAPI error: %s\n", err)
		return nil, err
	}

	request := &maps.NearbySearchRequest{}
	request.Location = location
	request.Radius = s.radius
	request.Keyword = keyword
	request.Language = language
	coffeeList := []datamodel.Coffee{}

	for {
		resp, err := c.NearbySearch(context.Background(), request)
		if err != nil {
			//#WARNING, study HOW it breaks
			fmt.Println("Search over!!")
			break
		}
		for i := 0; i < len(resp.Results); i++ {

			id := resp.Results[i].PlaceID
			Name := resp.Results[i].Name
			Rate := resp.Results[i].Rating

			cof := datamodel.Coffee{}
			cof.Id = id
			cof.Name = Name
			cof.Rate = Rate
			cof.Reviews = []datamodel.Review{}

			req := &maps.PlaceDetailsRequest{}
			req.PlaceID = id
			req.Language = language

			respd, err := c.PlaceDetails(context.Background(), req)
			if err != nil {
				log.Fatalf("fatal error: %s", err)
			}

			for j := 0; j < len(respd.Reviews); j++ {

				review := datamodel.Review{cof.Id, respd.Reviews[j].Text}
				cof.Reviews = append(cof.Reviews, review)
			}

			coffeeList = append(coffeeList, cof)
		}

		request.Location = nil
		request.Radius = 0
		request.Keyword = ""
		request.Language = ""
		if resp.NextPageToken == "" {
			break
		}
		request.PageToken = resp.NextPageToken

	}

	return coffeeList, nil
}

func (s *Search) Placesearchlayer1(keyword string, LAT, LNG float64) error {
	location := &maps.LatLng{Lat: LAT, Lng: LNG}
	language := "zh-TW"

	c, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		log.Printf("PlaceSearch MapsAPI error: %s\n", err)
		return err
	}

	request := &maps.NearbySearchRequest{}
	request.Location = location
	request.Radius = s.radius
	request.Keyword = keyword
	request.Language = language
	Pagecount := 0
	for {
		resp, err := c.NearbySearch(context.Background(), request)
		if err != nil {
			//#WARNING, study HOW it breaks
			fmt.Println(Pagecount)
			fmt.Println("Search over!!")
			break
		}

		request.Location = nil
		request.Radius = 0
		request.Keyword = ""
		request.Language = ""
		Pagecount++
		if resp.NextPageToken == "" {
			break
		}
		request.PageToken = resp.NextPageToken

	}

	return nil
}
func (s *Search) Placesearchlayer2(keyword string, LAT, LNG float64) error {
	location := &maps.LatLng{Lat: LAT, Lng: LNG}
	language := "zh-TW"

	c, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		log.Printf("PlaceSearch MapsAPI error: %s\n", err)
		return err
	}

	request := &maps.NearbySearchRequest{}
	request.Location = location
	request.Radius = s.radius
	request.Keyword = keyword
	request.Language = language
	Pagecount := 0

	for {
		resp, err := c.NearbySearch(context.Background(), request)
		if err != nil {
			//#WARNING, study HOW it breaks
			fmt.Println(Pagecount)
			fmt.Println("Search over!!")
			break
		}
		for i := 0; i < len(resp.Results); i++ {

			id := resp.Results[i].PlaceID

			req := &maps.PlaceDetailsRequest{}
			req.PlaceID = id
			req.Language = language

			_, err := c.PlaceDetails(context.Background(), req)
			if err != nil {
				log.Fatalf("fatal error: %s", err)
			}

			request.Location = nil
			request.Radius = 0
			request.Keyword = ""
			request.Language = ""
			Pagecount++
			if resp.NextPageToken == "" {
				break
			}
			request.PageToken = resp.NextPageToken

		}

	}
	return nil
}
