import json
import pandas as pd
import os
import threading
from concurrent.futures import ThreadPoolExecutor
from eight import process_eight_ball
from nine import process_nine_ball

file_lock = threading.Lock()

def json_to_dataframe(file_path):
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            json_data = json.load(file)
    except Exception as e:
        print(f"[ERROR] Failed to load JSON from {file_path}: {e}")
        raise

    try:
        df = pd.json_normalize(json_data, sep='_')
        df['dateTimeStamp'] = pd.to_datetime(df['dateTimeStamp'], errors='coerce')
        df = df.dropna(subset=['dateTimeStamp'])
        df = df.sort_values(by='dateTimeStamp')
        return df
    except Exception as e:
        print(f"[ERROR] Failed to process DataFrame from JSON: {e}")
        raise

def separate_dataframes(df):
    with ThreadPoolExecutor() as executor:
        future_eight = executor.submit(safe_process, process_eight_ball, df, "Eight Ball")
        future_nine = executor.submit(safe_process, process_nine_ball, df, "Nine Ball")

        df_eight_ball = future_eight.result()
        df_nine_ball = future_nine.result()

    return df_eight_ball, df_nine_ball

def safe_process(func, df, label):
    try:
        return func(df)
    except Exception as e:
        print(f"[ERROR] {label} processing failed: {e}")
        return pd.DataFrame()

def save_csv_files(df_eight_ball, df_nine_ball, file_path):
    dir_name = os.path.dirname(file_path)
    player_name = os.path.basename(dir_name).replace('\\', '/').split('/')[-1]

    eight_ball_filename = os.path.join(dir_name, "eightball.csv")
    nine_ball_filename = os.path.join(dir_name, "nineball.csv")

    with file_lock:
        df_eight_ball.to_csv(eight_ball_filename, index=False, float_format='%.2f')
        df_nine_ball.to_csv(nine_ball_filename, index=False, float_format='%.2f')

        try:
            print(f"{player_name} - 8-ball - {df_eight_ball['adjusted_skill_level'].iloc[-1]:.2f}")
        except IndexError:
            pass  # No 8-ball data, skip silently
        except Exception as e:
            print(f"{player_name} - 8-ball - ERROR: {e}")

        try:
            print(f"{player_name} - 9-ball - {df_nine_ball['adjusted_skill_level'].iloc[-1]:.2f}")
        except IndexError:
            pass  # No 9-ball data, skip silently
        except Exception as e:
            print(f"{player_name} - 9-ball - ERROR: {e}")


