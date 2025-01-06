package bot

import (
    "testing"
)


// Test function for illegalPlayback
func TestIllegalPlayback(t *testing.T) {
    // Define the team
    team := Team{Name: "Team A", Format: "8-Ball"}

    // Define test cases
    tests := []struct {
        name     string
        lineup   []int
        expected bool
    }{
        {
            name:     "Legal playback lineup",
            lineup:   []int{7, 5, 5, 3, 3}, // Legal because 2 or 3b could replace the last slot
            expected: false, // Playback is legal
        },
        {
            name:     "Illegal playback lineup",
            lineup:   []int{7, 7, 5, 2, 2}, // Illegal because no other combination works without playback
            expected: true, // Playback is illegal
        },
        {
            name:     "Another legal playback lineup",
            lineup:   []int{3, 5, 7, 5, 3}, // Legal because 2 could replace 3 or 5
            expected: false, // Playback is legal
        },
        {
            name:     "Legal lineup without a playback",
            lineup:   []int{2, 3, 5, 7, 2}, // No playback needed
            expected: false, // Playback is legal (no playback used)
        },
        {
            name:     "Lineup with more than one illegal playback",
            lineup:   []int{7, 7, 5, 5, 5}, // Illegal as 7, 7, 5, 5 cannot form a valid lineup without repetition
            expected: true, // Playback is illegal
        },
    }

    // Run the tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            illegal := illegalPlayback(tt.lineup, team)
            if illegal != tt.expected {
                t.Errorf("illegalPlayback(%v) = %v; want %v", tt.lineup, illegal, tt.expected)
            }
        })
    }
}
