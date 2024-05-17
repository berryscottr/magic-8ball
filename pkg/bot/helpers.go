package bot

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/stat/combin"
)

// seniorSkillRule returns a bool indicating if a lineup violates this rule
func seniorSkillRule(array []int) bool {
	numSeniors := 0
	for _, v := range array {
		if v >= SeniorSkillLevel {
			numSeniors++
		}
	}
	return numSeniors > 2
}

// sum returns the sum of the elements in the given int slice
func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

// containsSlice returns true if the given slice of int slices contains the given int slice
func containsSlice(slc [][]int, ele []int) bool {
	for _, v := range slc {
		if reflect.DeepEqual(v, ele) {
			return true
		}
	}
	return false
}

// validLineup returns true if the given lineup is valid
func validLineup(lineup []int, isPlayback bool, team Team) bool {
	if len(lineup) != 5 {
		return false
	}
	if team.Format == "8-Ball" {
		if sum(lineup) < 10 || sum(lineup) > 23 {
			return false
		}
	} else {
		if sum(lineup) < 5 || sum(lineup) > 23 {
			return false
		}
	}
	if seniorSkillRule(lineup) {
		return false
	}
	// if isPlayback {
		// if illegalPlayback(lineup, team) {
		// 	return false
		// }
	// }
	return true
}

func lineupsMsgLogic(availablePlayerSkills []int, team Team) string {
	var lineupsMsg string
	var noPlaybacksAllowed bool
	var isPlayback bool
	if len(availablePlayerSkills) < 5 {
		if len(availablePlayerSkills) == 4 && !noPlaybacksAllowed {
			isPlayback = true
			teamLineups := generateLineups(availablePlayerSkills, isPlayback, team)
			lineupsMsg = generateLineupsMsg(teamLineups, isPlayback)
		} else {
			lineupsMsg = "None found"
		}
	} else if len(availablePlayerSkills) >= 5 {
		teamLineups := generateLineups(availablePlayerSkills, isPlayback, team)
		if len(teamLineups) == 0 && !noPlaybacksAllowed {
			isPlayback = true
			teamLineups := generateLineups(availablePlayerSkills, isPlayback, team)
			lineupsMsg = generateLineupsMsg(teamLineups, isPlayback)
		} else {
			lineupsMsg = generateLineupsMsg(teamLineups, isPlayback)
		}
	}
	return lineupsMsg
}

// generateLineups returns a slice of TeamLineup structs
func generateLineups(skills []int, isPlayback bool, team Team) []TeamLineup {
	var lineups [][]int
	var teamLineups []TeamLineup
	generate := func(playerSkills []int) {
		sort.Sort(sort.Reverse(sort.IntSlice(playerSkills)))
		log.Info().Msgf("generating possible lineups of %v", playerSkills)
		gen := combin.NewPermutationGenerator(len(playerSkills), 5)
		for gen.Next() {
			permutationIndices := gen.Permutation(nil)
			var lineup []int
			for _, permutationIndex := range permutationIndices {
				lineup = append(lineup, playerSkills[permutationIndex])
			}
			sort.Sort(sort.Reverse(sort.IntSlice(lineup)))
			if validLineup(lineup, isPlayback, team) && !containsSlice(lineups, lineup) {
				lineups = append(lineups, lineup)
			}
		}
		for _, lineup := range lineups {
			teamLineups = append(teamLineups, TeamLineup{Lineup: lineup, Sum: sum(lineup)})
		}
		sort.Slice(teamLineups[:], func(i, j int) bool {
			return teamLineups[i].Sum > teamLineups[j].Sum
		})
	}
	if isPlayback {
		for i := range len(skills) {
			playbackPlayerSkills := append(skills, skills[i])
			generate(playbackPlayerSkills)
		}
	} else {
		generate(skills)
	}
	return teamLineups
}

// generateLineupsMsg returns a string of lineups and a string of playbackNote
func generateLineupsMsg(teamLineups []TeamLineup, isPlayback bool) string {
	var lineupsMsg string
	if len(teamLineups) == 0 {
		lineupsMsg = "None found"
		return lineupsMsg
	}
	for _, teamLineup := range teamLineups {
		lineupsMsg += fmt.Sprintf("%v %v\n", teamLineup.Lineup, teamLineup.Sum)
	}
	if isPlayback {
		lineupsMsg = "```" + lineupsMsg + "```" + "(eligible with playback)"
	} else {
		lineupsMsg = "```" + lineupsMsg + "```"
	}
	return lineupsMsg
}

// Check if playback is legal
// func illegalPlayback(lineup []int, team Team) bool {
// 	var illegalPlayback bool
// 	var teammates []TeammateSkill
// 	for _, teammate := range Teammates {
// 		for _, teammateTeam := range teammate.Teams {
// 			if teammateTeam.Name == team.Name {
// 				if teammateTeam.Format == "8-Ball" {
// 					teammates = append(teammates, TeammateSkill{Lastname: teammate.LastName, SkillLevel: teammate.SkillLevel.Eight})
// 				} else if teammateTeam.Format == "9-Ball" {
// 					teammates = append(teammates, TeammateSkill{Lastname: teammate.LastName, SkillLevel: teammate.SkillLevel.Nine})
// 				}
// 				break
// 			}
// 		}
// 	}
// 	var allPlayerSkills []int
// 	for _, teammate := range teammates {
// 		allPlayerSkills = append(allPlayerSkills, teammate.SkillLevel)
// 	}
// 	allLegalTeamLineups := generateLineups(allPlayerSkills, false, team)
// 	var allLegalLineups [][]int
// 	for _, legalTeamLineup := range allLegalTeamLineups {
// 		allLegalLineups = append(allLegalLineups, legalTeamLineup.Lineup)
// 	}
// 	if !containsSlice(allLegalLineups, lineup) {
// 		illegalPlayback = true
// 	}
// 	return illegalPlayback
// }
