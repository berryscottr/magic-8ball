import json
import pandas as pd
import os
import sys
from concurrent.futures import ThreadPoolExecutor

def process_player(playerID):
    try:
        players_dir = "data/players/"
        filename = os.path.join(players_dir, playerID, "matches.json")
        if os.path.isfile(filename):
            # Load the JSON data from file
            with open(filename, 'r') as f:
                data = json.load(f)

            # Convert to a DataFrame
            df = pd.DataFrame(data)

            # Convert the dateTimeStamp column to datetime objects
            df['dateTimeStamp'] = pd.to_datetime(df['dateTimeStamp'], errors='coerce')
            df = df.dropna(subset=['dateTimeStamp'])

            # Separate into EightBall and NineBall matches based on player's __typename
            df_eightball = df[df['player'].apply(lambda p: p.get('__typename') == 'EightBallPlayer')].copy()
            df_nineball = df[df['player'].apply(lambda p: p.get('__typename') == 'NineBallPlayer')].copy()

            # Process EightBall matches
            if not df_eightball.empty:
                # Sort eight-ball matches by date
                df_eightball = df_eightball.sort_values('dateTimeStamp')
                most_recent_eight = df_eightball.iloc[-1]
                current_eight_skill = most_recent_eight['skillLevel']

                # Filter eight-ball matches played at the current skill level
                df_eightball_current = df_eightball[df_eightball['skillLevel'] == current_eight_skill]
                # Limit to the most recent 60 matches
                df_eightball_current = df_eightball_current.sort_values('dateTimeStamp', ascending=False).head(60)

                # Group by opponent skill level and aggregate means and sample size
                expected_eightball = df_eightball_current.groupby('opponentSkillLevel').agg(
                    expectedEightBallPoints=('eightBallMatchPointsEarned', 'mean'),
                    expectedEightBallPointsLost=('eightBallMatchPointsLost', 'mean'),
                    sampleSize=('eightBallMatchPointsEarned', 'count')
                ).reset_index()

                # Round the expected values to two decimals
                expected_eightball['expectedEightBallPoints'] = expected_eightball['expectedEightBallPoints'].round(2)
                expected_eightball['expectedEightBallPointsLost'] = expected_eightball['expectedEightBallPointsLost'].round(2)

                # Save the output as CSV
                output_path = os.path.join(players_dir, playerID, "eight_matchups.csv")
                expected_eightball.to_csv(output_path, index=False)
                print(f"Player {playerID}: Written eight_matchups.csv, SL Matchups:\n{expected_eightball.to_string(index=False)}\n")
            else:
                print(f"Player {playerID}: No EightBall matches found.\n")

            # Process NineBall matches
            if not df_nineball.empty:
                # Sort nine-ball matches by date
                df_nineball = df_nineball.sort_values('dateTimeStamp')
                most_recent_nine = df_nineball.iloc[-1]
                current_nine_skill = most_recent_nine['skillLevel']

                # Filter nine-ball matches played at the current skill level
                df_nineball_current = df_nineball[df_nineball['skillLevel'] == current_nine_skill]
                # Limit to the most recent 60 matches
                df_nineball_current = df_nineball_current.sort_values('dateTimeStamp', ascending=False).head(60)

                # Group by opponent skill level and aggregate means and sample size
                expected_nineball = df_nineball_current.groupby('opponentSkillLevel').agg(
                    expectedNineBallPoints=('nineBallMatchPointsEarned', 'mean'),
                    sampleSize=('nineBallMatchPointsEarned', 'count')
                ).reset_index()

                # Round the expected values to two decimals
                expected_nineball['expectedNineBallPoints'] = expected_nineball['expectedNineBallPoints'].round(2)

                # Save the output as CSV
                output_path = os.path.join(players_dir, playerID, "nine_matchups.csv")
                expected_nineball.to_csv(output_path, index=False)
                print(f"Player {playerID}: Written nine_matchups.csv, SL Matchups:\n{expected_nineball.to_string(index=False)}\n")
            else:
                print(f"Player {playerID}: No NineBall matches found.\n")
        else:
            print(f"Player {playerID}: File {filename} does not exist.\n")
    except Exception as e:
        print(f"Error processing player {playerID}: {e}", file=sys.stderr)

def main():
    players_dir = "data/players/"
    if not os.path.exists(players_dir):
        print(f"Directory {players_dir} does not exist.\n")
        sys.exit(1)

    # Process each player's matches concurrently
    player_ids = [d for d in os.listdir(players_dir) if os.path.isdir(os.path.join(players_dir, d))]
    with ThreadPoolExecutor() as executor:
        executor.map(process_player, player_ids)

if __name__ == "__main__":
    main()
