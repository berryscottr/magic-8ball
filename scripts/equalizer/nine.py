import pandas as pd
import math

class NineBallPlayer:
    # Assume 8 foot table if not specified
    expected_points_nine_7_foot = {
        1: (0.0, (14/15)),
        2: ((14/15), (19/15)),
        3: ((19/15), (25/15)),
        4: ((25/15), (31/15)),
        5: ((31/15), (38/15)),
        6: ((38/15), (46/15)),
        7: ((46/16), (55/13)),
        8: ((55/13), (65/13)),
        9: ((65/13), (75/13))
    }
    expected_points_nine_8_foot = {
        1: (0.0, (14/16)),
        2: ((14/16), (19/16)),
        3: ((19/16), (25/16)),
        4: ((25/16), (31/16)),
        5: ((31/16), (38/16)),
        6: ((38/16), (46/16)),
        7: ((46/16), (55/14)),
        8: ((55/14), (65/14)),
        9: ((65/14), (75/14))
    }
    expected_points_nine_9_foot = {
        1: (0.0, (14/17)),
        2: ((14/17), (19/17)),
        3: ((19/17), (25/17)),
        4: ((25/17), (31/17)),
        5: ((31/17), (38/17)),
        6: ((38/17), (46/17)),
        7: ((46/17), (55/15)),
        8: ((55/15), (65/15)),
        9: ((65/15), (75/15))
    }
    expected_nine_scaler = 1.0

    def __init__(self, df):
        self.df_nine_ball = df[df['player___typename'] == 'NineBallPlayer'][['skillLevel', 'innings', 'defensiveShots', 'nineBallPoints', 'winLoss', 'tableSize']]
        self.df_nine_ball['nineBallPoints'] = pd.to_numeric(self.df_nine_ball['nineBallPoints'], errors='coerce').fillna(0).astype(int)

    def lookup_adjusted_skill_nine(self, running_avg):
        for skill, (low, high) in NineBallPlayer.expected_points_nine_8_foot.items():
            if low <= running_avg < high:
                return skill + (running_avg - low) / (high - low)
        return None

    def calculate_running_average_nine(self):
        running_averages = []
        adjusted_skill_levels = []
        for index, row in self.df_nine_ball.iterrows():
            prev_scores = self.df_nine_ball.loc[:index, 'adjusted_points_per_inning'].dropna().tail(20)
            n_scores = len(prev_scores)
            if n_scores == 0:
                best_scores_count = 0
            elif n_scores < 3:
                best_scores_count = n_scores
            elif n_scores < 20:
                best_scores_count = math.ceil(n_scores / 2)
            else:
                best_scores_count = 10
            best_scores = prev_scores.nlargest(best_scores_count)
            running_avg = best_scores.mean()
            running_averages.append(running_avg)
            adjusted_skill_levels.append(self.lookup_adjusted_skill_nine(running_avg))
        self.df_nine_ball['running_points_per_inning'] = running_averages
        self.df_nine_ball['adjusted_skill_level'] = adjusted_skill_levels

    def calculate_adjusted_points(self):
        adjusted_points = []
        for index, row in self.df_nine_ball.iterrows():
            skill_level = row.get('skillLevel', 1)
            if skill_level == 0:
                skill_level = 3
            points_per_inning = row.get('points_per_inning', 0)
            win_loss = row.get('winLoss', '')
            prev_scores = self.df_nine_ball.loc[:index, ['winLoss', 'points_per_inning']].dropna().tail(20)
            win_count = len(prev_scores[prev_scores['winLoss'] == 'W'])
            loss_count = len(prev_scores[prev_scores['winLoss'] == 'L'])
            win_percentage = win_count / (win_count + loss_count) * 100
            if row.get('tableSize', '8') == '7':
                low, high = NineBallPlayer.expected_points_nine_7_foot.get(skill_level, (0, 0))
            if row.get('tableSize', '8') == '9':
                low, high = NineBallPlayer.expected_points_nine_9_foot.get(skill_level, (0, 0))
            else:
                low, high = NineBallPlayer.expected_points_nine_8_foot.get(skill_level, (0, 0))
            low = low * NineBallPlayer.expected_nine_scaler
            high = high * NineBallPlayer.expected_nine_scaler
            adjusted_points_per_inning = low + (high - low) * (win_percentage / 100)
            if win_loss == 'W' and points_per_inning < adjusted_points_per_inning:
                adjusted_points.append(adjusted_points_per_inning)
            else:
                adjusted_points.append(points_per_inning)
        self.df_nine_ball['adjusted_points_per_inning'] = adjusted_points

    def process(self):
        if 'innings' in self.df_nine_ball.columns and 'defensiveShots' in self.df_nine_ball.columns:
            self.df_nine_ball['effective_innings'] = self.df_nine_ball['innings'] - self.df_nine_ball['defensiveShots']
        if 'nineBallPoints' in self.df_nine_ball.columns and 'effective_innings' in self.df_nine_ball.columns:
            self.df_nine_ball['points_per_inning'] = self.df_nine_ball['nineBallPoints'] / self.df_nine_ball['effective_innings']
        
        self.calculate_adjusted_points()
        self.calculate_running_average_nine()
        return self.df_nine_ball


def process_nine_ball(df):
    player = NineBallPlayer(df)
    return player.process()
