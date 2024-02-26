package model

import (
	"time"
	"strconv"
	data "github.com/6ixisgood/matrix-ticker/pkg/data"
)

type Matchup struct {
	Teams		[]Team
	StartTime	time.Time
	EndTime		time.Time
	Period		string
	PeriodMinRemaining string
	PeriodSecRemaining string
	PlayedStatus	string
	Extra			string
}

type Team struct {
	Nickname		string
	Abbreviation	string
	Score			string
	Location		string
	Wins			string
	Loses			string
	Ties			string
	Points			string
	Record			string
	LogoSrc			string
}

// FetchLeagueMatchupsByDate given a league and date, fetch the matchups for the day
func FetchLeagueMatchupsByDate(league string, date time.Time) []Matchup {
	client := data.SportsFeedClient()
	var matchups []Matchup
	switch league {
	case "nfl":
		resp := client.FetchDailyNFLGamesInfo(date)	

		// map team id to reference data for lookups
		idToRef := map[int]data.DailyGamesNFLTeamReference{}
		for _, ref := range(resp.References.TeamReferences) {
			idToRef[ref.ID] = ref
		}

		for _, game := range(resp.Games) {
			team1 := Team{
				Nickname: idToRef[game.Schedule.AwayTeam.ID].Name,
				Abbreviation: idToRef[game.Schedule.AwayTeam.ID].Abbreviation,
				Score: strconv.Itoa(game.Score.AwayScoreTotal),
				Location: idToRef[game.Schedule.AwayTeam.ID].City,
				Record: "",
				LogoSrc: idToRef[game.Schedule.AwayTeam.ID].OfficialLogoImageSrc,
			}
			team2 := Team{
				Nickname: idToRef[game.Schedule.HomeTeam.ID].Name,
				Abbreviation: idToRef[game.Schedule.HomeTeam.ID].Abbreviation,
				Score: strconv.Itoa(game.Score.HomeScoreTotal),
				Location: idToRef[game.Schedule.HomeTeam.ID].City,
				Record: "",
				LogoSrc: idToRef[game.Schedule.HomeTeam.ID].OfficialLogoImageSrc,
			}
			matchup := Matchup{
				Teams: []Team{team1, team2},
				StartTime: game.Schedule.StartTime,
				EndTime: game.Schedule.EndedTime,
				PlayedStatus: game.Schedule.PlayedStatus,
				Period: strconv.Itoa(game.Score.CurrentQuarter),
				PeriodMinRemaining: strconv.Itoa(int(game.Score.CurrentQuarterSecondsRemaining / 60)),
				PeriodSecRemaining: strconv.Itoa(int(game.Score.CurrentQuarterSecondsRemaining % 60)),

			}
			matchups = append(matchups, matchup)
		}
	}
	return matchups
}