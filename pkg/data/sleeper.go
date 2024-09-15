package data

import (
	"fmt"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"sort"
)

type Sleeper struct {
	Client *APIClient
}

type SleeperConfig struct {
	BaseUrl string
}

// Sleeper Singleton
var sleeperClient *Sleeper

func SleeperClient() *Sleeper {
	return sleeperClient
}

func InitSleeperClient(config SleeperConfig) {
	clientOptions := APIClientOptions{
		BaseURL: config.BaseUrl,
	}

	client := NewAPIClient(clientOptions)

	// set the singleton
	sleeperClient = &Sleeper{
		Client: client,
	}

}

// cache
var sleperPlayers map[string]SleeperPlayer

type SleeperLeagueSettings struct {
	ReserveAllowCov          int `json:"reserve_allow_cov"`
	ReserveSlots             int `json:"reserve_slots"`
	Leg                      int `json:"leg"`
	OffseasonAdds            int `json:"offseason_adds"`
	BenchLock                int `json:"bench_lock"`
	TradeReviewDays          int `json:"trade_review_days"`
	LeagueAverageMatch       int `json:"league_average_match"`
	WaiverType               int `json:"waiver_type"`
	MaxKeepers               int `json:"max_keepers"`
	Type                     int `json:"type"`
	PickTrading              int `json:"pick_trading"`
	DisableTrades            int `json:"disable_trades"`
	DailyWaivers             int `json:"daily_waivers"`
	TaxiYears                int `json:"taxi_years"`
	TradeDeadline            int `json:"trade_deadline"`
	VetoShowVotes            int `json:"veto_show_votes"`
	ReserveAllowSus          int `json:"reserve_allow_sus"`
	ReserveAllowOut          int `json:"reserve_allow_out"`
	PlayoffRoundType         int `json:"playoff_round_type"`
	WaiverDayOfWeek          int `json:"waiver_day_of_week"`
	TaxiAllowVets            int `json:"taxi_allow_vets"`
	ReserveAllowDnr          int `json:"reserve_allow_dnr"`
	VetoAutoPoll             int `json:"veto_auto_poll"`
	CommissionerDirectInvite int `json:"commissioner_direct_invite"`
	ReserveAllowDoubtful     int `json:"reserve_allow_doubtful"`
	WaiverClearDays          int `json:"waiver_clear_days"`
	PlayoffWeekStart         int `json:"playoff_week_start"`
	DailyWaiversDays         int `json:"daily_waivers_days"`
	TaxiSlots                int `json:"taxi_slots"`
	PlayoffType              int `json:"playoff_type"`
	DailyWaiversHour         int `json:"daily_waivers_hour"`
	NumTeams                 int `json:"num_teams"`
	VetoVotesNeeded          int `json:"veto_votes_needed"`
	PlayoffTeams             int `json:"playoff_teams"`
	PlayoffSeedType          int `json:"playoff_seed_type"`
	StartWeek                int `json:"start_week"`
	ReserveAllowNa           int `json:"reserve_allow_na"`
	DraftRounds              int `json:"draft_rounds"`
	TaxiDeadline             int `json:"taxi_deadline"`
	WaiverBidMin             int `json:"waiver_bid_min"`
	CapacityOverride         int `json:"capacity_override"`
	Divisions                int `json:"divisions"`
	DisableAdds              int `json:"disable_adds"`
	WaiverBudget             int `json:"waiver_budget"`
	BestBall                 int `json:"best_ball"`
}

type SleeperLeagueMetadata struct {
	LatestLeagueWinnerRosterID string `json:"latest_league_winner_roster_id"`
	KeeperDeadline             string `json:"keeper_deadline"`
	Division2                  string `json:"division_2"`
	Division1                  string `json:"division_1"`
	AutoContinue               string `json:"auto_continue"`
}
type SleeperLeagueScoringSettings struct {
	StFf         float64 `json:"st_ff"`
	PtsAllow713  float64 `json:"pts_allow_7_13"`
	DefStFf      float64 `json:"def_st_ff"`
	RecYd        float64 `json:"rec_yd"`
	FumRecTd     float64 `json:"fum_rec_td"`
	PtsAllow35P  float64 `json:"pts_allow_35p"`
	PtsAllow2834 float64 `json:"pts_allow_28_34"`
	Fum          float64 `json:"fum"`
	RushYd       float64 `json:"rush_yd"`
	PassTd       float64 `json:"pass_td"`
	BlkKick      float64 `json:"blk_kick"`
	PassYd       float64 `json:"pass_yd"`
	Safe         float64 `json:"safe"`
	DefTd        float64 `json:"def_td"`
	Fgm50P       float64 `json:"fgm_50p"`
	DefStTd      float64 `json:"def_st_td"`
	FumRec       float64 `json:"fum_rec"`
	Rush2Pt      float64 `json:"rush_2pt"`
	Xpm          float64 `json:"xpm"`
	PtsAllow2127 float64 `json:"pts_allow_21_27"`
	Fgm2029      float64 `json:"fgm_20_29"`
	PtsAllow16   float64 `json:"pts_allow_1_6"`
	FumLost      float64 `json:"fum_lost"`
	DefStFumRec  float64 `json:"def_st_fum_rec"`
	Int          float64 `json:"int"`
	Fgm019       float64 `json:"fgm_0_19"`
	PtsAllow1420 float64 `json:"pts_allow_14_20"`
	Rec          float64 `json:"rec"`
	Ff           float64 `json:"ff"`
	Fgmiss       float64 `json:"fgmiss"`
	StFumRec     float64 `json:"st_fum_rec"`
	Rec2Pt       float64 `json:"rec_2pt"`
	RushTd       float64 `json:"rush_td"`
	Xpmiss       float64 `json:"xpmiss"`
	Fgm3039      float64 `json:"fgm_30_39"`
	RecTd        float64 `json:"rec_td"`
	StTd         float64 `json:"st_td"`
	Pass2Pt      float64 `json:"pass_2pt"`
	PtsAllow0    float64 `json:"pts_allow_0"`
	PassInt      float64 `json:"pass_int"`
	Fgm4049      float64 `json:"fgm_40_49"`
	Sack         float64 `json:"sack"`
}

type SleeperLeague struct {
	TotalRosters          int                          `json:"total_rosters"`
	Status                string                       `json:"status"`
	Sport                 string                       `json:"sport"`
	Shard                 int                          `json:"shard"`
	Settings              SleeperLeagueSettings        `json:"settings"`
	SeasonType            string                       `json:"season_type"`
	Season                string                       `json:"season"`
	ScoringSettings       SleeperLeagueScoringSettings `json:"scoring_settings"`
	RosterPositions       []string                     `json:"roster_positions"`
	PreviousLeagueID      string                       `json:"previous_league_id"`
	Name                  string                       `json:"name"`
	Metadata              SleeperLeagueMetadata        `json:"metadata"`
	LoserBracketID        int                          `json:"loser_bracket_id"`
	LeagueID              string                       `json:"league_id"`
	LastReadID            string                       `json:"last_read_id"`
	LastPinnedMessageID   string                       `json:"last_pinned_message_id"`
	LastMessageTime       int64                        `json:"last_message_time"`
	LastMessageTextMap    interface{}                  `json:"last_message_text_map"`
	LastMessageID         string                       `json:"last_message_id"`
	LastMessageAttachment interface{}                  `json:"last_message_attachment"`
	LastAuthorIsBot       bool                         `json:"last_author_is_bot"`
	LastAuthorID          string                       `json:"last_author_id"`
	LastAuthorDisplayName string                       `json:"last_author_display_name"`
	LastAuthorAvatar      interface{}                  `json:"last_author_avatar"`
	GroupID               string                       `json:"group_id"`
	DraftID               string                       `json:"draft_id"`
	CompanyID             string                       `json:"company_id"`
	BracketID             int                          `json:"bracket_id"`
	Avatar                interface{}                  `json:"avatar"`
}

type SleeperLeagueMatchup struct {
	StartersPoints []float64          `json:"starters_points"`
	Starters       []string           `json:"starters"`
	RosterID       int                `json:"roster_id"`
	Points         float64            `json:"points"`
	PlayersPoints  map[string]float64 `json:"players_points"`
	Players        []string           `json:"players"`
	MatchupID      int                `json:"matchup_id"`
	CustomPoints   interface{}        `json:"custom_points"`
}

type SleeperPlayer struct {
	Status                string      `json:"status"`
	EspnID                interface{} `json:"espn_id"`
	College               interface{} `json:"college"`
	SwishID               interface{} `json:"swish_id"`
	FantasyPositions      interface{} `json:"fantasy_positions"`
	Position              string      `json:"position"`
	FullName              string      `json:"full_name"`
	InjuryStatus          interface{} `json:"injury_status"`
	BirthCity             interface{} `json:"birth_city"`
	GsisID                interface{} `json:"gsis_id"`
	RotowireID            interface{} `json:"rotowire_id"`
	Weight                string      `json:"weight"`
	LastName              string      `json:"last_name"`
	Metadata              interface{} `json:"metadata"`
	PracticeParticipation interface{} `json:"practice_participation"`
	Height                string      `json:"height"`
	InjuryNotes           interface{} `json:"injury_notes"`
	BirthCountry          interface{} `json:"birth_country"`
	Number                int         `json:"number"`
	Age                   interface{} `json:"age"`
	SearchRank            int         `json:"search_rank"`
	RotoworldID           int         `json:"rotoworld_id"`
	HighSchool            interface{} `json:"high_school"`
	StatsID               interface{} `json:"stats_id"`
	YearsExp              int         `json:"years_exp"`
	Hashtag               string      `json:"hashtag"`
	InjuryStartDate       interface{} `json:"injury_start_date"`
	SearchLastName        string      `json:"search_last_name"`
	BirthState            interface{} `json:"birth_state"`
	FantasyDataID         int         `json:"fantasy_data_id"`
	PracticeDescription   interface{} `json:"practice_description"`
	InjuryBodyPart        interface{} `json:"injury_body_part"`
	SearchFullName        string      `json:"search_full_name"`
	SportradarID          string      `json:"sportradar_id"`
	Team                  interface{} `json:"team"`
	FirstName             string      `json:"first_name"`
	DepthChartPosition    interface{} `json:"depth_chart_position"`
	Active                bool        `json:"active"`
	PlayerID              string      `json:"player_id"`
	BirthDate             interface{} `json:"birth_date"`
	SearchFirstName       string      `json:"search_first_name"`
	YahooID               interface{} `json:"yahoo_id"`
	Sport                 string      `json:"sport"`
	DepthChartOrder       interface{} `json:"depth_chart_order"`
	PandascoreID          interface{} `json:"pandascore_id"`
	NewsUpdated           interface{} `json:"news_updated"`
}

type SleeperLeagueUser struct {
	UserID   string      `json:"user_id"`
	Settings interface{} `json:"settings"`
	Metadata struct {
		UserMessagePn           string `json:"user_message_pn"`
		TransactionWaiver       string `json:"transaction_waiver"`
		TransactionTrade        string `json:"transaction_trade"`
		TransactionFreeAgent    string `json:"transaction_free_agent"`
		TransactionCommissioner string `json:"transaction_commissioner"`
		TradeBlockPn            string `json:"trade_block_pn"`
		TeamNameUpdate          string `json:"team_name_update"`
		TeamName                string `json:"team_name"`
		ShowMascots             string `json:"show_mascots"`
		PlayerNicknameUpdate    string `json:"player_nickname_update"`
		PlayerLikePn            string `json:"player_like_pn"`
		MentionPn               string `json:"mention_pn"`
		MascotMessage           string `json:"mascot_message"`
		JoinVoicePn             string `json:"join_voice_pn"`
		Avatar                  string `json:"avatar"`
		AllowPn                 string `json:"allow_pn"`
	} `json:"metadata"`
	LeagueID    string      `json:"league_id"`
	IsOwner     bool        `json:"is_owner"`
	IsBot       bool        `json:"is_bot"`
	DisplayName string      `json:"display_name"`
	Avatar      interface{} `json:"avatar"`
}

type SleeperLeagueRoster struct {
	Taxi     interface{} `json:"taxi"`
	Starters []string    `json:"starters"`
	Settings struct {
		Wins               int `json:"wins"`
		WaiverPosition     int `json:"waiver_position"`
		WaiverBudgetUsed   int `json:"waiver_budget_used"`
		TotalMoves         int `json:"total_moves"`
		Ties               int `json:"ties"`
		PptsDecimal        int `json:"ppts_decimal"`
		Ppts               int `json:"ppts"`
		Losses             int `json:"losses"`
		FptsDecimal        int `json:"fpts_decimal"`
		FptsAgainstDecimal int `json:"fpts_against_decimal"`
		FptsAgainst        int `json:"fpts_against"`
		Fpts               int `json:"fpts"`
		Division           int `json:"division"`
	} `json:"settings"`
	RosterID  int         `json:"roster_id"`
	Reserve   interface{} `json:"reserve"`
	Players   []string    `json:"players"`
	PlayerMap interface{} `json:"player_map"`
	OwnerID   string      `json:"owner_id"`
	Metadata  struct {
		Streak string `json:"streak"`
		Record string `json:"record"`
	} `json:"metadata"`
	LeagueID string      `json:"league_id"`
	Keepers  interface{} `json:"keepers"`
	CoOwners interface{} `json:"co_owners"`
}

type SleeperPlayerFormatted struct {
	PlayerID string
	Name     string
	Points   float64
	Position string
}

type SleeperTeamFormatted struct {
	UserID    string
	RosterID  int
	Name      string
	Avatar    string
	Score     float64
	Starters  []SleeperPlayerFormatted
	Bench     []SleeperPlayerFormatted
	Wins      int
	Losses    int
	Ties      int
	MatchupID int
}

type SleeperLeagueFormatted struct {
	Name              string
	StartingPositions []string
}

// Fetch the league info given the league's id
func (d *Sleeper) GetLeague(id string) SleeperLeague {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/league/%s", id),
	}

	var responseData SleeperLeague
	_, err := d.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// Fetch the league users for a given league id
func (d *Sleeper) GetUsers(league_id string) []SleeperLeagueUser {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/league/%s/users", league_id),
	}

	var responseData []SleeperLeagueUser
	_, err := d.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// Fetch the league rosters for a given league id
func (d *Sleeper) GetRosters(league_id string) []SleeperLeagueRoster {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/league/%s/rosters", league_id),
	}

	var responseData []SleeperLeagueRoster
	_, err := d.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// Fetch the league matchups for a given week and league id
func (d *Sleeper) GetMatchups(league_id string, week string) []SleeperLeagueMatchup {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/league/%s/matchups/%s", league_id, week),
	}

	var responseData []SleeperLeagueMatchup
	_, err := d.Client.DoAndUnmarshal(request, &responseData)
	if err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// Get a player's info by id
// Fetches data from a cached file
func (d *Sleeper) GetPlayer(id string) SleeperPlayer {
	if sleperPlayers == nil {
		// Create the full path to the file you want to reference
		absPath := filepath.Join(util.Config.CacheDir, "./sleeper_players.json")

		if err := util.ReadFileAndUnmarshal(absPath, &sleperPlayers); err != nil {
			log.Printf("Error %v", err)
		}
	}

	player := sleperPlayers[id]

	return player

}

func (d *Sleeper) GetMatchupsFormatted(leagueID string, week string) [][]SleeperTeamFormatted {
	// get all users in league
	userIDToTeam := make(map[string]*SleeperTeamFormatted)

	for _, rawUser := range d.GetUsers(leagueID) {
		team := SleeperTeamFormatted{
			UserID: rawUser.UserID,
			Name:   rawUser.Metadata.TeamName,
			Avatar: rawUser.Metadata.Avatar,
		}
		// add to dict for easier lookup
		userIDToTeam[team.UserID] = &team
	}

	// get all rosters in a league
	rosterIDToTeam := make(map[int]*SleeperTeamFormatted)

	for _, rawRoster := range d.GetRosters(leagueID) {
		team := userIDToTeam[rawRoster.OwnerID]
		team.RosterID = rawRoster.RosterID
		team.Wins = rawRoster.Settings.Wins
		team.Losses = rawRoster.Settings.Losses
		team.Ties = rawRoster.Settings.Ties
		rosterIDToTeam[team.RosterID] = team
	}

	// get all matchup info for each roster
	matchupToTeams := make(map[int][]SleeperTeamFormatted)

	for _, rawMatchup := range d.GetMatchups(leagueID, week) {
		team := rosterIDToTeam[rawMatchup.RosterID]

		// add the current score total
		team.Score = rawMatchup.Points
		// add matchup id to team
		team.MatchupID = rawMatchup.MatchupID

		// add players and their scores
		for id, points := range rawMatchup.PlayersPoints {
			player := SleeperPlayerFormatted{
				PlayerID: id,
				Points:   points,
			}

			// fetch player details
			pInfo := d.GetPlayer(id)

			firstName := byte(0)
			if len(pInfo.FirstName) > 0 {
				firstName = pInfo.FirstName[0]
			}
			player.Name = fmt.Sprintf("%c.%s", firstName, pInfo.LastName)
			player.Position = pInfo.Position

			// determine if starter
			if slices.IndexFunc(rawMatchup.Starters, func(s string) bool { return s == id }) >= 0 {
				team.Starters = append(team.Starters, player)
			} else {
				team.Bench = append(team.Bench, player)
			}
		}
		// sort starters by position
		sort.Slice(team.Starters, func(i, j int) bool {
			l := slices.IndexFunc(rawMatchup.Starters, func(s string) bool { return s == team.Starters[i].PlayerID })
			r := slices.IndexFunc(rawMatchup.Starters, func(s string) bool { return s == team.Starters[j].PlayerID })
			return l < r
		})

		// add to last map
		matchupToTeams[rawMatchup.MatchupID] = append(matchupToTeams[rawMatchup.MatchupID], *team)
	}

	// convert matchup map to list of lists for easier rendering in view
	var matchups [][]SleeperTeamFormatted
	for _, teams := range matchupToTeams {
		matchups = append(matchups, teams)
	}
	// sort teams by matchup
	sort.Slice(matchups, func(i, j int) bool {
		return matchups[i][0].MatchupID < matchups[j][0].MatchupID
	})

	return matchups
}

func (d *Sleeper) GetLeagueFormatted(leagueID string) SleeperLeagueFormatted {
	league := d.GetLeague(leagueID)

	var startingPositions []string
	for _, pos := range league.RosterPositions {
		if pos != "BN" {
			startingPositions = append(startingPositions, pos[:1])
		}
	}

	return SleeperLeagueFormatted{
		Name:              league.Name,
		StartingPositions: startingPositions,
	}

}
