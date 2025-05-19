#!/bin/bash

# seed_db.sh: Script to populate Book Tracker API with mock books.
# Usage: ./scripts/seed_db.sh
# Requires: curl, jq, API running at http://localhost:8080

BASE_URL="http://localhost:8080/api/v1"

add_book() {
    local title="$1"
    local author="$2"
    local status="$3"
    echo "Adding book: $title by $author (status: $status)"
    json_payload=$(printf '{"title":"%s","author":"%s","status":"%s"}' "$title" "$author" "$status")
    echo "Raw response:"
    response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/books" \
         -H "Content-Type: application/json" \
         -d "$json_payload")
    echo "$response"
    echo "$response" | grep -v "HTTP_STATUS" | jq . || echo "Failed to parse JSON response"
    echo ""
}

get_stats() {
    echo "Checking statistics:"
    curl -s "$BASE_URL/stats" | jq .
    echo ""
}

list_books() {
    echo "Listing all books":
    curl -s "$BASE_URL/books" | jq .
    echo ""
}

if ! command -v curl &> /dev/null; then
    echo "Error: curl is required but not installed."
    exit 1
fi
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed."
    exit 1
fi

echo "No need to clear database; it will be created fresh."

AUTHORS=("Jane Austen" "Herman Melville" "Alan Donovan" "George Orwell" "J.K. Rowling" "Ernest Hemingway" "Toni Morrison" "F. Scott Fitzgerald" "Virginia Woolf" "Mark Twain")
STATUSES=("unread" "reading" "complete")

for i in {1..50}; do
    AUTHOR=${AUTHORS[$((RANDOM % ${#AUTHORS[@]}))]}
    STATUS=${STATUSES[$((RANDOM % ${#STATUSES[@]}))]}

    TITLE="Book Title $i"

    add_book "$TITLE" "$AUTHOR" "$STATUS"
done

get_stats

list_books