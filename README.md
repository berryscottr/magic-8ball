# Magic 8Ball
APA Discord Bot for Wookie Mistakes Team

[![Go 1.18](https://img.shields.io/badge/golang-1.18-green.svg)](https://go.dev/dl/)
[![Python 3.10](https://img.shields.io/badge/python-3.10-blue.svg)](https://www.python.org/downloads/)

![Cat Pool](data/images/cat_pool.gif)

## Commands
### Lineups
`!line 23454567`

This command returns all eligible lineups for the first 8 given numbers in the message ranged between 2-7 to the #strategy channel.
### SL Matchups
`!heatsl`

This command returns the expected points of every skill level matchup as an image to the #strategy channel.

`!sl`

This command returns the expected points of every skill level matchup as text to the #strategy channel.
### Optimal Lineup
`!opt 65543 22235567`

This command returns the highest expected points lineups for that the second array (first 8) of numbers can 
respond to the first array (first 5) of numbers in the message. For this command to work, the arrays must be space separated.
### Game Day
`!game Wookie Mistakes`

This command returns a game day announcement to the #game-night channel that tracks incoming reactions for attendance.

## Dev Notes
- This bot is hard-coded to restart every 6th hour of the day in UTC time to enable full-time Github Workflow Action hosting.