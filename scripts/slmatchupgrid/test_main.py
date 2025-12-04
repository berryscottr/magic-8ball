import main as slmatchups


def test_main():
    games_to_win = slmatchups.workbook2df("../../data/GamesToWin.xlsx", True, True)
    assert games_to_win is not None
    sl_matchup_data_eight = slmatchups.workbook2df("../../data/SLMatchups.xlsx", True, False)
    assert sl_matchup_data_eight is not None
    games_to_win_nine = slmatchups.workbook2df("../../data/GamesToWinNine.xlsx", True, True)
    assert games_to_win_nine is not None
    sl_matchup_data_nine = slmatchups.workbook2df("../../data/SLMatchupsNine.xlsx", True, False)
    assert sl_matchup_data_nine is not None


if __name__ == '__main__':
    test_main()
