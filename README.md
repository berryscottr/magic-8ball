# Magic 8-Ball

APA Discord Bot for Wookie Mistakes and Safety Dance Teams

[![Go 1.20](https://img.shields.io/badge/golang-1.20-green.svg)](https://go.dev/dl/)
[![Python 3.11](https://img.shields.io/badge/python-3.11-blue.svg)](https://www.python.org/downloads/)
[![updateSkills](https://github.com/berryscottr/magic-8ball/actions/workflows/updateSkills.yml/badge.svg)](https://github.com/berryscottr/magic-8ball/actions/workflows/updateSkills.yml)
[![Build Magic 8Ball](https://github.com/berryscottr/magic-8ball/actions/workflows/build.yml/badge.svg?event=workflow_run)](https://github.com/berryscottr/magic-8ball/actions/workflows/build.yml)

![Cat Pool](data/images/cat_pool.gif)

## Commands

### Lineups

`!line 23454567`

This command returns all eligible lineups for the first 8 given numbers in the message ranged between 1-9 to the #strategy channel.

### SL Matchups

`!sl`

This command returns the expected points of every skill level matchup for 8-ball and 9-ball to the #strategy channel in both markdown text and links.

The embedded links have averages, medians, and modes for each skill level matchup.
and a link to the [8-ball heatmap](https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupAverages.svg) and [9-ball heatmap](https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupAveragesNine.svg).

### Inning Counts

`!inn`

This command returns everyone's effective innings per game towards their handicap to the #strategy channel.

### Optimal Lineup Eight-Ball

`!opt8 65543 23445567`

This command returns the highest expected points lineups for that the second array (first 8) of numbers can
respond to the first array (first 5) of numbers in the message. For this command to work, the arrays must be space separated.
There is also about 15 seconds of delay until a response.

### Optimal Lineup Nine-Ball

`!opt9 65543 12346789`

This command returns the highest expected points lineups for that the second array (first 8) of numbers can
respond to the first array (first 5) of numbers in the message. For this command to work, the arrays must be space separated.
There is also about 15 seconds of delay until a response.

### Optimal Playoff Lineup Eight-Ball

`!play 65543 22235567`

This command returns the highest differential expected points lineups for that the second array (first 8) of numbers can
respond to the first array (first 5) of numbers in the message. For this command to work, the arrays must be space separated.
There is also about 15 seconds of delay until a response.

### Game Day 8-Ball

`!8game Wookie Mistakes "custom message"`

This command returns a game day announcement to the #game-night-8 channel that tracks incoming reactions for attendance. If composed prior to Tuesday at 5am, it will schedule to be posted at 8:55am on Tuesday. If composed after Tuesday at 5am, it will post immediately.

A message can always be posted immediately by adding the `--now` keyword to the end of the command. This is useful for posting a game day announcement during the playoffs.

### Game Day 9-Ball

`!9game Safety Dance "custom message"`

This command returns a game day announcement to the #game-night-9 channel that tracks incoming reactions for attendance. If composed prior to Tuesday at 5am, it will schedule to be posted at 8:55am on Tuesday. If composed after Tuesday at 5am, it will post immediately.

A message can always be posted immediately by adding the `--now` keyword to the end of the command. This is useful for posting a game day announcement during the playoffs.

### Calendar

`!cal`

This command returns an embedded link to the calendar for the current season.

Further, the current session's team calendar with match results is located [here](data/schedules/Spring2023Schedule.csv).

## Dev Notes

- This bot is hosted on Railway App
- I've used this bot to benefit my own teams, but I encourage whoever may come across this to use this template for their own teams as well.
