#!/bin/sh

# Run the crawler application
/web-crawler &

# Define cleanup procedure
cleanup() {
    echo "Container stopped, performing cleanup..."
    
    # Copy the results.json file to the output directory if it exists
    if [ -f /scraped-data/results.json ]; then
        cp /scraped-data/results.json /output/results.json
        echo "results.json copied to /output directory."
    else
        echo "No results.json file found in /scraped-data."
    fi
}

# Trap SIGTERM and SIGINT signals and call cleanup
trap 'cleanup' SIGTERM SIGINT

# Wait for the application to finish
wait $!
