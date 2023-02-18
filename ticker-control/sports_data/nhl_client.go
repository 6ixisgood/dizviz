package sports_data

import (
	"fmt"
	"log"
	"time"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"net/http"
)

type Config struct {
	APIUsername string
	APIPassword string
}

const (
	format = "json"
	dailyGamesNHLURL	= "https://api.mysportsfeeds.com/v2.1/pull/nhl/%s/date/%s/games.json"
)

var (
	ClientConfig 				= &Config{}
)

func makeAPIRequest(url string) []byte {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
		return []byte("")
	}
	basicAuthEncoded := base64.StdEncoding.EncodeToString([]byte(ClientConfig.APIUsername + ":" + ClientConfig.APIPassword))
	req.Header.Add("Authorization", "Basic " + basicAuthEncoded)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return []byte("")
	}

	defer res.Body.Close()


	bodyBytes, _ := ioutil.ReadAll(res.Body)
	bodyString := string(bodyBytes)
	if 200 > res.StatusCode && res.StatusCode > 299 {
		log.Fatalln("Request failed:", res.StatusCode, bodyString)
		return []byte("")
	}
	
	return bodyBytes

}


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

func FetchDailyNHLGamesInfo(season string, date string) DailyGamesNHLResponse{
	url := fmt.Sprintf(dailyGamesNHLURL, season, date)
	res := makeAPIRequest(url)
	var respStruct DailyGamesNHLResponse
	json.Unmarshal(res, &respStruct)

	return respStruct
}