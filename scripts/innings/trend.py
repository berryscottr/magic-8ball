import matplotlib.pyplot as plt
import numpy as np

class RollingAveragePlotter:
    def __init__(self, data_points):
        self.data_points = data_points
        self.rolling_average = []

    def calculate_rolling_average(self):
        for i in range(len(self.data_points)):
            if i >= 19:
                last_20_sorted = sorted(self.data_points[i-19:i+1], reverse=True)[:10]
                avg = np.mean(last_20_sorted)
                self.rolling_average.append(avg)
            else:
                n = i + 1
                best_n_half = sorted(self.data_points[:i+1], reverse=True)[:int(np.ceil(n / 2))]
                avg = np.mean(best_n_half)
                self.rolling_average.append(avg)

    def plot(self):
        self.calculate_rolling_average()
        x_values = range(1, len(self.data_points) + 1)

        plt.figure(figsize=(10, 6))
        plt.scatter(x_values, self.data_points, color='blue', label='Data Points')
        plt.plot(x_values, self.rolling_average, color='red', label='Rolling Average of Good Games')
        plt.title('Plot of the Provided Data Points with Rolling Average')
        plt.xlabel('Index')
        plt.ylabel('Value')
        plt.legend()
        plt.grid(True, which='both', axis='both', linestyle='--', linewidth=0.5)
        plt.xticks(np.arange(min(x_values), max(x_values)+1, 1))
        plt.yticks(np.arange(min(self.data_points), max(self.data_points)+0.1, 0.1))
        plt.show()

# Provided data points
data_points = [
    1.06, 0.92, 1.48, 1.48, 1.16, 1.55, 1.94, 1.48, 2.82, 1.45, 
    1.48, 0.91, 1.48, 1.48, 1.55, 1.48, 2.21, 1.15, 1.48, 1.48, 
    2.07, 1.48, 1.48, 1.48, 2.38, 1.81, 2.24, 1.81, 1.10, 2.53, 
    1.81, 1.36, 2.11, 2.24, 2.11, 1.50, 1.84, 1.81, 1.81, 1.81, 
    1.10, 1.81, 1.81, 1.90, 1.30, 2.24, 4.22, 2.38, 1.71, 2.71, 
    1.69, 4.60, 1.14, 1.52, 0.96, 1.75, 1.68, 0.42, 1.60, 2.19,
    3.54, 2.19, 2.19, 2.19, 2.30, 3.54, 2.75, 3.24, 2.62, 2.89,
    3.31, 2.75, 3.44, 1.59, 2.26, 2.33, 2.62, 2.89
]

# Create an instance of RollingAveragePlotter
plotter = RollingAveragePlotter(data_points)
plotter.plot()
