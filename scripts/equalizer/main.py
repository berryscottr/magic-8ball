import os
import sys
from concurrent.futures import ThreadPoolExecutor
import util as equalizer_util

def main():
    players_dir = "data/players/"
    if not os.path.exists(players_dir):
        print(f"Directory {players_dir} does not exist.")
        sys.exit(1)

    def process_player(playerID):
        filename = os.path.join(players_dir, playerID, "matches.json")
        if os.path.isfile(filename):
            df = equalizer_util.json_to_dataframe(filename)
            df_eight_ball, df_nine_ball = equalizer_util.separate_dataframes(df)
            equalizer_util.save_csv_files(df_eight_ball, df_nine_ball, filename)

    with ThreadPoolExecutor() as executor:
        player_ids = os.listdir(players_dir)
        results = executor.map(process_player, player_ids)

if __name__ == "__main__":
    main()
    