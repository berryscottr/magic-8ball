import warnings
import eight_ball
import nine_ball
from pool_util import workbook2df

warnings.simplefilter(action='ignore', category=FutureWarning)


def main():
    games_to_win = workbook2df("../../data/GamesToWin.xlsx", True, True)
    sl_matchup_data_eight = workbook2df("../../data/SLMatchups.xlsx", True, False)
    # player_data = workbook2df("data/wookieMistakesPlayerData.xlsx", True, False)
    # game_data = workbook2df("data/wookieMistakesSpring2022Games.xlsx", True, False)
    eight_ball.get_sl_matchup_stats_eight(sl_matchup_data_eight, games_to_win)
    games_to_win_nine = workbook2df("../../data/GamesToWinNine.xlsx", True, True)
    sl_matchup_data_nine = workbook2df("../../data/SLMatchupsNine.xlsx", True, False)
    nine_ball.get_sl_matchup_stats_nine(sl_matchup_data_nine, games_to_win_nine)


if __name__ == '__main__':
    main()
