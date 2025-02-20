#!/bin/sh

FILE="vocabulary.bot"

# Check if the file exists
if [ ! -f "$FILE" ]; then
  echo "The file vocabulary.bot does not exist."
  exit 1
fi

TEMP_FILE=$(mktemp)

if [ -n "$1" ]; then
  SYMBOL="$1"
  
  # Confirm deletion of lines containing the specified symbol(s)
  printf "Are you sure you want to delete lines containing the symbol(s) '$SYMBOL', https://, http://, and @? (Y/n) "
  read REPLY
  
  if [ "$REPLY" != "y" ] && [ "$REPLY" != "Y" ]; then
    echo "Operation canceled."
    exit 0
  fi

  # Remove lines containing the specified symbol(s), https://, http://, and @, and sort them
  grep -vE "$SYMBOL|https://|http://|@" "$FILE" | sort | uniq > "$TEMP_FILE"
else
  # Remove lines containing https://, http://, and @, and sort them
  grep -vE "https://|http://|@" "$FILE" | sort | uniq > "$TEMP_FILE"
fi

# Replace the original file with the updated content
mv "$TEMP_FILE" "$FILE"

echo "Done."

exit 0
