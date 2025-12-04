import os
import sys
import glob
import math
import pandas as pd
import numpy as np
from sklearn.impute import SimpleImputer
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error
from concurrent.futures import ThreadPoolExecutor

class PlayerModel:
    def __init__(self, players_dir="data/players/"):
        self.players_dir = players_dir
        self.combined_df = pd.DataFrame()
        self.model = None
        self.imputer = None
        self.feature_cols = None

    def load_player_data(self, player_dir):
        """Load CSVs for one player, create next_skillLevel, convert winLoss to binary."""
        dfs = []
        csv_files = glob.glob(os.path.join(player_dir, "eightball.csv"))

        for csv_file in csv_files:
            df = pd.read_csv(csv_file)

            # Verify required columns exist
            if 'skillLevel' not in df.columns or 'winLoss' not in df.columns:
                print(f"⚠️ Skipping {csv_file} — missing 'skillLevel' or 'winLoss'")
                continue

            df['next_skillLevel'] = df['skillLevel'].shift(-1)
            df = df.dropna(subset=['next_skillLevel'])

            # Convert winLoss to binary
            df['winLossBinary'] = df['winLoss'].map({'W': 1, 'L': 0}).fillna(0)

            # Track player ID
            df['player_id'] = os.path.basename(player_dir)
            dfs.append(df)

        if dfs:
            return pd.concat(dfs, ignore_index=True)
        return pd.DataFrame()

    def load_all_players(self):
        if not os.path.exists(self.players_dir):
            print(f"Directory {self.players_dir} does not exist.")
            sys.exit(1)

        all_dfs = []
        with ThreadPoolExecutor() as executor:
            player_dirs = [os.path.join(self.players_dir, pid) for pid in os.listdir(self.players_dir)]
            results = executor.map(self.load_player_data, player_dirs)

        for df in results:
            if not df.empty:
                all_dfs.append(df)

        if not all_dfs:
            print("No player CSV data found.")
            sys.exit(1)

        self.combined_df = pd.concat(all_dfs, ignore_index=True)
        print(f"Combined dataset size: {self.combined_df.shape}")

    def train_global_model(self):
        """Train regression model on combined data including winLossBinary."""
        # Select features: everything except next_skillLevel, winLoss, player_id
        X = self.combined_df.drop(columns=['eightBallWins', 'winLossBinary', 'adjusted_innings_per_rack', 'eightBallLosses', 'opponentSkillLevel', 'skillLevel', 'next_skillLevel', 'winLoss', 'player_id', 'adjusted_equalizer_skill_level', 'innings', 'defensiveShots', 'tableSize', 'effective_innings', 'innings_per_rack'])
        y = self.combined_df['next_skillLevel']

        X = X.replace([np.inf, -np.inf], np.nan)
        X = X.dropna(axis=1, how='all')

        self.imputer = SimpleImputer(strategy="mean")
        X_imputed = self.imputer.fit_transform(X)

        if not np.isfinite(X_imputed).all():
            raise ValueError("NaN or Inf values remain after cleaning.")

        self.model = LinearRegression()
        self.model.fit(X_imputed, y)
        self.feature_cols = X.columns

    def evaluate_per_player(self):
        """Evaluate global model for each player and return metrics."""
        metrics = []
        for pid in self.combined_df['player_id'].unique():
            player_df = self.combined_df[self.combined_df['player_id'] == pid]
            X_player = player_df[self.feature_cols]
            y_player = player_df['next_skillLevel']

            X_player = X_player.replace([np.inf, -np.inf], np.nan)
            X_player_imputed = self.imputer.transform(X_player)

            preds = self.model.predict(X_player_imputed)
            preds_floor = np.floor(preds)

            mse = mean_squared_error(y_player, preds_floor)
            rmse = np.sqrt(mse)

            metrics.append({"player_id": pid, "rmse": rmse})
        return metrics

    def print_model_formula(self):
        """Print the linear regression formula."""
        coef_str = " + ".join([f"({c:.3f} * {f})" for c, f in zip(self.model.coef_, self.feature_cols)])
        print(f"\nLinear regression formula:\nnext_skillLevel ≈ {coef_str} + ({self.model.intercept_:.3f})")

    def run(self):
        self.load_all_players()
        self.train_global_model()
        self.print_model_formula()
        metrics = self.evaluate_per_player()
        print("\nPlayer RMSEs:")
        for m in metrics:
            print(f"{m['player_id']}: RMSE = {m['rmse']:.3f}")


if __name__ == "__main__":
    player_model = PlayerModel()
    player_model.run()
