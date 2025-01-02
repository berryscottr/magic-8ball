#!/bin/sh

# Function to check if it's Wednesday 8 AM
check_exit_condition() {
  while true; do
  sleep 10
    DAY=$(date +%u)    # 1=Monday, 7=Sunday
    HOUR=$(date +%H)
    if [ "$DAY" -eq 3 ] && [ "$HOUR" -eq 8 ]; then
      echo "Exiting because it's Wednesday at 8 AM."
      return 0
    fi
    echo "Not Wednesday 8 AM yet. Current day: $DAY, hour: $HOUR. Checking again in an hour..."
    sleep 3600  # Check every hour
  done
}

# Start the application and keep it running until the condition is met
echo "Starting the application..."
./magic-8ball &  # Run the application in the background
APP_PID=$!  # Capture the PID of the application process

# Monitor for the exit condition
check_exit_condition

# Gracefully stop the application
echo "Stopping the application..."
kill $APP_PID  # Send a termination signal to the application process
wait $APP_PID  # Wait for the application process to exit
echo "Application stopped."
