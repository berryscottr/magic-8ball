package scrapaer

import (
  "time"
)

var (
  // Location assumed table size
  TableSizes = map[int]int{
    87357: 7, // No Match This Week
    14771: 9, // Cues Billiards // has some 7 footers
    39855: 8, // Mr Cues II
    14773: 8, // Mr Cues II
    39856: 8, // Mr Cues II
    63333: 9, // The Independent
    83586: 7, // Amsterdam Cafe
    27131: 7, // Smith's Olde Bar
    14792: 7, // Burkhart's Pub
    93787: 7, // Friends
    73919: 7, // BJ Rooster
    65643: 7, // The Hideaway
    53393: 7, // Le Buzz 
    41205: 9, // Big Shots Billiards
    109808: 7, // Owlz Pub
    66994: 7, // Mazzy's Kennesaw
    14754: 7, // Kennesaw Billiards
    39958: 7, // Sidelines - Kennesaw
    110059: 7, // Varner's
    117721: 7, // The Pool Turtle
  }
	// GraphQL queries
  DivisionStandingsQuery = `
query divsionStandings($id: Int!) {
  division(id: $id) {
    id
    teams {
      id
      name
      number
      standing
      pointsLastWeek
      lastWeek
      sessionTotalPoints
      totalTeamMatchesPlayed
      isTied
      isBye
      league {
        id
        slug
        __typename
      }
      tournaments(time: FUTURE) {
        id
        tournamentTypeName
        __typename
      }
      __typename
    }
    __typename
  }
}`
  DivisionsDropdownQuery = `
query divisionsDropdown($id: Int, $session: Int) {
  league(id: $id) {
    id
    slug
    divisions(session: $session) {
      id
      name
      number
      __typename
    }
    __typename
  }
}`
	MemberStatsHeaderQuery = `
query getMemberStatsHeader($id: Int!) {
  member(id: $id) {
    id
    firstName
    lastName
    initials
    aliases(includeInactive: true) {
      id
      displayName
      memberNumber
      league {
        id
        slug
        name
        isDefault
        number
        __typename
      }
      formats
      stats {
        __typename
        skillLevel
        CLA
      }
      __typename
    }
    __typename
  }
}`
	PlayerTeamsQuery = `
query TeamStat($id: Int!, $limit: Int!, $offset: Int!) {
  alias(id: $id) {
    id
    pastTeams: players(current: false, active: null, limit: $limit, offset: $offset) {
      id
      ...EightBallTeam
      ...NineBallTeam
      ...MastersTeam
      __typename
    }
    currentTeams: players(current: true, active: null) {
      id
      ...EightBallTeam
      ...NineBallTeam
      ...MastersTeam
      __typename
    }
    __typename
  }
}

fragment NineBallTeam on NineBallPlayer {
  id
  isActive
  role
  rosterPosition
  nickName
  matchesPlayed
  matchesWon
  session {
    id
    name
    __typename
  }
  skillLevel
  rank
  team {
    id
    name
    division {
      id
      isTournament
      __typename
    }
    __typename
  }
  __typename
}

fragment EightBallTeam on EightBallPlayer {
  id
  isActive
  role
  rosterPosition
  nickName
  matchesPlayed
  matchesWon
  session {
    id
    name
    __typename
  }
  skillLevel
  rank
  team {
    id
    name
    division {
      id
      isTournament
      __typename
    }
    __typename
  }
  __typename
}

fragment MastersTeam on MastersPlayer {
  id
  isActive
  role
  rosterPosition
  nickName
  matchesPlayed
  matchesWon
  session {
    id
    name
    __typename
  }
  team {
    id
    name
    division {
      id
      isTournament
      __typename
    }
    __typename
  }
  __typename
}`
	TeamRosterQuery = `
query teamRoster($id: Int!) {
  team(id: $id) {
    ...rosterComponent
    __typename
  }
}

fragment rosterComponent on Team {
  id
  name
  number
  league {
    id
    slug
    __typename
  }
  division {
    id
    type
    __typename
  }
  roster {
    id
    memberNumber
    displayName
    matchesWon
    matchesPlayed
    ... on EightBallPlayer {
      pa
      ppm
      skillLevel
      __typename
    }
    ... on NineBallPlayer {
      pa
      ppm
      skillLevel
      __typename
    }
    member {
      id
      __typename
    }
    __typename
  }
  __typename
}`
	TeamScheduleQuery = `
query teamSchedule($id: Int!) {
  team(id: $id) {
    id
    sessionBonusPoints
    sessionPoints
    sessionTotalPoints
    division {
      id
      isTournament
      __typename
    }
    matches(unscheduled: true) {
      week
      type
      ...matchListItem
      __typename
    }
    __typename
  }
}

fragment matchListItem on Match {
  id
  isBye
  status
  scoresheet
  startTime
  isMine
  isPaid
  isPlayoff
  description
  results {
    homeAway
    points {
      total
      __typename
    }
    __typename
  }
  location {
    id
    name
    address {
      id
      name
      __typename
    }
    __typename
  }
  home {
    id
    name
    number
    isMine
    __typename
  }
  away {
    id
    name
    number
    isMine
    __typename
  }
  league {
    id
    isMine
    slug
    isElectronicPaymentsEnabled
    __typename
  }
  orderItems {
    id
    order {
      id
      member {
        id
        firstName
        lastName
        __typename
      }
      __typename
    }
    __typename
  }
  division {
    id
    scheduleInEdit
    isTournament
    __typename
  }
  __typename
}`
	MatchQuery = `
query MatchPage($id: Int!) {
  match(id: $id) {
    id
    division {
      id
      electronicScoringEnabled
      __typename
    }
    league {
      id
      esEnabled
      __typename
    }
    ...matchForCart
    __typename
  }
}

fragment matchForCart on Match {
  __typename
  id
  type
  startTime
  week
  isBye
  isMine
  isScored
  scoresheet
  isPaid
  location {
    ...googleMapComponent
    __typename
  }
  home {
    id
    name
    number
    isMine
    ...rosterComponent
    __typename
  }
  away {
    id
    name
    number
    isMine
    ...rosterComponent
    __typename
  }
  division {
    id
    scheduleInEdit
    type
    __typename
  }
  session {
    id
    name
    year
    __typename
  }
  league {
    id
    name
    currentSessionId
    isElectronicPaymentsEnabled
    country {
      id
      __typename
    }
    __typename
  }
  fees {
    amount
    tax
    total
    __typename
  }
  orderItems {
    id
    order {
      id
      member {
        id
        firstName
        lastName
        __typename
      }
      __typename
    }
    __typename
  }
  results {
    homeAway
    overUnder
    forfeits
    matchesWon
    matchesPlayed
    points {
      bonus
      penalty
      won
      adjustment
      sportsmanship
      total
      skillLevelViolationAdjustment
      __typename
    }
    scores {
      id
      player {
        id
        displayName
        __typename
      }
      matchPositionNumber
      playerPosition
      skillLevel
      innings
      defensiveShots
      eightBallWins
      eightOnBreak
      eightBallBreakAndRun
      nineBallPoints
      nineOnSnap
      nineBallBreakAndRun
      nineBallMatchPointsEarned
      mastersEightBallWins
      mastersNineBallWins
      winLoss
      matchForfeited
      doublesMatch
      dateTimeStamp
      teamSlot
      eightBallMatchPointsEarned
      incompleteMatch
      __typename
    }
    __typename
  }
}

fragment googleMapComponent on HostLocation {
  id
  phone
  name
  address {
    id
    name
    address1
    address2
    city
    zip
    latitude
    longitude
    __typename
  }
  __typename
}

fragment rosterComponent on Team {
  id
  name
  number
  league {
    id
    slug
    __typename
  }
  division {
    id
    type
    __typename
  }
  roster {
    id
    memberNumber
    displayName
    matchesWon
    matchesPlayed
    ... on EightBallPlayer {
      pa
      ppm
      skillLevel
      __typename
    }
    ... on NineBallPlayer {
      pa
      ppm
      skillLevel
      __typename
    }
    member {
      id
      __typename
    }
    __typename
  }
  __typename
}`
)

type DivisionStandings struct {
	Division struct {
		ID    int `json:"id"`
		Teams []struct {
			ID                     int    `json:"id"`
			Name                   string `json:"name"`
			Number                 string `json:"number"`
			Standing               int    `json:"standing"`
			PointsLastWeek         int    `json:"pointsLastWeek"`
			LastWeek               int    `json:"lastWeek"`
			SessionTotalPoints     int    `json:"sessionTotalPoints"`
			TotalTeamMatchesPlayed int    `json:"totalTeamMatchesPlayed"`
			IsTied                 bool   `json:"isTied"`
			IsBye                  bool   `json:"isBye"`
			League                 struct {
				ID       int    `json:"id"`
				Slug     string `json:"slug"`
				Typename string `json:"__typename"`
			} `json:"league"`
			Tournaments []interface{} `json:"tournaments"`
			Typename    string        `json:"__typename"`
		} `json:"teams"`
		Typename string `json:"__typename"`
	} `json:"division"`
}

type DivisionsDropdown struct {
	League struct {
		ID        int    `json:"id"`
		Slug      string `json:"slug"`
		Divisions []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Number   string `json:"number"`
			Typename string `json:"__typename"`
		} `json:"divisions"`
		Typename string `json:"__typename"`
	} `json:"league"`
}

type PlayerMatch struct {
	ID     int `json:"id"`
	Player struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
		Typename    string `json:"__typename"`
	} `json:"player"`
  TableSize                  int         `json:"tableSize"`
	MatchPositionNumber        int         `json:"matchPositionNumber"`
	PlayerPosition             int         `json:"playerPosition"`
	SkillLevel                 int         `json:"skillLevel"`
  OpponentSkillLevel         int         `json:"opponentSkillLevel"`
	Innings                    int         `json:"innings"`
	DefensiveShots             int         `json:"defensiveShots"`
  OpponentDefensiveShots     int         `json:"opponentDefensiveShots"`
	EightBallWins              interface{} `json:"eightBallWins"`
	EightBallLosses            interface{} `json:"eightBallLosses"`
	EightOnBreak               interface{} `json:"eightOnBreak"`
	EightBallBreakAndRun       interface{} `json:"eightBallBreakAndRun"`
	NineBallPoints             interface{} `json:"nineBallPoints"`
	NineOnSnap                 interface{} `json:"nineOnSnap"`
	NineBallBreakAndRun        interface{} `json:"nineBallBreakAndRun"`
	NineBallMatchPointsEarned  interface{} `json:"nineBallMatchPointsEarned"`
	MastersEightBallWins       interface{} `json:"mastersEightBallWins"`
	MastersNineBallWins        interface{} `json:"mastersNineBallWins"`
	WinLoss                    string      `json:"winLoss"`
	MatchForfeited             bool        `json:"matchForfeited"`
	DoublesMatch               bool        `json:"doublesMatch"`
	DateTimeStamp              time.Time   `json:"dateTimeStamp"`
	TeamSlot                   string      `json:"teamSlot"`
	EightBallMatchPointsEarned interface{} `json:"eightBallMatchPointsEarned"`
  EightBallMatchPointsLost   interface{} `json:"eightBallMatchPointsLost"`
	IncompleteMatch            bool        `json:"incompleteMatch"`
	Typename                   string      `json:"__typename"`
}

type Player struct {
	MemberID int
	PlayerID int
  Name    string
	TeamIDs  []int
	Matches  []PlayerMatch
}

type MemberStatsHeader struct {
	Member struct {
		ID        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Initials  string `json:"initials"`
		Aliases   []struct {
			ID           int    `json:"id"`
			DisplayName  string `json:"displayName"`
			MemberNumber string `json:"memberNumber"`
			League       struct {
				ID        int    `json:"id"`
				Slug      string `json:"slug"`
				Name      string `json:"name"`
				IsDefault bool   `json:"isDefault"`
				Number    string `json:"number"`
				Typename  string `json:"__typename"`
			} `json:"league"`
			Formats []string `json:"formats"`
			Stats   []struct {
				Typename   string `json:"__typename"`
				SkillLevel int    `json:"skillLevel"`
				CLA        int    `json:"CLA"`
			} `json:"stats"`
			Typename string `json:"__typename"`
		} `json:"aliases"`
		Typename string `json:"__typename"`
	} `json:"member"`
}

type PlayerTeams struct {
	Alias struct {
		ID        int `json:"id"`
		PastTeams []struct {
			ID             int         `json:"id"`
			IsActive       bool        `json:"isActive"`
			Role           string      `json:"role"`
			RosterPosition int         `json:"rosterPosition"`
			NickName       interface{} `json:"nickName"`
			MatchesPlayed  int         `json:"matchesPlayed"`
			MatchesWon     int         `json:"matchesWon"`
			Session        struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Typename string `json:"__typename"`
			} `json:"session"`
			SkillLevel int `json:"skillLevel"`
			Rank       int `json:"rank"`
			Team       struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Division struct {
					ID           int    `json:"id"`
					IsTournament bool   `json:"isTournament"`
					Typename     string `json:"__typename"`
				} `json:"division"`
				Typename string `json:"__typename"`
			} `json:"team"`
			Typename string `json:"__typename"`
		} `json:"pastTeams"`
		CurrentTeams []struct {
			ID             int         `json:"id"`
			IsActive       bool        `json:"isActive"`
			Role           string      `json:"role"`
			RosterPosition int         `json:"rosterPosition"`
			NickName       interface{} `json:"nickName"`
			MatchesPlayed  int         `json:"matchesPlayed"`
			MatchesWon     int         `json:"matchesWon"`
			Session        struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Typename string `json:"__typename"`
			} `json:"session"`
			SkillLevel int `json:"skillLevel"`
			Rank       int `json:"rank"`
			Team       struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Division struct {
					ID           int    `json:"id"`
					IsTournament bool   `json:"isTournament"`
					Typename     string `json:"__typename"`
				} `json:"division"`
				Typename string `json:"__typename"`
			} `json:"team"`
			Typename string `json:"__typename"`
		} `json:"currentTeams"`
		Typename string `json:"__typename"`
	} `json:"alias"`
}

type TeamRoster struct {
	Team struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Number string `json:"number"`
		League struct {
			ID       int    `json:"id"`
			Slug     string `json:"slug"`
			Typename string `json:"__typename"`
		} `json:"league"`
		Division struct {
			ID       int    `json:"id"`
			Type     string `json:"type"`
			Typename string `json:"__typename"`
		} `json:"division"`
		Roster []struct {
			ID            int     `json:"id"`
			MemberNumber  string  `json:"memberNumber"`
			DisplayName   string  `json:"displayName"`
			MatchesWon    int     `json:"matchesWon"`
			MatchesPlayed int     `json:"matchesPlayed"`
			Pa            float64 `json:"pa"`
			Ppm           float64 `json:"ppm"`
			SkillLevel    int     `json:"skillLevel"`
			Typename      string  `json:"__typename"`
			Member        struct {
				ID       int    `json:"id"`
				Typename string `json:"__typename"`
			} `json:"member"`
		} `json:"roster"`
		Typename string `json:"__typename"`
	} `json:"team"`
}

type TeamSchedule struct {
	Team struct {
		ID                 int `json:"id"`
		SessionBonusPoints int `json:"sessionBonusPoints"`
		SessionPoints      int `json:"sessionPoints"`
		SessionTotalPoints int `json:"sessionTotalPoints"`
		Division           struct {
			ID           int    `json:"id"`
			IsTournament bool   `json:"isTournament"`
			Typename     string `json:"__typename"`
		} `json:"division"`
		Matches []struct {
			Week        interface{}   `json:"week"`
			Type        interface{}   `json:"type"`
			ID          interface{}   `json:"id"`
			IsBye       bool          `json:"isBye"`
			Status      string        `json:"status"`
			Scoresheet  interface{}   `json:"scoresheet"`
			StartTime   time.Time     `json:"startTime"`
			IsMine      bool          `json:"isMine"`
			IsPaid      interface{}   `json:"isPaid"`
			IsPlayoff   bool          `json:"isPlayoff"`
			Description string        `json:"description"`
			Results     []interface{} `json:"results"`
			Location    interface{}   `json:"location"`
			Home        interface{}   `json:"home"`
			Away        interface{}   `json:"away"`
			League      interface{}   `json:"league"`
			OrderItems  interface{}   `json:"orderItems"`
			Division    struct {
				ID             int    `json:"id"`
				ScheduleInEdit bool   `json:"scheduleInEdit"`
				IsTournament   bool   `json:"isTournament"`
				Typename       string `json:"__typename"`
			} `json:"division"`
			Typename string `json:"__typename"`
		} `json:"matches"`
		Typename string `json:"__typename"`
	} `json:"team"`
}

type MatchResult struct {
  HomeAway      string `json:"homeAway"`
  // OverUnder     int    `json:"overUnder"`
  Forfeits      int    `json:"forfeits"`
  MatchesWon    int    `json:"matchesWon"`
  MatchesPlayed int    `json:"matchesPlayed"`
  Points        struct {
    Bonus                         int    `json:"bonus"`
    Penalty                       int    `json:"penalty"`
    Won                           int    `json:"won"`
    Adjustment                    int    `json:"adjustment"`
    Sportsmanship                 int    `json:"sportsmanship"`
    Total                         int    `json:"total"`
    SkillLevelViolationAdjustment int    `json:"skillLevelViolationAdjustment"`
    Typename                      string `json:"__typename"`
  } `json:"points"`
  Scores []struct {
    ID     int `json:"id"`
    Player struct {
      ID          int    `json:"id"`
      DisplayName string `json:"displayName"`
      Typename    string `json:"__typename"`
    } `json:"player"`
    TableSize                  int         `json:"tableSize"`
    MatchPositionNumber        int         `json:"matchPositionNumber"`
    PlayerPosition             int         `json:"playerPosition"`
    SkillLevel                 int         `json:"skillLevel"`
    OpponentSkillLevel         int         `json:"opponentSkillLevel"`
    Innings                    int         `json:"innings"`
    DefensiveShots             int         `json:"defensiveShots"`
    OpponentDefensiveShots     int         `json:"opponentDefensiveShots"`
    EightBallWins              interface{} `json:"eightBallWins"`
    EightBallLosses            interface{} `json:"eightBallLosses"`
    EightOnBreak               interface{} `json:"eightOnBreak"`
    EightBallBreakAndRun       interface{} `json:"eightBallBreakAndRun"`
    NineBallPoints             interface{} `json:"nineBallPoints"`
    NineOnSnap                 interface{} `json:"nineOnSnap"`
    NineBallBreakAndRun        interface{} `json:"nineBallBreakAndRun"`
    NineBallMatchPointsEarned  interface{} `json:"nineBallMatchPointsEarned"`
    MastersEightBallWins       interface{} `json:"mastersEightBallWins"`
    MastersNineBallWins        interface{} `json:"mastersNineBallWins"`
    WinLoss                    string      `json:"winLoss"`
    MatchForfeited             bool        `json:"matchForfeited"`
    DoublesMatch               bool        `json:"doublesMatch"`
    DateTimeStamp              time.Time   `json:"dateTimeStamp"`
    TeamSlot                   string      `json:"teamSlot"`
    EightBallMatchPointsEarned interface{} `json:"eightBallMatchPointsEarned"`
    EightBallMatchPointsLost   interface{} `json:"eightBallMatchPointsLost"`
    IncompleteMatch            bool        `json:"incompleteMatch"`
    Typename                   string      `json:"__typename"`
  } `json:"scores"`
  Typename string `json:"__typename"`
}

type MatchPage struct {
	Match struct {
		ID       int `json:"id"`
		Division struct {
			ID                       int    `json:"id"`
			ElectronicScoringEnabled bool   `json:"electronicScoringEnabled"`
			Typename                 string `json:"__typename"`
			ScheduleInEdit           bool   `json:"scheduleInEdit"`
			Type                     string `json:"type"`
		} `json:"division"`
		League struct {
			ID                          int    `json:"id"`
			EsEnabled                   bool   `json:"esEnabled"`
			Typename                    string `json:"__typename"`
			Name                        string `json:"name"`
			CurrentSessionID            int    `json:"currentSessionId"`
			IsElectronicPaymentsEnabled bool   `json:"isElectronicPaymentsEnabled"`
			Country                     struct {
				ID       int    `json:"id"`
				Typename string `json:"__typename"`
			} `json:"country"`
		} `json:"league"`
		Typename   string      `json:"__typename"`
		Type       string      `json:"type"`
		StartTime  string      `json:"startTime"`
		Week       int         `json:"week"`
		IsBye      bool        `json:"isBye"`
		IsMine     bool        `json:"isMine"`
		IsScored   bool        `json:"isScored"`
		Scoresheet interface{} `json:"scoresheet"`
		IsPaid     bool        `json:"isPaid"`
		Location   struct {
			ID      int    `json:"id"`
			Phone   string `json:"phone"`
			Name    string `json:"name"`
			Address struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Address1  string  `json:"address1"`
				Address2  string  `json:"address2"`
				City      string  `json:"city"`
				Zip       string  `json:"zip"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Typename  string  `json:"__typename"`
			} `json:"address"`
			Typename string `json:"__typename"`
		} `json:"location"`
		Home struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Number string `json:"number"`
			IsMine bool   `json:"isMine"`
			League struct {
				ID       int    `json:"id"`
				Slug     string `json:"slug"`
				Typename string `json:"__typename"`
			} `json:"league"`
			Division struct {
				ID       int    `json:"id"`
				Type     string `json:"type"`
				Typename string `json:"__typename"`
			} `json:"division"`
			Roster []struct {
				ID            int     `json:"id"`
				MemberNumber  string  `json:"memberNumber"`
				DisplayName   string  `json:"displayName"`
				MatchesWon    int     `json:"matchesWon"`
				MatchesPlayed int     `json:"matchesPlayed"`
				Pa            float64 `json:"pa"`
				Ppm           float64 `json:"ppm"`
				SkillLevel    int     `json:"skillLevel"`
				Typename      string  `json:"__typename"`
				Member        struct {
					ID       int    `json:"id"`
					Typename string `json:"__typename"`
				} `json:"member"`
			} `json:"roster"`
			Typename string `json:"__typename"`
		} `json:"home"`
		Away struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Number string `json:"number"`
			IsMine bool   `json:"isMine"`
			League struct {
				ID       int    `json:"id"`
				Slug     string `json:"slug"`
				Typename string `json:"__typename"`
			} `json:"league"`
			Division struct {
				ID       int    `json:"id"`
				Type     string `json:"type"`
				Typename string `json:"__typename"`
			} `json:"division"`
			Roster []struct {
				ID            int     `json:"id"`
				MemberNumber  string  `json:"memberNumber"`
				DisplayName   string  `json:"displayName"`
				MatchesWon    int     `json:"matchesWon"`
				MatchesPlayed int     `json:"matchesPlayed"`
				Pa            float64 `json:"pa"`
				Ppm           float64 `json:"ppm"`
				SkillLevel    int     `json:"skillLevel"`
				Typename      string  `json:"__typename"`
				Member        struct {
					ID       int    `json:"id"`
					Typename string `json:"__typename"`
				} `json:"member"`
			} `json:"roster"`
			Typename string `json:"__typename"`
		} `json:"away"`
		Session struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Year     int    `json:"year"`
			Typename string `json:"__typename"`
		} `json:"session"`
		Fees struct {
			Amount   int    `json:"amount"`
			Tax      int    `json:"tax"`
			Total    int    `json:"total"`
			Typename string `json:"__typename"`
		} `json:"fees"`
		OrderItems []interface{} `json:"orderItems"`
		Results    []MatchResult `json:"results"`
	} `json:"match"`
}
