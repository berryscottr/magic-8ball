import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import seaborn as sns


class SLMatchesNine:
    def __init__(self, SL, opponentSL):
        self.SL = SL
        self.opponentSL = opponentSL
        self.points = 0
        self.games = 0
        self.indices_checked = []
        self.average = None

    def addpoints(self, points):
        self.points += points

    def addgame(self, index):
        self.games += 1
        self.indices_checked.append(index)

    def getaverage(self):
        try:
            self.average = self.points / self.games
        except ZeroDivisionError:
            self.average = None


def get_sl_matchup_stats_nine(df, games2win):
    slrange = range(1, 10)
    slmatches_average = games2win
    for p1skill in slrange:
        for p2skill in slrange:
            matchup_data = SLMatchesNine(p1skill, p2skill)
            if p1skill == p2skill:
                slmatches_average.loc[p1skill, p2skill] = 10.00
            else:
                for index, row in df.iterrows():
                    if row['Player_1'] == p1skill and row['Player_2'] == p2skill:
                        matchup_data.addpoints(row['Points_1'])
                        matchup_data.addgame(index)
                        if p1skill == p2skill:
                            matchup_data.addpoints(row['Points_2'])
                for index, row in df.iterrows():
                    if index not in matchup_data.indices_checked:
                        if row['Player_2'] == p1skill and row['Player_1'] == p2skill:
                            matchup_data.addpoints(row['Points_2'])
                            matchup_data.addgame(index)
                matchup_data.getaverage()
                if matchup_data.average is None:
                    slmatches_average.loc[p1skill, p2skill] = 0.00  # alternatively could set to None
                else:
                    try:
                        slmatches_average.loc[p1skill, p2skill] = round(matchup_data.average, 2)
                    except TypeError:
                        slmatches_average.loc[p1skill, p2skill] = matchup_data.average
    sls = pd.DataFrame(slmatches_average, index=slrange, columns=slrange, dtype=float)
    sns.heatmap(sls, annot=True, cmap=sns.color_palette("coolwarm", 12), vmin=0, vmax=20, fmt=".2f",
                linewidths=.2, cbar_kws={"label": "Average Points"})
    plt.title("Opponent SL", size=10)
    plt.xlabel("Opponent SL")
    plt.ylabel("Player SL")
    plt.tick_params(axis='both', which='major', labelsize=10, labelbottom=False, bottom=False, top=False, labeltop=True,
                    left=False, right=False)
    plt.savefig("../../data/images/slMatchupAveragesNine.svg")
    slmatches_average.to_excel(r'../../data/SLMatchupAveragesNine.xlsx', index=True, header=True, sheet_name="Nine")
    print(slmatches_average)
    plt.clf()

def get_sl_matchup_stats_nine_jason_version(df):
    total_score_matrix = [ [0]*10 for i in range(10)]
    matchup_count_matrix = [ [0]*10 for i in range(10)]
    average_score_matrix = [ [0]*10 for i in range(10)]
    for index, row in df.iterrows():
        sl_1 = row['Player_1']
        sl_2 = row['Player_2']
        points_1 = row['Points_1']
        points_2 = row['Points_2']

        total_score_matrix[sl_1][sl_2] += points_1
        total_score_matrix[sl_2][sl_1] += points_2
        matchup_count_matrix[sl_1][sl_2] += 1
        matchup_count_matrix[sl_2][sl_1] += 1
    
    for x in range(1, 10):
        for y in range(1, 10):
            if (matchup_count_matrix[x][y] != 0):
                average_score_matrix[x][y] = round(
                    total_score_matrix[x][y] / matchup_count_matrix[x][y],
                    2)
            else:
                average_score_matrix[x][y] = -1
    
    print(np.matrix(average_score_matrix))
