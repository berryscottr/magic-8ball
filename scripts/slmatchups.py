import openpyxl
import pandas as pd
import warnings
warnings.simplefilter(action='ignore', category=FutureWarning)


def workbook2df(path, first_row_header, first_column_index):
    workbook = openpyxl.load_workbook(path).active
    rows = list(workbook.iter_rows(values_only=True))
    df = pd.DataFrame(rows)
    if first_row_header:
        new_header = df.iloc[0]
        df = df[1:]
        df.columns = new_header
    if first_column_index:
        df = df.set_index("SL")
    return df


class SLMatches:
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


def get_sl_matchup_stats(df, games2win):
    slmatches = games2win
    slrange = range(2, 8)
    for p1skill in slrange:
        for p2skill in slrange:
            matchup_data = SLMatches(p1skill, p2skill)
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
            if matchup_data.games < 6:
                matchup_data.average = "X{}".format(round(matchup_data.average, 2))
            try:
                slmatches.loc[p1skill, p2skill] = round(matchup_data.average, 2)
            except TypeError:
                slmatches.loc[p1skill, p2skill] = matchup_data.average
    # for sl in slrange:
    #     plt.plot(np.array(slrange), np.array([float(str(x).replace("X", "", 1)) for x in slmatches.values[sl-2]]))
    #     title = 'Matchups SL {}'.format(sl)
    #     imgtitle = "../data/slmatchups{}.jpg".format(sl)
    #     plt.title(title)
    #     plt.xlabel('Opponent SL')
    #     plt.ylabel('Average Points')
    #     plt.savefig(imgtitle)
    #     plt.show()
    slmatches.to_excel(r'../data/SLMatchupAverages.xlsx', index=True, header=True)
    print(slmatches)


def main():
    games_to_win = workbook2df("../data/GamesToWin.xlsx", True, True)
    sl_matchup_data = workbook2df("../data/SLMatchups.xlsx", True, False)
    # player_data = workbook2df("data/wookieMistakesPlayerData.xlsx", True, False)
    # game_data = workbook2df("data/wookieMistakesSpring2022Games.xlsx", True, False)
    get_sl_matchup_stats(sl_matchup_data, games_to_win)


if __name__ == '__main__':
    main()
