package content_data

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"net/http"
)

type NHLRequestConfig struct {
	APIUsername string
	APIPassword string
}

const (
	format = "json"
	dailyGamesNHLURL	= "https://scrambled-api.mysportsfeeds.com/v2.1/pull/nhl/%s/date/%s/games.json"
)

var (
	NHLClientConfig 				= &NHLRequestConfig{}
)

func makeAPIRequest(url string) []byte {
	log.Printf("Making API Requst to %v", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
		return []byte("")
	}
	basicAuthEncoded := base64.StdEncoding.EncodeToString([]byte(NHLClientConfig.APIUsername + ":" + NHLClientConfig.APIPassword))
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

func FetchDailyNHLGamesInfo(date string) DailyGamesNHLResponse{
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

	url := fmt.Sprintf(dailyGamesNHLURL, season, date)
	fmt.Sprintf(url)
	res := makeAPIRequest(url)
	var respStruct DailyGamesNHLResponse
	json.Unmarshal(res, &respStruct)

	// j := `{"lastUpdatedOn":"2023-02-28T21:10:39.477Z","games":[{"schedule":{"id":78209,"startTime":"2022-10-12T23:00:00.000Z","endedTime":"2022-10-13T01:33:30.000Z","awayTeam":{"id":11,"abbreviation":"BOS"},"homeTeam":{"id":5,"abbreviation":"WSH"},"venue":{"id":13,"name":"Capital One Arena"},"venueAllegiance":"HOME","scheduleStatus":"NORMAL","originalStartTime":null,"delayedOrPostponedReason":null,"playedStatus":"COMPLETED","attendance":18573,"officials":[{"id":29421,"title":"Referee","firstName":"Michael","lastName":"Markovic"},{"id":22266,"title":"Referee","firstName":"Kevin","lastName":"Pollock"},{"id":20744,"title":"Linesman","firstName":"Jonny","lastName":"Murray"},{"id":20754,"title":"Linesman","firstName":"Brad","lastName":"Kovachik"}],"broadcasters":["TNT"],"weather":{"type":"OVERCAST","description":"overcast clouds","wind":{"speed":{"milesPerHour":4,"kilometersPerHour":6},"direction":{"degrees":287,"label":"WNW"}},"temperature":{"fahrenheit":67,"celsius":19},"precipitation":{"type":null,"percent":null,"amount":{"millimeters":null,"centimeters":null,"inches":null,"feet":null}},"humidityPercent":29}},"score":{"currentPeriod":null,"currentPeriodSecondsRemaining":null,"currentIntermission":null,"awayScoreTotal":5,"awayShotsTotal":25,"homeScoreTotal":2,"homeShotsTotal":33,"periods":[{"periodNumber":1,"awayScore":2,"awayShots":13,"homeScore":0,"homeShots":12},{"periodNumber":2,"awayScore":1,"awayShots":7,"homeScore":2,"homeShots":8},{"periodNumber":3,"awayScore":2,"awayShots":5,"homeScore":0,"homeShots":13}]}},{"schedule":{"id":78210,"startTime":"2022-10-12T23:00:00.000Z","endedTime":"2022-10-13T01:51:54.000Z","awayTeam":{"id":19,"abbreviation":"CBJ"},"homeTeam":{"id":3,"abbreviation":"CAR"},"venue":{"id":21,"name":"PNC Arena"},"venueAllegiance":"HOME","scheduleStatus":"NORMAL","originalStartTime":null,"delayedOrPostponedReason":null,"playedStatus":"COMPLETED","attendance":18824,"officials":[{"id":20633,"title":"Referee","firstName":"Kyle","lastName":"Rehman"},{"id":29361,"title":"Referee","firstName":"Kendrick","lastName":"Nicholson"},{"id":29401,"title":"Linesman","firstName":"Julien","lastName":"Fournier"},{"id":35291,"title":"Linesman","firstName":"Ben","lastName":"O'Quinn"}],"broadcasters":["BSSO","BSOH"],"weather":{"type":"OVERCAST","description":"overcast clouds","wind":{"speed":{"milesPerHour":7,"kilometersPerHour":11},"direction":{"degrees":158,"label":"ESE"}},"temperature":{"fahrenheit":63,"celsius":17},"precipitation":{"type":null,"percent":null,"amount":{"millimeters":null,"centimeters":null,"inches":null,"feet":null}},"humidityPercent":66}},"score":{"currentPeriod":null,"currentPeriodSecondsRemaining":null,"currentIntermission":null,"awayScoreTotal":1,"awayShotsTotal":32,"homeScoreTotal":4,"homeShotsTotal":43,"periods":[{"periodNumber":1,"awayScore":0,"awayShots":10,"homeScore":0,"homeShots":10},{"periodNumber":2,"awayScore":1,"awayShots":14,"homeScore":2,"homeShots":13},{"periodNumber":3,"awayScore":0,"awayShots":8,"homeScore":2,"homeShots":20}]}},{"schedule":{"id":78214,"startTime":"2022-10-12T23:00:00.000Z","endedTime":"2022-10-13T01:55:32.000Z","awayTeam":{"id":12,"abbreviation":"TOR"},"homeTeam":{"id":14,"abbreviation":"MTL"},"venue":{"id":4,"name":"Bell Centre"},"venueAllegiance":"HOME","scheduleStatus":"NORMAL","originalStartTime":null,"delayedOrPostponedReason":null,"playedStatus":"COMPLETED","attendance":21105,"officials":[{"id":20520,"title":"Referee","firstName":"Dan","lastName":"O'Rourke"},{"id":29378,"title":"Referee","firstName":"Jake","lastName":"Brenk"},{"id":22271,"title":"Linesman","firstName":"Libor","lastName":"Suchanek"},{"id":22267,"title":"Linesman","firstName":"Michel","lastName":"Cormier"}],"broadcasters":["TVAS","SN"],"weather":{"type":"OVERCAST","description":"overcast clouds","wind":{"speed":{"milesPerHour":10,"kilometersPerHour":16},"direction":{"degrees":162,"label":"S"}},"temperature":{"fahrenheit":61,"celsius":16},"precipitation":{"type":null,"percent":null,"amount":{"millimeters":null,"centimeters":null,"inches":null,"feet":null}},"humidityPercent":60}},"score":{"currentPeriod":null,"currentPeriodSecondsRemaining":null,"currentIntermission":null,"awayScoreTotal":3,"awayShotsTotal":32,"homeScoreTotal":4,"homeShotsTotal":23,"periods":[{"periodNumber":1,"awayScore":1,"awayShots":8,"homeScore":0,"homeShots":8},{"periodNumber":2,"awayScore":1,"awayShots":13,"homeScore":2,"homeShots":9},{"periodNumber":3,"awayScore":1,"awayShots":11,"homeScore":2,"homeShots":6}]}},{"schedule":{"id":78211,"startTime":"2022-10-13T01:30:00.000Z","endedTime":"2022-10-13T04:32:20.000Z","awayTeam":{"id":20,"abbreviation":"CHI"},"homeTeam":{"id":22,"abbreviation":"COL"},"venue":{"id":19,"name":"Ball Arena"},"venueAllegiance":"HOME","scheduleStatus":"NORMAL","originalStartTime":null,"delayedOrPostponedReason":null,"playedStatus":"COMPLETED","attendance":18143,"officials":[{"id":29400,"title":"Referee","firstName":"Chris","lastName":"Schlenker"},{"id":20651,"title":"Referee","firstName":"Brian","lastName":"Pochmara"},{"id":29375,"title":"Linesman","firstName":"Bevan","lastName":"Mills"},{"id":20744,"title":"Linesman","firstName":"Jonny","lastName":"Murray"}],"broadcasters":["TNT"],"weather":{"type":"SUNNY","description":"sky is clear","wind":{"speed":{"milesPerHour":17,"kilometersPerHour":27},"direction":{"degrees":57,"label":"ENE"}},"temperature":{"fahrenheit":64,"celsius":18},"precipitation":{"type":null,"percent":null,"amount":{"millimeters":null,"centimeters":null,"inches":null,"feet":null}},"humidityPercent":15}},"score":{"currentPeriod":null,"currentPeriodSecondsRemaining":null,"currentIntermission":null,"awayScoreTotal":2,"awayShotsTotal":15,"homeScoreTotal":5,"homeShotsTotal":30,"periods":[{"periodNumber":1,"awayScore":1,"awayShots":3,"homeScore":2,"homeShots":10},{"periodNumber":2,"awayScore":0,"awayShots":6,"homeScore":2,"homeShots":8},{"periodNumber":3,"awayScore":1,"awayShots":6,"homeScore":1,"homeShots":12}]}},{"schedule":{"id":78212,"startTime":"2022-10-13T02:00:00.000Z","endedTime":"2022-10-13T04:43:08.000Z","awayTeam":{"id":143,"abbreviation":"SEA"},"homeTeam":{"id":29,"abbreviation":"ANA"},"venue":{"id":8,"name":"Honda Center"},"venueAllegiance":"HOME","scheduleStatus":"NORMAL","originalStartTime":null,"delayedOrPostponedReason":null,"playedStatus":"COMPLETED","attendance":17530,"officials":[{"id":20650,"title":"Referee","firstName":"Francois","lastName":"StLaurent"},{"id":20738,"title":"Referee","firstName":"TJ","lastName":"Luxmore"},{"id":35186,"title":"Linesman","firstName":"Caleb","lastName":"Apperson"},{"id":20637,"title":"Linesman","firstName":"Bryan","lastName":"Pancich"}],"broadcasters":["BSSC","ROOT-NW"],"weather":{"type":"LIGHT_RAIN","description":"light rain","wind":{"speed":{"milesPerHour":10,"kilometersPerHour":16},"direction":{"degrees":254,"label":"W"}},"temperature":{"fahrenheit":74,"celsius":23},"precipitation":{"type":"RAIN","percent":null,"amount":{"millimeters":0,"centimeters":null,"inches":null,"feet":null}},"humidityPercent":61}},"score":{"currentPeriod":null,"currentPeriodSecondsRemaining":null,"currentIntermission":null,"awayScoreTotal":4,"awayShotsTotal":48,"homeScoreTotal":5,"homeShotsTotal":27,"periods":[{"periodNumber":1,"awayScore":1,"awayShots":20,"homeScore":1,"homeShots":10},{"periodNumber":2,"awayScore":2,"awayShots":13,"homeScore":1,"homeShots":8},{"periodNumber":3,"awayScore":1,"awayShots":14,"homeScore":2,"homeShots":8},{"periodNumber":4,"awayScore":0,"awayShots":1,"homeScore":1,"homeShots":1}]}},{"schedule":{"id":78213,"startTime":"2022-10-13T02:00:00.000Z","endedTime":"2022-10-13T05:02:07.000Z","awayTeam":{"id":21,"abbreviation":"VAN"},"homeTeam":{"id":24,"abbreviation":"EDM"},"venue":{"id":22,"name":"Rogers Place"},"venueAllegiance":"HOME","scheduleStatus":"NORMAL","originalStartTime":null,"delayedOrPostponedReason":null,"playedStatus":"COMPLETED","attendance":18347,"officials":[{"id":29418,"title":"Referee","firstName":"Garrett","lastName":"Rank"},{"id":20755,"title":"Referee","firstName":"Kelly","lastName":"Sutherland"},{"id":29382,"title":"Linesman","firstName":"Brandon","lastName":"Gawryletz"},{"id":22275,"title":"Linesman","firstName":"Mark","lastName":"Shewchyk"}],"broadcasters":["SN"],"weather":{"type":"SUNNY","description":"sky is clear","wind":{"speed":{"milesPerHour":15,"kilometersPerHour":24},"direction":{"degrees":315,"label":"NW"}},"temperature":{"fahrenheit":49,"celsius":9},"precipitation":{"type":null,"percent":null,"amount":{"millimeters":null,"centimeters":null,"inches":null,"feet":null}},"humidityPercent":39}},"score":{"currentPeriod":null,"currentPeriodSecondsRemaining":null,"currentIntermission":null,"awayScoreTotal":3,"awayShotsTotal":36,"homeScoreTotal":5,"homeShotsTotal":25,"periods":[{"periodNumber":1,"awayScore":2,"awayShots":14,"homeScore":0,"homeShots":8},{"periodNumber":2,"awayScore":1,"awayShots":14,"homeScore":3,"homeShots":7},{"periodNumber":3,"awayScore":0,"awayShots":8,"homeScore":2,"homeShots":10}]}}],"references":{"teamReferences":[{"id":3,"city":"Carolina","name":"Hurricanes","abbreviation":"CAR","homeVenue":{"id":21,"name":"PNC Arena"},"teamColoursHex":["#cc0000","#000000","#a2aaad","#76232f"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"CanesNHL"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/12.svg"},{"id":5,"city":"Washington","name":"Capitals","abbreviation":"WSH","homeVenue":{"id":13,"name":"Capital One Arena"},"teamColoursHex":["#041e42","#c8102e","#ffffff"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"Capitals"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/15.svg"},{"id":11,"city":"Boston","name":"Bruins","abbreviation":"BOS","homeVenue":{"id":27,"name":"TD Garden"},"teamColoursHex":["#ffb81c","#000000"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"NHLBruins"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/6.svg"},{"id":12,"city":"Toronto","name":"Maple Leafs","abbreviation":"TOR","homeVenue":{"id":1,"name":"Scotiabank Arena"},"teamColoursHex":["#00205b","#ffffff"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"MapleLeafs"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/10.svg"},{"id":14,"city":"Montreal","name":"Canadiens","abbreviation":"MTL","homeVenue":{"id":4,"name":"Bell Centre"},"teamColoursHex":["#af1e2d","#192168"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"CanadiensMTL"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/8.svg"},{"id":19,"city":"Columbus","name":"Blue Jackets","abbreviation":"CBJ","homeVenue":{"id":17,"name":"Nationwide Arena"},"teamColoursHex":["#002654","#ce1126","#a4a9ad"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"BlueJacketsNHL"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/29.svg"},{"id":20,"city":"Chicago","name":"Blackhawks","abbreviation":"CHI","homeVenue":{"id":28,"name":"United Center"},"teamColoursHex":["#cf0a2c","#ff671b","#00833e","#ffd100","#d18a00","#001970","#000000","#ffffff"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"NHLBlackhawks"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/16.svg"},{"id":21,"city":"Vancouver","name":"Canucks","abbreviation":"VAN","homeVenue":{"id":7,"name":"Rogers Arena"},"teamColoursHex":["#00205b","#00843d","#041c2c","#99999a","#ffffff"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"Canucks"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/23.svg"},{"id":22,"city":"Colorado","name":"Avalanche","abbreviation":"COL","homeVenue":{"id":19,"name":"Ball Arena"},"teamColoursHex":["#6f263d","#236192","#a2aaad","#000000"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"Avalanche"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/21.svg"},{"id":24,"city":"Edmonton","name":"Oilers","abbreviation":"EDM","homeVenue":{"id":22,"name":"Rogers Place"},"teamColoursHex":["#041e42","#ff4c00"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"EdmontonOilers"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/22.svg"},{"id":29,"city":"Anaheim","name":"Ducks","abbreviation":"ANA","homeVenue":{"id":8,"name":"Honda Center"},"teamColoursHex":["#f47a38","#b9975b","#c1c6c8","#000000"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"AnaheimDucks"}],"officialLogoImageSrc":"https://www-league.nhlstatic.com/images/logos/teams-current-primary-light/24.svg"},{"id":143,"city":"Seattle","name":"Kraken","abbreviation":"SEA","homeVenue":{"id":194,"name":"Climate Pledge Arena"},"teamColoursHex":["#355464","#68a2b9","#99d9d9","#e9072b"],"socialMediaAccounts":[{"mediaType":"TWITTER","value":"SeattleKraken"}],"officialLogoImageSrc":"https://cms.nhl.bamgrid.com/images/assets/binary/317578370/binary-file/file.svg"}],"venueReferences":[{"id":4,"name":"Bell Centre","city":"Montreal, QC","country":"Canada","geoCoordinates":{"latitude":45.496111,"longitude":-73.569444},"capacitiesByEventType":[{"eventType":"ICE_HOCKEY","capacity":21302},{"eventType":"BASKETBALL","capacity":22114},{"eventType":"AMPITHEATRE","capacity":14000},{"eventType":"HEMICYCLE","capacity":3500},{"eventType":"MMA","capacity":23152},{"eventType":"CONCERT","capacity":15000},{"eventType":"THEATRE","capacity":9000}],"playingSurface":null,"baseballDimensions":[],"hasRoof":true,"hasRetractableRoof":false},{"id":8,"name":"Honda Center","city":"Anaheim, CA","country":"USA","geoCoordinates":{"latitude":33.807778,"longitude":-117.876667},"capacitiesByEventType":[{"eventType":"ICE_HOCKEY","capacity":17174},{"eventType":"BASKETBALL","capacity":18336},{"eventType":"CONCERT","capacity":18900},{"eventType":"THEATRE","capacity":8400}],"playingSurface":null,"baseballDimensions":[],"hasRoof":true,"hasRetractableRoof":false},{"id":13,"name":"Capital One Arena","city":"Washington, DC","country":"USA","geoCoordinates":{"latitude":38.89801,"longitude":-77.0231909},"capacitiesByEventType":[{"eventType":"ICE_HOCKEY","capacity":18506},{"eventType":"BASKETBALL","capacity":20356}],"playingSurface":"Multiple","baseballDimensions":[],"hasRoof":true,"hasRetractableRoof":false},{"id":19,"name":"Ball Arena","city":"Denver, CO","country":"USA","geoCoordinates":{"latitude":39.748611,"longitude":-105.0075},"capacitiesByEventType":[{"eventType":"ICE_HOCKEY","capacity":17809},{"eventType":"BASKETBALL","capacity":19520},{"eventType":"ARENA_FOOTBALL","capacity":17417},{"eventType":"LACROSSE","capacity":17809},{"eventType":"CONCERT","capacity":20000}],"playingSurface":"Multiple","baseballDimensions":[],"hasRoof":true,"hasRetractableRoof":false},{"id":21,"name":"PNC Arena","city":"Raleigh, NC","country":"USA","geoCoordinates":{"latitude":35.803333,"longitude":-78.721944},"capacitiesByEventType":[{"eventType":"ICE_HOCKEY","capacity":18680},{"eventType":"BASKETBALL","capacity":19722},{"eventType":"CONCERT","capacity":19500}],"playingSurface":"Multiple","baseballDimensions":[],"hasRoof":true,"hasRetractableRoof":false},{"id":22,"name":"Rogers Place","city":"Edmonton, AB","country":"Canada","geoCoordinates":{"latitude":53.546944,"longitude":-113.497778},"capacitiesByEventType":[{"eventType":"ICE_HOCKEY","capacity":18347},{"eventType":"BASKETBALL","capacity":19500},{"eventType":"CONCERT","capacity":20734}],"playingSurface":null,"baseballDimensions":[],"hasRoof":true,"hasRetractableRoof":false}]}}`
	// json.Unmarshal([]byte(j), &respStruct)

	return respStruct
}

func checkErr(e error) {
	if (e != nil) {
		log.Printf("Error: %s", e)
	}
}