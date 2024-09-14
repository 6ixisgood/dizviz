package data

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SportsFeed struct {
	Client *APIClient
}

type SportsFeedConfig struct {
	BaseUrl  string
	Username string
	Password string
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

// sportsfeedDateFormat take in a time.Time and return the proper seasona and date format for sportsfeed
func sportsfeedDateFormat(date time.Time) (string, string) {
	// come up with date id for sportsfeed
	dateStr := date.Format("20060102")
	year, _ := strconv.Atoi(dateStr[:4])
	month, _ := strconv.Atoi(dateStr[4:6])

	var season string
	if month > 6 {
		season = fmt.Sprintf("%d-%d-regular", year, year+1)
	} else {
		season = fmt.Sprintf("%d-%d-regular", year-1, year)
	}

	return season, dateStr
}

// FetchDailyNHLGamesInfo get daily NHL games response from Sportsfeed
func (s *SportsFeed) FetchDailyNHLGamesInfo(date time.Time) DailyGamesNHLResponse {
	season, dateStr := sportsfeedDateFormat(date)

	// fetch the data
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/pull/nhl/%s/date/%s/games.json", season, dateStr),
	}

	var responseData DailyGamesNHLResponse
	_, err := s.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// FetchDailyNFLGamesInfo get daily NFL games response from Sportsfeed
func (s *SportsFeed) FetchDailyNFLGamesInfo(date time.Time) DailyGamesNFLResponse {
	season, dateStr := sportsfeedDateFormat(date)

	// fetch the data
	fmt.Println(season)
	fmt.Println(dateStr)
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/pull/nfl/%s/date/%s/games.json", season, dateStr),
	}

	var responseData DailyGamesNFLResponse
	_, err := s.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

func (s *SportsFeed) FetchNFLCurrentSeason() NFLCurrentSeasonResponse {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/pull/nfl/current_season.json"),
	}

	var responseData NFLCurrentSeasonResponse
	_, err := s.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

func (s *SportsFeed) FetchNFLCurrentWeek() NFLCurrentWeek {
	var week NFLCurrentWeek

	res := s.FetchNFLCurrentSeason()

	if len(res.Seasons) == 0 {
		return week
	}

	week.SeasonSlug = res.Seasons[0].Slug

	duration := time.Now().Sub(res.Seasons[0].StartDate.Time)
	days := int(duration.Hours() / 24)
	weeks := days / 7
	week.Week = fmt.Sprintf("%d", weeks+1)

	log.Println(days, weeks, week.Week)

	return week
}

func (s *SportsFeed) FetchNFLWeeklyGames(seasonSlug string, week string) NFLWeeklyGamesResponse {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/pull/nfl/%s/week/%s/games.json", seasonSlug, week),
	}

	var responseData NFLWeeklyGamesResponse
	_, err := s.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

func (s *SportsFeed) FetchNFLWeeklyGamesFormatted(seasonSlug string, week string) []NFLBoxScoreResponseFormatted {
	res := s.FetchNFLWeeklyGames(seasonSlug, week)

	var games []NFLBoxScoreResponseFormatted

	for _, game := range res.Games {
		secRem := game.Score.CurrentQuarterSecondsRemaining
		sec := int(secRem % 60)
		min := int(secRem / 60)

		games = append(games, NFLBoxScoreResponseFormatted{
			GameID:              fmt.Sprintf("%d", game.Schedule.ID),
			HomeAbbreviation:    game.Schedule.HomeTeam.Abbreviation,
			AwayAbbreviation:    game.Schedule.AwayTeam.Abbreviation,
			HomeScore:           game.Score.HomeScoreTotal,
			AwayScore:           game.Score.AwayScoreTotal,
			Quarter:             game.Score.CurrentQuarter,
			QuarterMinRemaining: min,
			QuarterSecRemaining: sec,
		})
	}
	return games
}

// FetchNFLBoxScore Fetch info about specific NFL game from Sportsfeed
func (s *SportsFeed) FetchNFLBoxScore(matchup string, date time.Time) (NFLBoxScoreResponseFormatted, error) {
	season, dateStr := sportsfeedDateFormat(date)

	var path string
	if strings.Contains(matchup, "-") {
		path = fmt.Sprintf("/pull/nfl/%s/games/%s-%s/boxscore.json", season, dateStr, matchup)
	} else {
		path = fmt.Sprintf("/pull/nfl/%s/games/%s/boxscore.json", season, matchup)
	}

	// fetch the data
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: path,
	}

	var responseData NFLBoxScoreResponse
	var formattedGameData NFLBoxScoreResponseFormatted
	statusCode, err := s.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}
	if statusCode == 204 {
		return formattedGameData, errors.New("204 received from SportsFeed, game isn't ready yet")
	}

	// clean up time remaining
	secRem := responseData.Scoring.CurrentQuarterSecondsRemaining
	sec := int(secRem % 60)
	min := int(secRem / 60)

	formattedGameData = NFLBoxScoreResponseFormatted{
		HomeAbbreviation:    responseData.Game.HomeTeam.Abbreviation,
		AwayAbbreviation:    responseData.Game.AwayTeam.Abbreviation,
		HomeScore:           responseData.Scoring.HomeScoreTotal,
		AwayScore:           responseData.Scoring.AwayScoreTotal,
		HomeLogo:            responseData.References.TeamReferences[0].OfficialLogoImageSrc,
		AwayLogo:            responseData.References.TeamReferences[1].OfficialLogoImageSrc,
		Quarter:             responseData.Scoring.CurrentQuarter,
		QuarterMinRemaining: min,
		QuarterSecRemaining: sec,
		Down:                responseData.Scoring.CurrentDown,
		YardsRemaining:      responseData.Scoring.CurrentYardsRemaining,
		LineOfScrimmage:     responseData.Scoring.LineOfScrimmage.YardLine,
		PlayedStatus:        responseData.Game.PlayedStatus,
		StartTime:           responseData.Game.StartTime,
		HomePassYards:       responseData.Stats.Home.TeamStats[0].Passing.PassNetYards,
		AwayPassYards:       responseData.Stats.Away.TeamStats[0].Passing.PassNetYards,
		HomeRushYards:       responseData.Stats.Home.TeamStats[0].Rushing.RushYards,
		AwayRushYards:       responseData.Stats.Away.TeamStats[0].Rushing.RushYards,
		HomeSacks:           responseData.Stats.Home.TeamStats[0].Tackles.Sacks,
		AwaySacks:           responseData.Stats.Away.TeamStats[0].Tackles.Sacks,
		HomeWins:            responseData.Stats.Home.TeamStats[0].Standings.Wins,
		AwayWins:            responseData.Stats.Away.TeamStats[0].Standings.Wins,
		HomeLosses:          responseData.Stats.Home.TeamStats[0].Standings.Losses,
		AwayLosses:          responseData.Stats.Away.TeamStats[0].Standings.Losses,
		HomeTies:            responseData.Stats.Home.TeamStats[0].Standings.Ties,
		AwayTies:            responseData.Stats.Away.TeamStats[0].Standings.Ties,
	}

	return formattedGameData, nil
}

func checkErr(e error) {
	if e != nil {
		log.Printf("Error: %s", e)
	}
}
