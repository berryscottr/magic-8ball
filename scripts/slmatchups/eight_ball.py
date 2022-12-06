import dataframe_image as dfi
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns


class SLMatchesEight:
    def __init__(self, SL, opponentSL):
        self.SL = SL
        self.opponentSL = opponentSL
        self.points = 0
        self.games = 0
        self.indices_checked = []
        self.average = None
        self.median = None
        self.mode = None
        self.resultcounts = {
            "0/3": 0,
            "0/2": 0,
            "1/2": 0,
            "2/1": 0,
            "2/0": 0,
            "3/0": 0,
        }

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

    def getmedian(self):
        median_game_index = int(sum(self.resultcounts.values()) / 2)
        count = 0
        for k, v in self.resultcounts.items():
            if count <= median_game_index < count + v:
                self.median = k
            count += v

    def getmode(self):
        self.mode = max(self.resultcounts, key=self.resultcounts.get)


def get_sl_matchup_stats_eight(df, games2win):
    slrange = range(2, 8)
    # get average
    slmatches_average = games2win
    for p1skill in slrange:
        for p2skill in slrange:
            matchup_data = SLMatchesEight(p1skill, p2skill)
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
            if p1skill == p2skill:
                matchup_data.average /= 2
            # if matchup_data.games < 6:
            #     matchup_data.average = "X{}".format(round(matchup_data.average, 2))
            try:
                slmatches_average.loc[p1skill, p2skill] = round(matchup_data.average, 2)
            except TypeError:
                slmatches_average.loc[p1skill, p2skill] = matchup_data.average
    sls = pd.DataFrame(slmatches_average, index=slrange, columns=slrange, dtype=float)
    sns.heatmap(sls, annot=True, cmap=sns.color_palette("coolwarm", 12), vmin=0, vmax=3, fmt=".2f",
                linewidths=.2, cbar_kws={"label": "Average Points"})
    plt.title("Opponent SL", size=10)
    plt.xlabel("Opponent SL")
    plt.ylabel("Player SL")
    plt.tick_params(axis='both', which='major', labelsize=10, labelbottom=False, bottom=False, top=False, labeltop=True,
                    left=False, right=False)
    plt.savefig("../../data/images/slMatchupAverages.svg")
    slmatches_average.to_excel(r'../../data/SLMatchupAverages.xlsx', index=True, header=True)
    print(slmatches_average)
    # get median and mode
    slmatches_median = games2win
    slmatches_mode = games2win
    for p1skill in slrange:
        for p2skill in slrange:
            matchup_data = SLMatchesEight(p1skill, p2skill)
            for index, row in df.iterrows():
                if row['Player_1'] == p1skill and row['Player_2'] == p2skill:
                    if row['Points_1'] == 0 and row['Points_2'] == 3:
                        matchup_data.resultcounts["0/3"] += 1
                    elif row['Points_1'] == 0 and row['Points_2'] == 2:
                        matchup_data.resultcounts["0/2"] += 1
                    elif row['Points_1'] == 1 and row['Points_2'] == 2:
                        matchup_data.resultcounts["1/2"] += 1
                    elif row['Points_1'] == 2 and row['Points_2'] == 1:
                        matchup_data.resultcounts["2/1"] += 1
                    elif row['Points_1'] == 2 and row['Points_2'] == 0:
                        matchup_data.resultcounts["2/0"] += 1
                    elif row['Points_1'] == 3 and row['Points_2'] == 0:
                        matchup_data.resultcounts["3/0"] += 1
                    matchup_data.addgame(index)
            for index, row in df.iterrows():
                if index not in matchup_data.indices_checked:
                    if row['Player_2'] == p1skill and row['Player_1'] == p2skill:
                        if row['Points_2'] == 0 and row['Points_1'] == 3:
                            matchup_data.resultcounts["0/3"] += 1
                        elif row['Points_2'] == 0 and row['Points_1'] == 2:
                            matchup_data.resultcounts["0/2"] += 1
                        elif row['Points_2'] == 1 and row['Points_1'] == 2:
                            matchup_data.resultcounts["1/2"] += 1
                        elif row['Points_2'] == 2 and row['Points_1'] == 1:
                            matchup_data.resultcounts["2/1"] += 1
                        elif row['Points_2'] == 2 and row['Points_1'] == 0:
                            matchup_data.resultcounts["2/0"] += 1
                        elif row['Points_2'] == 3 and row['Points_1'] == 0:
                            matchup_data.resultcounts["3/0"] += 1
                        matchup_data.addgame(index)
            matchup_data.getmedian()
            try:
                slmatches_median.loc[p1skill, p2skill] = matchup_data.median
            except TypeError:
                slmatches_median.loc[p1skill, p2skill] = matchup_data.median
            matchup_data.getmode()
            try:
                slmatches_mode.loc[p1skill, p2skill] = matchup_data.mode
            except TypeError:
                slmatches_mode.loc[p1skill, p2skill] = matchup_data.mode
    sls = pd.DataFrame(slmatches_median, index=slrange, columns=slrange, dtype=str)
    dfi.export(sls, "../../data/images/slMatchupMedians.png")
    slmatches_median.to_excel(r'../../data/SLMatchupMedians.xlsx', index=True, header=True)
    print(slmatches_median)
    sls = pd.DataFrame(slmatches_mode, index=slrange, columns=slrange, dtype=str)
    dfi.export(sls, "../../data/images/slMatchupModes.png")
    slmatches_mode.to_excel(r'../../data/SLMatchupModes.xlsx', index=True, header=True)
    print(slmatches_mode)
    plt.clf()

    # for sl in slrange:
    #     plt.plot(np.array(slrange), np.array([float(str(x).replace("X", "", 1)) for x in slmatches.values[sl-2]]))
    #     title = 'Matchups SL {}'.format(sl)
    #     imgtitle = "../data/slmatchups{}.jpg".format(sl)
    #     plt.title(title)
    #     plt.xlabel('Opponent SL')
    #     plt.ylabel('Average Points')
    #     plt.savefig(imgtitle)
    #     plt.show()
