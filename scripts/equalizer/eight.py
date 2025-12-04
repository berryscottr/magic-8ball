import numpy as np
import pandas as pd
import math

class EightBallPlayer:
    expected_innings_eight = {
        7: (1.00, 2.00),
        6: (2.00, 3.00),
        5: (3.00, 4.00),
        4: (4.00, 5.00),
        3: (5.00, 6.00),
        2: (7.00, 7.00)
    }
    skill_levels_eight = {
        7: (0.00, 2.00),
        6: (2.01, 3.00),
        5: (3.01, 4.00),
        4: (4.01, 5.00),
        3: (5.01, 7.00),
        2: (7.01, 15.00)
    }
    
    expected_eight_scaler = 1.0

    def __init__(self, df):
        self.df_eight_ball = df[df['player___typename'] == 'EightBallPlayer'][['skillLevel', 'innings', 'defensiveShots', 'eightBallWins', 'eightBallLosses', 'winLoss', 'opponentSkillLevel', 'tableSize']]
        self.df_eight_ball['eightBallWins'] = pd.to_numeric(self.df_eight_ball['eightBallWins'], errors='coerce').fillna(0).astype(int)
        self.df_eight_ball['eightBallLosses'] = pd.to_numeric(self.df_eight_ball['eightBallLosses'], errors='coerce').fillna(0).astype(int)
        self.process()

    @staticmethod
    def apa_8_ball_race(skill1, skill2):
        if skill1 == 0:
            skill1 = 3
        if skill2 == 0:
            skill2 = 3
        race_chart = {
            2: {2: 2, 3: 2, 4: 2, 5: 2, 6: 2, 7: 2},
            3: {2: 3, 3: 2, 4: 2, 5: 2, 6: 2, 7: 2},
            4: {2: 4, 3: 3, 4: 3, 5: 3, 6: 3, 7: 2},
            5: {2: 5, 3: 4, 4: 4, 5: 4, 6: 4, 7: 3},
            6: {2: 6, 3: 5, 4: 5, 5: 5, 6: 5, 7: 4},
            7: {2: 7, 3: 6, 4: 5, 5: 5, 6: 5, 7: 5},
        }
        try:
            games1 = race_chart[skill1][skill2]
            games2 = race_chart[skill2][skill1]
            return games1, games2
        except KeyError:
            raise ValueError("Skill levels must be between 2 and 7 inclusive.")


    def lookup_adjusted_skill(self, running_avg):
        running_avg = round(running_avg, 2)
        for skill, (low, high) in self.skill_levels_eight.items():
            if low <= running_avg <= high:
                low, high = self.skill_levels_eight[skill][0], self.skill_levels_eight[skill][1]
                adjusted_skill = skill + 1 - (running_avg - low) / (high - low)
                return adjusted_skill
        return None

    def calculate_running_average(self):
        running_averages = []
        adjusted_skill_levels = []
        adjusted_equalizer_skill_levels = []
        for index, row in self.df_eight_ball.iterrows():
            prev_scores = self.df_eight_ball.loc[:index, 'adjusted_innings_per_rack'].dropna().tail(20)
            n_scores = len(prev_scores)
            if n_scores == 0:
                best_scores_count = 0
            elif n_scores < 3:
                best_scores_count = n_scores
            elif n_scores < 20:
                best_scores_count = math.ceil(n_scores / 2)
            else:
                best_scores_count = 10
            best_scores = prev_scores.nsmallest(best_scores_count)
            running_avg = best_scores.mean()
            running_averages.append(running_avg)
            adj_skill = self.lookup_adjusted_skill(running_avg)
            adjusted_skill_levels.append(adj_skill)
            adjusted_eq = (
                # 0.677 * row.get('skillLevel', 0) +
                # 0.000 * row.get('innings', 0) +
                # -0.000 * row.get('defensiveShots', 0) +
                # 0.081 * row.get('eightBallWins', 0) +
                # -0.026 * row.get('eightBallLosses', 0) +
                # 0.041 * row.get('opponentSkillLevel', 0) +
                # 0.005 * row.get('tableSize', 0) +
                # 0.001 * row.get('effective_innings', 0) +
                # -0.005 * row.get('innings_per_rack', 0) +
                # -0.147 * row.get('adjusted_innings_per_rack', 0) +
                -2.101 * row.get('running_win_rate', 0) +
                -0.116 * running_avg +
                0.924 * adj_skill +
                # -0.746 * (1 if row.get('winLoss', 'L') == 'W' else 0) +
                1.326
            )
            adjusted_equalizer_skill_levels.append(adjusted_eq)
        self.df_eight_ball['running_innings_per_rack'] = running_averages
        self.df_eight_ball['adjusted_skill_level'] = adjusted_skill_levels
        # self.df_eight_ball['adjusted_equalizer_skill_level'] = adjusted_equalizer_skill_levels

    def calculate_adjusted_innings(self):
        adjusted_points = []
        running_win_rates = []
        for index, row in self.df_eight_ball.iterrows():
            skill_level = row.get('skillLevel', 1)
            if skill_level == 0:
                skill_level = 3
            match_wins_needed, match_wins_opponent_needed = self.apa_8_ball_race(skill_level, row.get('opponentSkillLevel', 3))
            eight_ball_wins = row.get('eightBallWins', 0)
            eight_ball_losses = row.get('eightBallLosses', 0)
            # if eight_ball_wins < match_wins_needed and eight_ball_losses < match_wins_opponent_needed:
            #     adjusted_points.append(np.nan)
            #     continue
            is_win = row.get('winLoss', '') == 'W'
            prev_scores = self.df_eight_ball.loc[:index, ['winLoss', 'innings_per_rack']].dropna().tail(20)
            win_count = len(prev_scores[prev_scores['winLoss'] == 'W'])
            loss_count = len(prev_scores[prev_scores['winLoss'] == 'L'])
            total_games = win_count + loss_count
            win_percentage = (win_count / total_games * 100) if total_games > 0 else 0.0
            running_win_rates.append(win_percentage / 100.0)
            low, high = EightBallPlayer.expected_innings_eight.get(skill_level, (0, 0))
            low = low * self.expected_eight_scaler
            high = high * self.expected_eight_scaler
            adjusted_innings_per_rack = high - (high - low) * (win_percentage / 100)

            loss_penalty_ceiling = low + high

            if loss_penalty_ceiling == 0:
                table_size_adjustment = 0.0
            elif row.get('tableSize', '7') == 8:
                table_size_adjustment = 1 / loss_penalty_ceiling
            elif row.get('tableSize', '7') == 9:
                table_size_adjustment = 2 / loss_penalty_ceiling
            else:
                table_size_adjustment = 0.0
            innings_per_rack = row.get('innings_per_rack', 0) - table_size_adjustment
            if is_win:
                innings_to_use = min(innings_per_rack, adjusted_innings_per_rack)
                adjusted_points.append(innings_to_use)
            else:
                total_racks = eight_ball_wins + eight_ball_losses
                win_rate_in_match = eight_ball_wins / total_racks if total_racks > 0 else 0.0
                if eight_ball_wins < match_wins_needed and eight_ball_losses < match_wins_opponent_needed:
                    adjusted_innings_per_rack = loss_penalty_ceiling
                else:
                    # percent_of_wins_reached = eight_ball_wins / match_wins_needed
                    # adjusted_innings_per_rack = loss_penalty_ceiling - (loss_penalty_ceiling - innings_per_rack) * percent_of_wins_reached
                    if win_rate_in_match == 0.0:
                        adjusted_innings_per_rack = loss_penalty_ceiling
                    else:
                        adjusted_innings_per_rack = loss_penalty_ceiling - (loss_penalty_ceiling - innings_per_rack) * win_rate_in_match
                adjusted_points.append(adjusted_innings_per_rack)
        self.df_eight_ball['adjusted_innings_per_rack'] = adjusted_points
        self.df_eight_ball['running_win_rate'] = running_win_rates

    def process(self):
        if 'innings' in self.df_eight_ball.columns and 'defensiveShots' in self.df_eight_ball.columns:
            self.df_eight_ball['effective_innings'] = self.df_eight_ball['innings'] - self.df_eight_ball['defensiveShots']
        if 'eightBallWins' in self.df_eight_ball.columns and 'eightBallLosses' in self.df_eight_ball.columns and 'effective_innings' in self.df_eight_ball.columns:
            self.df_eight_ball['innings_per_rack'] = self.df_eight_ball['effective_innings'] / (self.df_eight_ball['eightBallWins'] + self.df_eight_ball['eightBallLosses'])
        self.calculate_adjusted_innings()
        self.calculate_running_average()
        return self.df_eight_ball

def process_eight_ball(df):
    player = EightBallPlayer(df)
    return player.process()
