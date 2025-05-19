#!/bin/bash

# test_metrics.sh
# Tests the Prometheus /metrics endpoint using curl

set -e

METRICS_URL="http://localhost:8080/metrics"

echo "Testing Prometheus metrics endpoint at $METRICS_URL..."

response=$(curl -s -o response.txt -w "%{http_code}" "$METRICS_URL")

if [ "$response" != "200" ]; then
    echo "Error: Expected status 200, got $response"
    cat response.txt
    rm -f response.txt
    exit 1
fi

if ! grep -q "http_requests_total" response.txt; then
    echo "Error: Response does not contain 'http_requests_total' metric"
    cat response.txt
    rm -f response.txt
    exit 1
fi

echo "Success: /metrics endpoint returned 200 and contains expected metrics"
rm -f response.txt