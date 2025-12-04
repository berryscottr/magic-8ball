package scrapaer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GraphQLClient struct {
	BaseURL string
	Headers map[string]string
}

func NewGraphQLClient(baseURL string, headers map[string]string) *GraphQLClient {
	return &GraphQLClient{BaseURL: baseURL, Headers: headers}
}

type GraphQLResponse struct {
	Data   json.RawMessage     `json:"data"`
	Errors []map[string]string `json:"errors,omitempty"`
}

func (c *GraphQLClient) Query(query string, variables map[string]interface{}, operationName string) (json.RawMessage, error) {
	jsonData := map[string]interface{}{
		"query":         query,
		"variables":     variables,
		"operationName": operationName,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		gqlResp.Errors = []map[string]string{{"message": string(body)}}
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL returned errors: %s", gqlResp.Errors[0]["message"])
	}

	return gqlResp.Data, err
}

type PoolPlayersAPI struct {
	Client *GraphQLClient
}

func NewPoolPlayersAPI(client *GraphQLClient) *PoolPlayersAPI {
	return &PoolPlayersAPI{Client: client}
}

func (api *PoolPlayersAPI) fetchData(query string, variables map[string]interface{}, operationName string, result interface{}) error {
	jsonData, err := api.Client.Query(query, variables, operationName)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, result)
}

func (api *PoolPlayersAPI) GetMemberStatsHeader(memberID int) (MemberStatsHeader, error) {
	var result MemberStatsHeader
	err := api.fetchData(MemberStatsHeaderQuery, map[string]interface{}{"id": memberID}, "getMemberStatsHeader", &result)
	return result, err
}

func (api *PoolPlayersAPI) GetPlayerTeams(playerID int) (PlayerTeams, error) {
	var result PlayerTeams
	err := api.fetchData(PlayerTeamsQuery, map[string]interface{}{
		"id":     playerID,
		"limit":  10000,
		"offset": 0,
	}, "TeamStat", &result)
	return result, err
}

func (api *PoolPlayersAPI) GetTeamRoster(teamID int) (TeamRoster, error) {
	var result TeamRoster
	err := api.fetchData(TeamRosterQuery, map[string]interface{}{"id": teamID}, "teamRoster", &result)
	return result, err
}

func (api *PoolPlayersAPI) GetTeamSchedule(teamID int) (TeamSchedule, error) {
	var result TeamSchedule
	err := api.fetchData(TeamScheduleQuery, map[string]interface{}{"id": teamID}, "teamSchedule", &result)
	return result, err
}

func (api *PoolPlayersAPI) GetMatch(matchID int) (MatchPage, error) {
	var result MatchPage
	err := api.fetchData(MatchQuery, map[string]interface{}{"id": matchID}, "MatchPage", &result)
	return result, err
}

func (api *PoolPlayersAPI) GetDivisionStandings(divisionID int) (DivisionStandings, error) {
	var result DivisionStandings
	err := api.fetchData(DivisionStandingsQuery, map[string]interface{}{"id": divisionID}, "DivisionStandings", &result)
	return result, err
}

func (api *PoolPlayersAPI) GetDivisionsDropdown(leagueID int, sessionID int) (DivisionsDropdown, error) {
	var result DivisionsDropdown
	err := api.fetchData(DivisionsDropdownQuery, map[string]interface{}{
		"id":      leagueID,
		"session": sessionID,
	}, "DivisionsDropdown", &result)
	return result, err
}
