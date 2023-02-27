package bot

import (
	"reflect"
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
func validLineup(lineup []int) bool {
	if len(lineup) != 5 {
		return false
	}
	if sum(lineup) < 10 || sum(lineup) > 23 {
		return false
	}
	if seniorSkillRule(lineup) {
		return false
	}
	return true
}
