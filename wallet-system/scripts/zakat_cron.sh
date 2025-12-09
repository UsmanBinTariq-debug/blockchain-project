#!/bin/bash

# Zakat Scheduler Script
# This script processes zakat deductions on the 1st of each month

# Set variables
BACKEND_URL="http://localhost:8080/api"
ADMIN_TOKEN="${JWT_TOKEN}"

# Function to process zakat for all wallets
process_zakat() {
    echo "Processing monthly zakat deductions..."
    
    # Call backend endpoint to process zakat
    curl -X POST "${BACKEND_URL}/zakat/process" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" \
        -H "Content-Type: application/json"
    
    echo "Zakat processing completed"
}

# Check if today is the 1st of the month
DAY=$(date +%d)

if [ "$DAY" = "01" ]; then
    process_zakat
else
    echo "Not the 1st of the month. Skipping zakat processing."
    echo "Current day: $DAY"
fi
