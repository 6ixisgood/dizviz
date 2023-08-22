package data

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"path/filepath"
)

type SleeperConfig struct {
}


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
	TotalRosters int    `json:"total_rosters"`
	Status       string `json:"status"`
	Sport        string `json:"sport"`
	Shard        int    `json:"shard"`
	Settings     SleeperLeagueSettings `json:"settings"`
	SeasonType      string `json:"season_type"`
	Season          string `json:"season"`
	ScoringSettings SleeperLeagueScoringSettings `json:"scoring_settings"`
	RosterPositions  []string `json:"roster_positions"`
	PreviousLeagueID string   `json:"previous_league_id"`
	Name             string   `json:"name"`
	Metadata         SleeperLeagueMetadata `json:"metadata"`
	LoserBracketID        string		`json:"loser_bracket_id"`
	LeagueID              string      `json:"league_id"`
	LastReadID            string		 `json:"last_read_id"`
	LastPinnedMessageID   string      `json:"last_pinned_message_id"`
	LastMessageTime       int64       `json:"last_message_time"`
	LastMessageTextMap    interface{} `json:"last_message_text_map"`
	LastMessageID         string      `json:"last_message_id"`
	LastMessageAttachment interface{} `json:"last_message_attachment"`
	LastAuthorIsBot       bool        `json:"last_author_is_bot"`
	LastAuthorID          string      `json:"last_author_id"`
	LastAuthorDisplayName string      `json:"last_author_display_name"`
	LastAuthorAvatar      interface{} `json:"last_author_avatar"`
	GroupID               string		 `json:"group_id"`
	DraftID               string      `json:"draft_id"`
	CompanyID             string		 `json:"company_id"`
	BracketID             string	 `json:"bracket_id"`
	Avatar                interface{} `json:"avatar"`
}

type SleeperLeagueMatchup struct {
	Starters     []string    `json:"starters"`
	RosterID     int         `json:"roster_id"`
	Players      []string    `json:"players"`
	MatchupID    int         `json:"matchup_id"`
	Points       float64     `json:"points"`
	CustomPoints interface{} `json:"custom_points"`
}

type SleeperPlayer struct {
	Status                string      `json:"status"`
	EspnID                interface{} `json:"espn_id"`
	College               interface{} `json:"college"`
	SwishID               interface{} `json:"swish_id"`
	FantasyPositions      interface{} `json:"fantasy_positions"`
	Position              interface{} `json:"position"`
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



type Sleeper struct {
	Client			*APIClient
}

// Get a new client for Sleeper
func NewSleeper(baseUrl string) Sleeper {
	clientOptions := APIClientOptions {
		BaseURL: "https://api.sleeper.app/v1",
	}

	client := NewAPIClient(clientOptions)

	return Sleeper{
		Client: client,	
	}
}

// Fetch the league info given the league's id
func (d *Sleeper) GetLeague(id string) SleeperLeague {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/league/%s", id),
	}

	var responseData SleeperLeague
	if err := d.Client.DoAndUnmarshal(request, &responseData); err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// Fetch the league matchups for a given week and id
func (d *Sleeper) GetMatchups(id string, week string) []SleeperLeagueMatchup {
	request := &APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/league/%s/matchups/%s", id, week),
	}

	var responseData []SleeperLeagueMatchup
	if err := d.Client.DoAndUnmarshal(request, &responseData); err != nil {
		log.Fatalf("Failed to retrieve and unmarshal data. Error: %v", err)
	}

	return responseData
}

// Get a player's info by id
// Fetches data from a cached file 
func (d* Sleeper) GetPlayer(id string) SleeperPlayer {
	_, filePath, _, _ := runtime.Caller(0)
	// Get the directory of the current file
	basePath := filepath.Dir(filePath)
	// Create the full path to the file you want to reference
	absPath := filepath.Join(basePath, "./cache/sleeper_players.json")

	var responseData map[string]SleeperPlayer

	if err := ReadFileAndUnmarshal(absPath, &responseData); err != nil {
		log.Printf("Error %v", err)
	}

	player := responseData[id]

	return player

}
