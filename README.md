# Magic 8-Ball
APA Discord Bot for Wookie Mistakes Team

[![Go 1.18](https://img.shields.io/badge/golang-1.18-green.svg)](https://go.dev/dl/)
[![Python 3.7](https://img.shields.io/badge/python-3.7-blue.svg)](https://www.python.org/downloads/)
[![updateSkills](https://github.com/berryscottr/magic-8ball/actions/workflows/updateSkills.yml/badge.svg)](https://github.com/berryscottr/magic-8ball/actions/workflows/updateSkills.yml)
[![Build Magic 8Ball](https://github.com/berryscottr/magic-8ball/actions/workflows/build.yml/badge.svg?event=workflow_run)](https://github.com/berryscottr/magic-8ball/actions/workflows/build.yml)

![Cat Pool](data/images/cat_pool.gif)

## Commands
### Lineups
`!line 23454567`

This command returns all eligible lineups for the first 8 given numbers in the message ranged between 1-9 to the #strategy channel.
### SL Matchups
`!sl`

This command returns the expected points of every skill level matchup to the #strategy channel in both markdown text and links.

The embedded links have averages, medians, and modes for each skill level matchup.
and a link to the [heatmap](https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupAverages.svg).
### Lineups
`!inn`

This command returns everyone's effective innings per game towards their handicap to the #strategy channel.
### Optimal Lineup
`!opt 65543 22235567`

This command returns the highest expected points lineups for that the second array (first 8) of numbers can 
respond to the first array (first 5) of numbers in the message. For this command to work, the arrays must be space separated.
There is also about 15 seconds of delay until a response.
### Optimal Playoff Lineup
`!play 65543 22235567`

This command returns the highest differential expected points lineups for that the second array (first 8) of numbers can
respond to the first array (first 5) of numbers in the message. For this command to work, the arrays must be space separated.
There is also about 15 seconds of delay until a response.
### Game Day 8-Ball
`!8game Wookie Mistakes`

This command returns a game day announcement to the #game-night-8 channel that tracks incoming reactions for attendance.

### Game Day 9-Ball
`!9game Wookie Mistakes`

This command returns a game day announcement to the #game-night-8 channel that tracks incoming reactions for attendance.

## Dev Notes
- This bot is hard-coded to restart every 6th hour of the day in UTC time to enable full-time Github Workflow Action hosting.