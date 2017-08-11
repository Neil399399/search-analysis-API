package main

import (
	"fmt"
	"log"
	"search-analysis-API/datamodel"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	APIKey   string // "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	Lat, Lng float64
	keyword  string //"coffee"

)

func PlaceSearch(KEYWORD string, LAT, LNG float64) ([]datamodel.Coffee, error) {
	var radius uint
	location := &maps.LatLng{Lat: LAT, Lng: LNG}
	radius = 500
	keyword = KEYWORD
	language := "zh-TW"

	c, err := maps.NewClient(maps.WithAPIKey(APIKey))
	if err != nil {
		log.Printf("PlaceSearch MapsAPI error: %s\n", err)
		return nil, err
	}

	request := &maps.NearbySearchRequest{}
	request.Location = location
	request.Radius = radius
	request.Keyword = keyword
	request.Language = language
	coffeeList := []datamodel.Coffee{}
	for {

		resp, err := c.NearbySearch(context.Background(), request)
		if err != nil {
			//#WARNING, study HOW it breaks
			fmt.Println("Search over!!", err)
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
