#!/bin/bash

output_file="code.txt"

# Remove existing code.txt if it exists
rm -f "$output_file"

# Find all .go files and process them
find . -type f -name "*.vue" | while read -r file; do
    # Write opening tag
    echo "<$file>" >> "$output_file"
    
    # Write file contents
    cat "$file" >> "$output_file"
    
    # Write closing tag
    echo "</$file>" >> "$output_file"
    
    # Add newline for readability
    echo "" >> "$output_file"
done