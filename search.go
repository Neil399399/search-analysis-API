package main

import (
	"fmt"
	"log"
	"search-analysis-API/datamodel"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func PlaceSearch() ([]datamodel.Coffee, error) {
	var radius uint
	APIKey := "AIzaSyAFictx33AgxsMkYF-fHCkeakTlBiIZIV4"
	location := &maps.LatLng{Lat: 25.054989, Lng: 121.533359}
	radius = 500
	keyword := "coffee"
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
