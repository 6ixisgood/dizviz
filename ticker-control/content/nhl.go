package content

import (
	"fmt"
	"github.com/sixisgoood/matrix-ticker/sports_data"
)

const (
	// AWAY_TEAM AWAY_SCORE - HOME_TEAM HOME_SCORE TIME_REMAINING PERIOD
	gameScoreFeed = "%v %v - %v %v %v %v"
)


func NHLDailyGamesTicker(season string, date string) string {
	data := sports_data.FetchDailyNHLGamesInfo(season, date)
	
	games := []string{}	

	for _, game := range data.Games {
		awayTeam := game.Schedule.AwayTeam.Abbreviation
		awayTeamScore := game.Score.AwayScoreTotal
		homeTeam := game.Schedule.HomeTeam.Abbreviation
		homeTeamScore := game.Score.HomeScoreTotal
		timeRemaining := game.Score.CurrentPeriodSecondsRemaining
		if timeRemaining == nil {
			timeRemaining = ""
		}
		currentPeriod := game.Score.CurrentPeriod
		if currentPeriod == nil && game.Schedule.PlayedStatus == "COMPLETED" {
			currentPeriod = "FINAL"
		} else if currentPeriod == nil {
			currentPeriod = ""
		}



		games = append(games, fmt.Sprintf(gameScoreFeed, 
											awayTeam,
											awayTeamScore,
											homeTeam,
											homeTeamScore,
											timeRemaining,
											currentPeriod))
	}

	feed := ""
	for _, gameString := range games {
		feed += gameString
		feed += "  \u2022  " 
	}

	return feed
} 