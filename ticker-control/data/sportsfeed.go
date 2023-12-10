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

type SportsFeedConfig struct {
	BaseUrl			string
	Username		string
	Password		string
}

// Sports Feed Singleton
var sportsFeedClient *SportsFeed

func SportsFeedClient() *SportsFeed {
	return sportsFeedClient
}

func InitSportsFeedClient(config SportsFeedConfig) {
	clientOptions := APIClientOptions{
		BaseURL: config.BaseUrl,
		BasicAuth: &BasicAuthCredentials{
			Username: config.Username,
			Password: config.Password,
		},
	}

	client := NewAPIClient(clientOptions)

	// set the singleton
	sportsFeedClient = &SportsFeed{
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

type NFLBoxScoreResponse struct {
	LastUpdatedOn time.Time `json:"lastUpdatedOn"`
	Game          struct {
		ID        int         `json:"id"`
		Week      int         `json:"week"`
		StartTime time.Time   `json:"startTime"`
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
		VenueAllegiance          string        `json:"venueAllegiance"`
		ScheduleStatus           string        `json:"scheduleStatus"`
		OriginalStartTime        interface{}   `json:"originalStartTime"`
		DelayedOrPostponedReason interface{}   `json:"delayedOrPostponedReason"`
		PlayedStatus             string        `json:"playedStatus"`
		Attendance               interface{}   `json:"attendance"`
		Officials                []interface{} `json:"officials"`
		Broadcasters             []string      `json:"broadcasters"`
		Weather                  struct {
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
	} `json:"game"`
	Scoring struct {
		CurrentQuarter                 int         `json:"currentQuarter"`
		CurrentQuarterSecondsRemaining int         `json:"currentQuarterSecondsRemaining"`
		CurrentIntermission            interface{} `json:"currentIntermission"`
		TeamInPossession               struct {
			ID           int    `json:"id"`
			Abbreviation string `json:"abbreviation"`
		} `json:"teamInPossession"`
		CurrentDown           int `json:"currentDown"`
		CurrentYardsRemaining int `json:"currentYardsRemaining"`
		LineOfScrimmage       struct {
			Team struct {
				ID           int    `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
			YardLine int `json:"yardLine"`
		} `json:"lineOfScrimmage"`
		AwayScoreTotal int `json:"awayScoreTotal"`
		HomeScoreTotal int `json:"homeScoreTotal"`
		Quarters       []struct {
			QuarterNumber int `json:"quarterNumber"`
			AwayScore     int `json:"awayScore"`
			HomeScore     int `json:"homeScore"`
			ScoringPlays  []struct {
				QuarterSecondsElapsed int `json:"quarterSecondsElapsed"`
				Team                  struct {
					ID           int    `json:"id"`
					Abbreviation string `json:"abbreviation"`
				} `json:"team"`
				ScoreChange     int    `json:"scoreChange"`
				AwayScore       int    `json:"awayScore"`
				HomeScore       int    `json:"homeScore"`
				PlayDescription string `json:"playDescription"`
			} `json:"scoringPlays"`
		} `json:"quarters"`
	} `json:"scoring"`
	Stats struct {
		Away struct {
			TeamStats []struct {
				Passing struct {
					PassAttempts    int     `json:"passAttempts"`
					PassCompletions int     `json:"passCompletions"`
					PassPct         float64 `json:"passPct"`
					PassGrossYards  int     `json:"passGrossYards"`
					PassNetYards    int     `json:"passNetYards"`
					PassAvg         float64 `json:"passAvg"`
					PassYardsPerAtt float64 `json:"passYardsPerAtt"`
					PassTD          int     `json:"passTD"`
					PassTDPct       float64 `json:"passTDPct"`
					PassInt         int     `json:"passInt"`
					PassIntPct      float64 `json:"passIntPct"`
					PassLng         int     `json:"passLng"`
					Pass20Plus      int     `json:"pass20Plus"`
					Pass40Plus      int     `json:"pass40Plus"`
					PassSacks       int     `json:"passSacks"`
					PassSackY       int     `json:"passSackY"`
					QbRating        float64 `json:"qbRating"`
				} `json:"passing"`
				Rushing struct {
					RushAttempts    int     `json:"rushAttempts"`
					RushYards       int     `json:"rushYards"`
					RushAverage     float64 `json:"rushAverage"`
					RushTD          int     `json:"rushTD"`
					RushLng         int     `json:"rushLng"`
					Rush1StDowns    int     `json:"rush1stDowns"`
					Rush1StDownsPct float64 `json:"rush1stDownsPct"`
					Rush20Plus      int     `json:"rush20Plus"`
					Rush40Plus      int     `json:"rush40Plus"`
					RushFumbles     int     `json:"rushFumbles"`
				} `json:"rushing"`
				Receiving struct {
					Receptions  int     `json:"receptions"`
					RecYards    int     `json:"recYards"`
					RecAverage  float64 `json:"recAverage"`
					RecTD       int     `json:"recTD"`
					RecLng      int     `json:"recLng"`
					Rec1StDowns int     `json:"rec1stDowns"`
					Rec20Plus   int     `json:"rec20Plus"`
					Rec40Plus   int     `json:"rec40Plus"`
					RecFumbles  int     `json:"recFumbles"`
				} `json:"receiving"`
				Tackles struct {
					TackleSolo     int `json:"tackleSolo"`
					TackleTotal    int `json:"tackleTotal"`
					TackleAst      int `json:"tackleAst"`
					Sacks          int `json:"sacks"`
					SackYds        int `json:"sackYds"`
					TacklesForLoss int `json:"tacklesForLoss"`
				} `json:"tackles"`
				Interceptions struct {
					Interceptions  int     `json:"interceptions"`
					IntTD          int     `json:"intTD"`
					IntYds         int     `json:"intYds"`
					IntAverage     float64 `json:"intAverage"`
					IntLng         int     `json:"intLng"`
					PassesDefended int     `json:"passesDefended"`
					Stuffs         int     `json:"stuffs"`
					StuffYds       int     `json:"stuffYds"`
					KB             int     `json:"kB"`
					Safeties       int     `json:"safeties"`
				} `json:"interceptions"`
				Fumbles struct {
					Fumbles     int `json:"fumbles"`
					FumLost     int `json:"fumLost"`
					FumForced   int `json:"fumForced"`
					FumOwnRec   int `json:"fumOwnRec"`
					FumOppRec   int `json:"fumOppRec"`
					FumRecYds   int `json:"fumRecYds"`
					FumTotalRec int `json:"fumTotalRec"`
					FumTD       int `json:"fumTD"`
				} `json:"fumbles"`
				KickoffReturns struct {
					KrRet    int     `json:"krRet"`
					KrYds    int     `json:"krYds"`
					KrAvg    float64 `json:"krAvg"`
					KrLng    int     `json:"krLng"`
					KrTD     int     `json:"krTD"`
					Kr20Plus int     `json:"kr20Plus"`
					Kr40Plus int     `json:"kr40Plus"`
					KrFC     int     `json:"krFC"`
					KrFum    int     `json:"krFum"`
				} `json:"kickoffReturns"`
				PuntReturns struct {
					PrRet    int     `json:"prRet"`
					PrYds    int     `json:"prYds"`
					PrAvg    float64 `json:"prAvg"`
					PrLng    int     `json:"prLng"`
					PrTD     int     `json:"prTD"`
					Pr20Plus int     `json:"pr20Plus"`
					Pr40Plus int     `json:"pr40Plus"`
					PrFC     int     `json:"prFC"`
					PrFum    int     `json:"prFum"`
				} `json:"puntReturns"`
				FieldGoals struct {
					FgBlk        int     `json:"fgBlk"`
					FgMade       int     `json:"fgMade"`
					FgAtt        int     `json:"fgAtt"`
					FgPct        float64 `json:"fgPct"`
					FgMade119    int     `json:"fgMade1_19"`
					FgAtt119     int     `json:"fgAtt1_19"`
					Fg119Pct     float64 `json:"fg1_19Pct"`
					FgMade2029   int     `json:"fgMade20_29"`
					FgAtt2029    int     `json:"fgAtt20_29"`
					Fg2029Pct    float64 `json:"fg20_29Pct"`
					FgMade3039   int     `json:"fgMade30_39"`
					FgAtt3039    int     `json:"fgAtt30_39"`
					Fg3039Pct    float64 `json:"fg30_39Pct"`
					FgMade4049   int     `json:"fgMade40_49"`
					FgAtt4049    int     `json:"fgAtt40_49"`
					Fg4049Pct    float64 `json:"fg40_49Pct"`
					FgMade50Plus int     `json:"fgMade50Plus"`
					FgAtt50Plus  int     `json:"fgAtt50Plus"`
					Fg50PlusPct  float64 `json:"fg50PlusPct"`
					FgLng        int     `json:"fgLng"`
				} `json:"fieldGoals"`
				ExtraPointAttempt struct {
					XpBlk      int     `json:"xpBlk"`
					XpMade     int     `json:"xpMade"`
					XpAtt      int     `json:"xpAtt"`
					XpPct      float64 `json:"xpPct"`
					FgAndXpPts int     `json:"fgAndXpPts"`
				} `json:"extraPointAttempt"`
				Kickoffs struct {
					Kickoffs    int     `json:"kickoffs"`
					KoYds       int     `json:"koYds"`
					KoOOB       int     `json:"koOOB"`
					KoAvg       float64 `json:"koAvg"`
					KoTB        int     `json:"koTB"`
					KoRet       int     `json:"koRet"`
					KoRetYds    int     `json:"koRetYds"`
					KoRetAvgYds float64 `json:"koRetAvgYds"`
					KoTD        int     `json:"koTD"`
					KoOS        int     `json:"koOS"`
					KoOSR       int     `json:"koOSR"`
				} `json:"kickoffs"`
				Punting struct {
					Punts       int     `json:"punts"`
					PuntYds     int     `json:"puntYds"`
					PuntNetYds  int     `json:"puntNetYds"`
					PuntLng     int     `json:"puntLng"`
					PuntAvg     float64 `json:"puntAvg"`
					PuntNetAvg  float64 `json:"puntNetAvg"`
					PuntBlk     int     `json:"puntBlk"`
					PuntOOB     int     `json:"puntOOB"`
					PuntDown    int     `json:"puntDown"`
					PuntIn20    int     `json:"puntIn20"`
					PuntIn20Pct float64 `json:"puntIn20Pct"`
					PuntTB      int     `json:"puntTB"`
					PuntTBPct   float64 `json:"puntTBPct"`
					PuntFC      int     `json:"puntFC"`
					PuntRet     int     `json:"puntRet"`
					PuntRetYds  int     `json:"puntRetYds"`
					PuntRetAvg  float64 `json:"puntRetAvg"`
				} `json:"punting"`
				Miscellaneous struct {
					FirstDownsTotal   int     `json:"firstDownsTotal"`
					FirstDownsPass    int     `json:"firstDownsPass"`
					FirstDownsRush    int     `json:"firstDownsRush"`
					FirstDownsPenalty int     `json:"firstDownsPenalty"`
					ThirdDowns        int     `json:"thirdDowns"`
					ThirdDownsAtt     int     `json:"thirdDownsAtt"`
					ThirdDownsPct     float64 `json:"thirdDownsPct"`
					FourthDowns       int     `json:"fourthDowns"`
					FourthDownsAtt    int     `json:"fourthDownsAtt"`
					FourthDownsPct    float64 `json:"fourthDownsPct"`
					Penalties         int     `json:"penalties"`
					PenaltyYds        int     `json:"penaltyYds"`
					OffensePlays      int     `json:"offensePlays"`
					OffenseYds        int     `json:"offenseYds"`
					OffenseAvgYds     float64 `json:"offenseAvgYds"`
					TotalTD           int     `json:"totalTD"`
				} `json:"miscellaneous"`
				Standings struct {
					Wins              int     `json:"wins"`
					Losses            int     `json:"losses"`
					Ties              int     `json:"ties"`
					OtWins            int     `json:"otWins"`
					OtLosses          int     `json:"otLosses"`
					WinPct            float64 `json:"winPct"`
					PointsFor         int     `json:"pointsFor"`
					PointsAgainst     int     `json:"pointsAgainst"`
					PointDifferential int     `json:"pointDifferential"`
				} `json:"standings"`
				TwoPointAttempts struct {
					TwoPtAtt      int `json:"twoPtAtt"`
					TwoPtMade     int `json:"twoPtMade"`
					TwoPtPassAtt  int `json:"twoPtPassAtt"`
					TwoPtPassMade int `json:"twoPtPassMade"`
					TwoPtRushAtt  int `json:"twoPtRushAtt"`
					TwoPtRushMade int `json:"twoPtRushMade"`
				} `json:"twoPointAttempts"`
				SnapCounts struct {
					OffenseSnaps     int `json:"offenseSnaps"`
					DefenseSnaps     int `json:"defenseSnaps"`
					SpecialTeamSnaps int `json:"specialTeamSnaps"`
				} `json:"snapCounts"`
			} `json:"teamStats"`
			Players []struct {
				Player struct {
					ID           int    `json:"id"`
					FirstName    string `json:"firstName"`
					LastName     string `json:"lastName"`
					Position     string `json:"position"`
					JerseyNumber int    `json:"jerseyNumber"`
				} `json:"player"`
				PlayerStats []struct {
					Tackles struct {
						TackleSolo     int     `json:"tackleSolo"`
						TackleTotal    int     `json:"tackleTotal"`
						TackleAst      int     `json:"tackleAst"`
						Sacks          float64 `json:"sacks"`
						SackYds        int     `json:"sackYds"`
						TacklesForLoss int     `json:"tacklesForLoss"`
					} `json:"tackles"`
					Interceptions struct {
						Interceptions  int     `json:"interceptions"`
						IntTD          int     `json:"intTD"`
						IntYds         int     `json:"intYds"`
						IntAverage     float64 `json:"intAverage"`
						IntLng         int     `json:"intLng"`
						PassesDefended int     `json:"passesDefended"`
						Stuffs         int     `json:"stuffs"`
						StuffYds       int     `json:"stuffYds"`
						Safeties       int     `json:"safeties"`
						KB             int     `json:"kB"`
					} `json:"interceptions"`
					Fumbles struct {
						Fumbles     int `json:"fumbles"`
						FumLost     int `json:"fumLost"`
						FumForced   int `json:"fumForced"`
						FumOwnRec   int `json:"fumOwnRec"`
						FumOppRec   int `json:"fumOppRec"`
						FumRecYds   int `json:"fumRecYds"`
						FumTotalRec int `json:"fumTotalRec"`
						FumTD       int `json:"fumTD"`
					} `json:"fumbles"`
					KickoffReturns struct {
						KrRet    int     `json:"krRet"`
						KrYds    int     `json:"krYds"`
						KrAvg    float64 `json:"krAvg"`
						KrLng    int     `json:"krLng"`
						KrTD     int     `json:"krTD"`
						Kr20Plus int     `json:"kr20Plus"`
						Kr40Plus int     `json:"kr40Plus"`
						KrFC     int     `json:"krFC"`
						KrFum    int     `json:"krFum"`
					} `json:"kickoffReturns"`
					PuntReturns struct {
						PrRet    int     `json:"prRet"`
						PrYds    int     `json:"prYds"`
						PrAvg    float64 `json:"prAvg"`
						PrLng    int     `json:"prLng"`
						PrTD     int     `json:"prTD"`
						Pr20Plus int     `json:"pr20Plus"`
						Pr40Plus int     `json:"pr40Plus"`
						PrFC     int     `json:"prFC"`
						PrFum    int     `json:"prFum"`
					} `json:"puntReturns"`
					FieldGoals struct {
						FgBlk        int     `json:"fgBlk"`
						FgMade       int     `json:"fgMade"`
						FgAtt        int     `json:"fgAtt"`
						FgPct        float64 `json:"fgPct"`
						FgMade119    int     `json:"fgMade1_19"`
						FgAtt119     int     `json:"fgAtt1_19"`
						Fg119Pct     float64 `json:"fg1_19Pct"`
						FgMade2029   int     `json:"fgMade20_29"`
						FgAtt2029    int     `json:"fgAtt20_29"`
						Fg2029Pct    float64 `json:"fg20_29Pct"`
						FgMade3039   int     `json:"fgMade30_39"`
						FgAtt3039    int     `json:"fgAtt30_39"`
						Fg3039Pct    float64 `json:"fg30_39Pct"`
						FgMade4049   int     `json:"fgMade40_49"`
						FgAtt4049    int     `json:"fgAtt40_49"`
						Fg4049Pct    float64 `json:"fg40_49Pct"`
						FgMade50Plus int     `json:"fgMade50Plus"`
						FgAtt50Plus  int     `json:"fgAtt50Plus"`
						Fg50PlusPct  float64 `json:"fg50PlusPct"`
						FgLng        int     `json:"fgLng"`
					} `json:"fieldGoals"`
					ExtraPointAttempts struct {
						XpBlk      int     `json:"xpBlk"`
						XpMade     int     `json:"xpMade"`
						XpAtt      int     `json:"xpAtt"`
						XpPct      float64 `json:"xpPct"`
						FgAndXpPts int     `json:"fgAndXpPts"`
					} `json:"extraPointAttempts"`
					Kickoffs struct {
						Kickoffs    int     `json:"kickoffs"`
						KoYds       int     `json:"koYds"`
						KoOOB       int     `json:"koOOB"`
						KoAvg       float64 `json:"koAvg"`
						KoTB        int     `json:"koTB"`
						KoRet       int     `json:"koRet"`
						KoRetYds    int     `json:"koRetYds"`
						KoRetAvgYds float64 `json:"koRetAvgYds"`
						KoTD        int     `json:"koTD"`
						KoOS        int     `json:"koOS"`
						KoOSR       int     `json:"koOSR"`
					} `json:"kickoffs"`
					Punting struct {
						Punts       int     `json:"punts"`
						PuntYds     int     `json:"puntYds"`
						PuntNetYds  int     `json:"puntNetYds"`
						PuntLng     int     `json:"puntLng"`
						PuntAvg     float64 `json:"puntAvg"`
						PuntNetAvg  float64 `json:"puntNetAvg"`
						PuntBlk     int     `json:"puntBlk"`
						PuntOOB     int     `json:"puntOOB"`
						PuntDown    int     `json:"puntDown"`
						PuntIn20    int     `json:"puntIn20"`
						PuntIn20Pct float64 `json:"puntIn20Pct"`
						PuntTB      int     `json:"puntTB"`
						PuntTBPct   float64 `json:"puntTBPct"`
						PuntFC      int     `json:"puntFC"`
						PuntRet     int     `json:"puntRet"`
						PuntRetYds  int     `json:"puntRetYds"`
						PuntRetAvg  float64 `json:"puntRetAvg"`
					} `json:"punting"`
					Miscellaneous struct {
						GamesStarted int `json:"gamesStarted"`
					} `json:"miscellaneous"`
					SnapCounts struct {
						OffenseSnaps     int `json:"offenseSnaps"`
						DefenseSnaps     int `json:"defenseSnaps"`
						SpecialTeamSnaps int `json:"specialTeamSnaps"`
					} `json:"snapCounts"`
				} `json:"playerStats"`
			} `json:"players"`
		} `json:"away"`
		Home struct {
			TeamStats []struct {
				Passing struct {
					PassAttempts    int     `json:"passAttempts"`
					PassCompletions int     `json:"passCompletions"`
					PassPct         float64 `json:"passPct"`
					PassGrossYards  int     `json:"passGrossYards"`
					PassNetYards    int     `json:"passNetYards"`
					PassAvg         float64 `json:"passAvg"`
					PassYardsPerAtt float64 `json:"passYardsPerAtt"`
					PassTD          int     `json:"passTD"`
					PassTDPct       float64 `json:"passTDPct"`
					PassInt         int     `json:"passInt"`
					PassIntPct      float64 `json:"passIntPct"`
					PassLng         int     `json:"passLng"`
					Pass20Plus      int     `json:"pass20Plus"`
					Pass40Plus      int     `json:"pass40Plus"`
					PassSacks       int     `json:"passSacks"`
					PassSackY       int     `json:"passSackY"`
					QbRating        float64 `json:"qbRating"`
				} `json:"passing"`
				Rushing struct {
					RushAttempts    int     `json:"rushAttempts"`
					RushYards       int     `json:"rushYards"`
					RushAverage     float64 `json:"rushAverage"`
					RushTD          int     `json:"rushTD"`
					RushLng         int     `json:"rushLng"`
					Rush1StDowns    int     `json:"rush1stDowns"`
					Rush1StDownsPct float64 `json:"rush1stDownsPct"`
					Rush20Plus      int     `json:"rush20Plus"`
					Rush40Plus      int     `json:"rush40Plus"`
					RushFumbles     int     `json:"rushFumbles"`
				} `json:"rushing"`
				Receiving struct {
					Receptions  int     `json:"receptions"`
					RecYards    int     `json:"recYards"`
					RecAverage  float64 `json:"recAverage"`
					RecTD       int     `json:"recTD"`
					RecLng      int     `json:"recLng"`
					Rec1StDowns int     `json:"rec1stDowns"`
					Rec20Plus   int     `json:"rec20Plus"`
					Rec40Plus   int     `json:"rec40Plus"`
					RecFumbles  int     `json:"recFumbles"`
				} `json:"receiving"`
				Tackles struct {
					TackleSolo     int `json:"tackleSolo"`
					TackleTotal    int `json:"tackleTotal"`
					TackleAst      int `json:"tackleAst"`
					Sacks          int `json:"sacks"`
					SackYds        int `json:"sackYds"`
					TacklesForLoss int `json:"tacklesForLoss"`
				} `json:"tackles"`
				Interceptions struct {
					Interceptions  int     `json:"interceptions"`
					IntTD          int     `json:"intTD"`
					IntYds         int     `json:"intYds"`
					IntAverage     float64 `json:"intAverage"`
					IntLng         int     `json:"intLng"`
					PassesDefended int     `json:"passesDefended"`
					Stuffs         int     `json:"stuffs"`
					StuffYds       int     `json:"stuffYds"`
					KB             int     `json:"kB"`
					Safeties       int     `json:"safeties"`
				} `json:"interceptions"`
				Fumbles struct {
					Fumbles     int `json:"fumbles"`
					FumLost     int `json:"fumLost"`
					FumForced   int `json:"fumForced"`
					FumOwnRec   int `json:"fumOwnRec"`
					FumOppRec   int `json:"fumOppRec"`
					FumRecYds   int `json:"fumRecYds"`
					FumTotalRec int `json:"fumTotalRec"`
					FumTD       int `json:"fumTD"`
				} `json:"fumbles"`
				KickoffReturns struct {
					KrRet    int     `json:"krRet"`
					KrYds    int     `json:"krYds"`
					KrAvg    float64 `json:"krAvg"`
					KrLng    int     `json:"krLng"`
					KrTD     int     `json:"krTD"`
					Kr20Plus int     `json:"kr20Plus"`
					Kr40Plus int     `json:"kr40Plus"`
					KrFC     int     `json:"krFC"`
					KrFum    int     `json:"krFum"`
				} `json:"kickoffReturns"`
				PuntReturns struct {
					PrRet    int     `json:"prRet"`
					PrYds    int     `json:"prYds"`
					PrAvg    float64 `json:"prAvg"`
					PrLng    int     `json:"prLng"`
					PrTD     int     `json:"prTD"`
					Pr20Plus int     `json:"pr20Plus"`
					Pr40Plus int     `json:"pr40Plus"`
					PrFC     int     `json:"prFC"`
					PrFum    int     `json:"prFum"`
				} `json:"puntReturns"`
				FieldGoals struct {
					FgBlk        int     `json:"fgBlk"`
					FgMade       int     `json:"fgMade"`
					FgAtt        int     `json:"fgAtt"`
					FgPct        float64 `json:"fgPct"`
					FgMade119    int     `json:"fgMade1_19"`
					FgAtt119     int     `json:"fgAtt1_19"`
					Fg119Pct     float64 `json:"fg1_19Pct"`
					FgMade2029   int     `json:"fgMade20_29"`
					FgAtt2029    int     `json:"fgAtt20_29"`
					Fg2029Pct    float64 `json:"fg20_29Pct"`
					FgMade3039   int     `json:"fgMade30_39"`
					FgAtt3039    int     `json:"fgAtt30_39"`
					Fg3039Pct    float64 `json:"fg30_39Pct"`
					FgMade4049   int     `json:"fgMade40_49"`
					FgAtt4049    int     `json:"fgAtt40_49"`
					Fg4049Pct    float64 `json:"fg40_49Pct"`
					FgMade50Plus int     `json:"fgMade50Plus"`
					FgAtt50Plus  int     `json:"fgAtt50Plus"`
					Fg50PlusPct  float64 `json:"fg50PlusPct"`
					FgLng        int     `json:"fgLng"`
				} `json:"fieldGoals"`
				ExtraPointAttempt struct {
					XpBlk      int     `json:"xpBlk"`
					XpMade     int     `json:"xpMade"`
					XpAtt      int     `json:"xpAtt"`
					XpPct      float64 `json:"xpPct"`
					FgAndXpPts int     `json:"fgAndXpPts"`
				} `json:"extraPointAttempt"`
				Kickoffs struct {
					Kickoffs    int     `json:"kickoffs"`
					KoYds       int     `json:"koYds"`
					KoOOB       int     `json:"koOOB"`
					KoAvg       float64 `json:"koAvg"`
					KoTB        int     `json:"koTB"`
					KoRet       int     `json:"koRet"`
					KoRetYds    int     `json:"koRetYds"`
					KoRetAvgYds float64 `json:"koRetAvgYds"`
					KoTD        int     `json:"koTD"`
					KoOS        int     `json:"koOS"`
					KoOSR       int     `json:"koOSR"`
				} `json:"kickoffs"`
				Punting struct {
					Punts       int     `json:"punts"`
					PuntYds     int     `json:"puntYds"`
					PuntNetYds  int     `json:"puntNetYds"`
					PuntLng     int     `json:"puntLng"`
					PuntAvg     float64 `json:"puntAvg"`
					PuntNetAvg  float64 `json:"puntNetAvg"`
					PuntBlk     int     `json:"puntBlk"`
					PuntOOB     int     `json:"puntOOB"`
					PuntDown    int     `json:"puntDown"`
					PuntIn20    int     `json:"puntIn20"`
					PuntIn20Pct float64 `json:"puntIn20Pct"`
					PuntTB      int     `json:"puntTB"`
					PuntTBPct   float64 `json:"puntTBPct"`
					PuntFC      int     `json:"puntFC"`
					PuntRet     int     `json:"puntRet"`
					PuntRetYds  int     `json:"puntRetYds"`
					PuntRetAvg  float64 `json:"puntRetAvg"`
				} `json:"punting"`
				Miscellaneous struct {
					FirstDownsTotal   int     `json:"firstDownsTotal"`
					FirstDownsPass    int     `json:"firstDownsPass"`
					FirstDownsRush    int     `json:"firstDownsRush"`
					FirstDownsPenalty int     `json:"firstDownsPenalty"`
					ThirdDowns        int     `json:"thirdDowns"`
					ThirdDownsAtt     int     `json:"thirdDownsAtt"`
					ThirdDownsPct     float64 `json:"thirdDownsPct"`
					FourthDowns       int     `json:"fourthDowns"`
					FourthDownsAtt    int     `json:"fourthDownsAtt"`
					FourthDownsPct    float64 `json:"fourthDownsPct"`
					Penalties         int     `json:"penalties"`
					PenaltyYds        int     `json:"penaltyYds"`
					OffensePlays      int     `json:"offensePlays"`
					OffenseYds        int     `json:"offenseYds"`
					OffenseAvgYds     float64 `json:"offenseAvgYds"`
					TotalTD           int     `json:"totalTD"`
				} `json:"miscellaneous"`
				Standings struct {
					Wins              int     `json:"wins"`
					Losses            int     `json:"losses"`
					Ties              int     `json:"ties"`
					OtWins            int     `json:"otWins"`
					OtLosses          int     `json:"otLosses"`
					WinPct            float64 `json:"winPct"`
					PointsFor         int     `json:"pointsFor"`
					PointsAgainst     int     `json:"pointsAgainst"`
					PointDifferential int     `json:"pointDifferential"`
				} `json:"standings"`
				TwoPointAttempts struct {
					TwoPtAtt      int `json:"twoPtAtt"`
					TwoPtMade     int `json:"twoPtMade"`
					TwoPtPassAtt  int `json:"twoPtPassAtt"`
					TwoPtPassMade int `json:"twoPtPassMade"`
					TwoPtRushAtt  int `json:"twoPtRushAtt"`
					TwoPtRushMade int `json:"twoPtRushMade"`
				} `json:"twoPointAttempts"`
				SnapCounts struct {
					OffenseSnaps     int `json:"offenseSnaps"`
					DefenseSnaps     int `json:"defenseSnaps"`
					SpecialTeamSnaps int `json:"specialTeamSnaps"`
				} `json:"snapCounts"`
			} `json:"teamStats"`
			Players []struct {
				Player struct {
					ID           int    `json:"id"`
					FirstName    string `json:"firstName"`
					LastName     string `json:"lastName"`
					Position     string `json:"position"`
					JerseyNumber int    `json:"jerseyNumber"`
				} `json:"player"`
				PlayerStats []struct {
					Tackles struct {
						TackleSolo     int     `json:"tackleSolo"`
						TackleTotal    int     `json:"tackleTotal"`
						TackleAst      int     `json:"tackleAst"`
						Sacks          float64 `json:"sacks"`
						SackYds        int     `json:"sackYds"`
						TacklesForLoss int     `json:"tacklesForLoss"`
					} `json:"tackles"`
					Interceptions struct {
						Interceptions  int     `json:"interceptions"`
						IntTD          int     `json:"intTD"`
						IntYds         int     `json:"intYds"`
						IntAverage     float64 `json:"intAverage"`
						IntLng         int     `json:"intLng"`
						PassesDefended int     `json:"passesDefended"`
						Stuffs         int     `json:"stuffs"`
						StuffYds       int     `json:"stuffYds"`
						Safeties       int     `json:"safeties"`
						KB             int     `json:"kB"`
					} `json:"interceptions"`
					Fumbles struct {
						Fumbles     int `json:"fumbles"`
						FumLost     int `json:"fumLost"`
						FumForced   int `json:"fumForced"`
						FumOwnRec   int `json:"fumOwnRec"`
						FumOppRec   int `json:"fumOppRec"`
						FumRecYds   int `json:"fumRecYds"`
						FumTotalRec int `json:"fumTotalRec"`
						FumTD       int `json:"fumTD"`
					} `json:"fumbles"`
					KickoffReturns struct {
						KrRet    int     `json:"krRet"`
						KrYds    int     `json:"krYds"`
						KrAvg    float64 `json:"krAvg"`
						KrLng    int     `json:"krLng"`
						KrTD     int     `json:"krTD"`
						Kr20Plus int     `json:"kr20Plus"`
						Kr40Plus int     `json:"kr40Plus"`
						KrFC     int     `json:"krFC"`
						KrFum    int     `json:"krFum"`
					} `json:"kickoffReturns"`
					PuntReturns struct {
						PrRet    int     `json:"prRet"`
						PrYds    int     `json:"prYds"`
						PrAvg    float64 `json:"prAvg"`
						PrLng    int     `json:"prLng"`
						PrTD     int     `json:"prTD"`
						Pr20Plus int     `json:"pr20Plus"`
						Pr40Plus int     `json:"pr40Plus"`
						PrFC     int     `json:"prFC"`
						PrFum    int     `json:"prFum"`
					} `json:"puntReturns"`
					Miscellaneous struct {
						GamesStarted int `json:"gamesStarted"`
					} `json:"miscellaneous"`
					SnapCounts struct {
						OffenseSnaps     int `json:"offenseSnaps"`
						DefenseSnaps     int `json:"defenseSnaps"`
						SpecialTeamSnaps int `json:"specialTeamSnaps"`
					} `json:"snapCounts"`
				} `json:"playerStats"`
			} `json:"players"`
		} `json:"home"`
	} `json:"stats"`
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
		PlayerReferences []struct {
			ID              int    `json:"id"`
			FirstName       string `json:"firstName"`
			LastName        string `json:"lastName"`
			PrimaryPosition string `json:"primaryPosition"`
			JerseyNumber    int    `json:"jerseyNumber"`
			CurrentTeam     struct {
				ID           int    `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"currentTeam"`
			CurrentRosterStatus string        `json:"currentRosterStatus"`
			CurrentInjury       interface{}   `json:"currentInjury"`
			Height              string        `json:"height"`
			Weight              int           `json:"weight"`
			BirthDate           string        `json:"birthDate"`
			Age                 int           `json:"age"`
			BirthCity           interface{}   `json:"birthCity"`
			BirthCountry        interface{}   `json:"birthCountry"`
			Rookie              bool          `json:"rookie"`
			HighSchool          interface{}   `json:"highSchool"`
			College             string        `json:"college"`
			Handedness          interface{}   `json:"handedness"`
			OfficialImageSrc    string        `json:"officialImageSrc"`
			SocialMediaAccounts []interface{} `json:"socialMediaAccounts"`
		} `json:"playerReferences"`
		PlayerStatReferences []struct {
			Category     string `json:"category"`
			FullName     string `json:"fullName"`
			Description  string `json:"description"`
			Abbreviation string `json:"abbreviation"`
			Type         string `json:"type"`
		} `json:"playerStatReferences"`
		TeamStatReferences []struct {
			Category     string `json:"category"`
			FullName     string `json:"fullName"`
			Description  string `json:"description"`
			Abbreviation string `json:"abbreviation"`
			Type         string `json:"type"`
		} `json:"teamStatReferences"`
	} `json:"references"`
}

type NFLBoxScoreResponseFormatted struct {
	HomeAbbreviation			string
	AwayAbbreviation			string
	HomeScore					int
	AwayScore					int
	HomeLogo					string
	AwayLogo					string
	Quarter						int
	QuarterMinRemaining			int
	QuarterSecRemaining			int
	Down						int
	YardsRemaining				int
	LineOfScrimmage				int	
	PlayedStatus				string
	StartTime					time.Time
	HomePassYards				int
	AwayPassYards				int
	HomeRushYards				int
	AwayRushYards				int
	HomeSacks					int
	AwaySacks					int
	HomeWins					int
	AwayWins					int
	HomeLosses					int
	AwayLosses					int
	HomeTies					int
	AwayTies					int
}

func (s *SportsFeed) FetchNFLBoxScore(game string) NFLBoxScoreResponseFormatted {
	// determine the "season" to use
	year, err := strconv.Atoi(game[:4])
	checkErr(err)
	month, err := strconv.Atoi(game[4:6])
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
		Endpoint: fmt.Sprintf("/pull/nfl/%s/games/%s/boxscore.json", season, game),
	}

	var responseData NFLBoxScoreResponse
	if err := s.Client.DoAndUnmarshal(request, &responseData); err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	// clean up time remaining
	secRem := responseData.Scoring.CurrentQuarterSecondsRemaining 
	sec := int(secRem % 60)
	min := int(secRem / 60)

	formattedGameData := NFLBoxScoreResponseFormatted{
		HomeAbbreviation: responseData.Game.HomeTeam.Abbreviation,
		AwayAbbreviation: responseData.Game.AwayTeam.Abbreviation,
		HomeScore: responseData.Scoring.HomeScoreTotal,
		AwayScore: responseData.Scoring.AwayScoreTotal,
		HomeLogo: responseData.References.TeamReferences[0].OfficialLogoImageSrc,
		AwayLogo: responseData.References.TeamReferences[1].OfficialLogoImageSrc,
		Quarter: responseData.Scoring.CurrentQuarter,
		QuarterMinRemaining: min,
		QuarterSecRemaining: sec,
		Down: responseData.Scoring.CurrentDown,
		YardsRemaining: responseData.Scoring.CurrentYardsRemaining,
		LineOfScrimmage: responseData.Scoring.LineOfScrimmage.YardLine,
		PlayedStatus: responseData.Game.PlayedStatus,
		StartTime: responseData.Game.StartTime,
		HomePassYards: responseData.Stats.Home.TeamStats[0].Passing.PassNetYards,
		AwayPassYards: responseData.Stats.Away.TeamStats[0].Passing.PassNetYards,
		HomeRushYards: responseData.Stats.Home.TeamStats[0].Rushing.RushYards,
		AwayRushYards: responseData.Stats.Away.TeamStats[0].Rushing.RushYards,
		HomeSacks: responseData.Stats.Home.TeamStats[0].Tackles.Sacks,
		AwaySacks: responseData.Stats.Away.TeamStats[0].Tackles.Sacks,
		HomeWins: responseData.Stats.Home.TeamStats[0].Standings.Wins,
		AwayWins: responseData.Stats.Away.TeamStats[0].Standings.Wins,
		HomeLosses: responseData.Stats.Home.TeamStats[0].Standings.Losses,
		AwayLosses: responseData.Stats.Away.TeamStats[0].Standings.Losses,
		HomeTies: responseData.Stats.Home.TeamStats[0].Standings.Ties,
		AwayTies: responseData.Stats.Away.TeamStats[0].Standings.Ties,

	}

	return formattedGameData
}



func checkErr(e error) {
	if (e != nil) {
		log.Printf("Error: %s", e)
	}
}