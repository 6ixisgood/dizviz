package data

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"net/http"
)

type SportsFeed struct {
	Client			*APIClient
}

func NewSportsFeed(baseUrl string, auth BasicAuthCredentials) SportsFeed {
	clientOptions := APIClientOptions {
		BaseURL: "https://scrambled-api.mysportsfeeds.com/v2.1/",
		BasicAuth: &auth,
	}

	client := NewAPIClient(clientOptions)

	return SportsFeed{
		Client: client,	
	}	
}

// -----------------------------------------
// NHL
// -----------------------------------------

type DailyGamesNHLResponse struct {
	LastUpdatedOn time.Time `json:"lastUpdatedOn"`
	Games         []struct {
		Schedule struct {
			ID        int       `json:"id"`
			StartTime time.Time `json:"startTime"`
			EndedTime time.Time `json:"endedTime"`
			AwayTeam  struct {
				ID           int    `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"awayTeam"`
			HomeTeam struct {
				ID           int    `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"homeTeam"`
			Venue struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"venue"`
			VenueAllegiance          string      `json:"venueAllegiance"`
			ScheduleStatus           string      `json:"scheduleStatus"`
			OriginalStartTime        interface{} `json:"originalStartTime"`
			DelayedOrPostponedReason interface{} `json:"delayedOrPostponedReason"`
			PlayedStatus             string      `json:"playedStatus"`
			Attendance               int         `json:"attendance"`
			Officials                []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"officials"`
			Broadcasters []string `json:"broadcasters"`
			Weather      struct {
				Type        string `json:"type"`
				Description string `json:"description"`
				Wind        struct {
					Speed struct {
						MilesPerHour      int `json:"milesPerHour"`
						KilometersPerHour int `json:"kilometersPerHour"`
					} `json:"speed"`
					Direction struct {
						Degrees int    `json:"degrees"`
						Label   string `json:"label"`
					} `json:"direction"`
				} `json:"wind"`
				Temperature struct {
					Fahrenheit int `json:"fahrenheit"`
					Celsius    int `json:"celsius"`
				} `json:"temperature"`
				Precipitation struct {
					Type    interface{} `json:"type"`
					Percent interface{} `json:"percent"`
					Amount  struct {
						Millimeters interface{} `json:"millimeters"`
						Centimeters interface{} `json:"centimeters"`
						Inches      interface{} `json:"inches"`
						Feet        interface{} `json:"feet"`
					} `json:"amount"`
				} `json:"precipitation"`
				HumidityPercent int `json:"humidityPercent"`
			} `json:"weather"`
		} `json:"schedule"`
		Score struct {
			CurrentPeriod                 interface{} `json:"currentPeriod"`
			CurrentPeriodSecondsRemaining interface{} `json:"currentPeriodSecondsRemaining"`
			CurrentIntermission           interface{} `json:"currentIntermission"`
			AwayScoreTotal                int         `json:"awayScoreTotal"`
			AwayShotsTotal                int         `json:"awayShotsTotal"`
			HomeScoreTotal                int         `json:"homeScoreTotal"`
			HomeShotsTotal                int         `json:"homeShotsTotal"`
			Periods                       []struct {
				PeriodNumber int `json:"periodNumber"`
				AwayScore    int `json:"awayScore"`
				AwayShots    int `json:"awayShots"`
				HomeScore    int `json:"homeScore"`
				HomeShots    int `json:"homeShots"`
			} `json:"periods"`
		} `json:"score"`
	} `json:"games"`
	References struct {
		TeamReferences []struct {
			ID           int    `json:"id"`
			City         string `json:"city"`
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
			HomeVenue    struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"homeVenue"`
			TeamColoursHex      []string `json:"teamColoursHex"`
			SocialMediaAccounts []struct {
				MediaType string `json:"mediaType"`
				Value     string `json:"value"`
			} `json:"socialMediaAccounts"`
			OfficialLogoImageSrc string `json:"officialLogoImageSrc"`
		} `json:"teamReferences"`
		VenueReferences []struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			City           string `json:"city"`
			Country        string `json:"country"`
			GeoCoordinates struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"geoCoordinates"`
			CapacitiesByEventType []struct {
				EventType string `json:"eventType"`
				Capacity  int    `json:"capacity"`
			} `json:"capacitiesByEventType"`
			PlayingSurface     interface{}   `json:"playingSurface"`
			BaseballDimensions []interface{} `json:"baseballDimensions"`
			HasRoof            bool          `json:"hasRoof"`
			HasRetractableRoof bool          `json:"hasRetractableRoof"`
		} `json:"venueReferences"`
	} `json:"references"`
}

func (s *SportsFeed) FetchDailyNHLGamesInfo(date string) DailyGamesNHLResponse{
	// determine the "season" to use
	year, err := strconv.Atoi(date[:4])
	checkErr(err)
	month, err := strconv.Atoi(date[4:6])
	checkErr(err)

	var season string
	if month > 6 {
		season = fmt.Sprintf("%d-%d-regular", year, year+1)
	} else {
		season = fmt.Sprintf("%d-%d-regular", year-1, year)
	}

	// fetch the data
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("pull/nhl/%s/date/%s/games.json", season, date),
	}

	var responseData DailyGamesNHLResponse
	if err := s.Client.DoAndUnmarshal(request, &responseData); err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}
// -----------------------------------------
// NFL 
// -----------------------------------------

type DailyGamesNFLResponse struct {
	LastUpdatedOn time.Time `json:"lastUpdatedOn"`
	Games         []struct {
		Schedule struct {
			ID        int         `json:"id"`
			Week      int         `json:"week"`
			StartTime time.Time   `json:"startTime"`
			EndedTime interface{} `json:"endedTime"`
			AwayTeam  struct {
				ID           int    `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"awayTeam"`
			HomeTeam struct {
				ID           int    `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"homeTeam"`
			Venue struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"venue"`
			VenueAllegiance          string        `json:"venueAllegiance"`
			ScheduleStatus           string        `json:"scheduleStatus"`
			OriginalStartTime        interface{}   `json:"originalStartTime"`
			DelayedOrPostponedReason interface{}   `json:"delayedOrPostponedReason"`
			PlayedStatus             string        `json:"playedStatus"`
			Attendance               interface{}   `json:"attendance"`
			Officials                []interface{} `json:"officials"`
			Broadcasters             []string      `json:"broadcasters"`
			Weather                  interface{}   `json:"weather"`
		} `json:"schedule"`
		Score struct {
			CurrentQuarter                 interface{}   `json:"currentQuarter"`
			CurrentQuarterSecondsRemaining interface{}   `json:"currentQuarterSecondsRemaining"`
			CurrentIntermission            interface{}   `json:"currentIntermission"`
			TeamInPossession               interface{}   `json:"teamInPossession"`
			CurrentDown                    interface{}   `json:"currentDown"`
			CurrentYardsRemaining          interface{}   `json:"currentYardsRemaining"`
			LineOfScrimmage                interface{}   `json:"lineOfScrimmage"`
			AwayScoreTotal                 interface{}   `json:"awayScoreTotal"`
			HomeScoreTotal                 interface{}   `json:"homeScoreTotal"`
			Quarters                       []interface{} `json:"quarters"`
		} `json:"score"`
	} `json:"games"`
	References struct {
		TeamReferences []struct {
			ID           int    `json:"id"`
			City         string `json:"city"`
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
			HomeVenue    struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"homeVenue"`
			TeamColoursHex      []string `json:"teamColoursHex"`
			SocialMediaAccounts []struct {
				MediaType string `json:"mediaType"`
				Value     string `json:"value"`
			} `json:"socialMediaAccounts"`
			OfficialLogoImageSrc string `json:"officialLogoImageSrc"`
		} `json:"teamReferences"`
		VenueReferences []struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			City           string `json:"city"`
			Country        string `json:"country"`
			GeoCoordinates struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"geoCoordinates"`
			CapacitiesByEventType []struct {
				EventType string `json:"eventType"`
				Capacity  int    `json:"capacity"`
			} `json:"capacitiesByEventType"`
			PlayingSurface     string        `json:"playingSurface"`
			BaseballDimensions []interface{} `json:"baseballDimensions"`
			HasRoof            bool          `json:"hasRoof"`
			HasRetractableRoof bool          `json:"hasRetractableRoof"`
		} `json:"venueReferences"`
	} `json:"references"`
}


func (s *SportsFeed) FetchDailyNFLGamesInfo(date string) DailyGamesNFLResponse{
	// determine the "season" to use
	year, err := strconv.Atoi(date[:4])
	checkErr(err)
	month, err := strconv.Atoi(date[4:6])
	checkErr(err)

	var season string
	if month > 6 {
		season = fmt.Sprintf("%d-%d-regular", year, year+1)
	} else {
		season = fmt.Sprintf("%d-%d-regular", year-1, year)
	}

	// fetch the data
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("pull/nfl/%s/date/%s/games.json", season, date),
	}

	var responseData DailyGamesNFLResponse
	if err := s.Client.DoAndUnmarshal(request, &responseData); err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}



func checkErr(e error) {
	if (e != nil) {
		log.Printf("Error: %s", e)
	}
}