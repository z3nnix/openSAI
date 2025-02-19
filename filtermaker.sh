#!/bin/sh

FILE="vocabulary.bot"

if [ ! -f "$FILE" ]; then
  echo "Файл '$FILE' не существует."
  exit 1
fi

if [ -n "$1" ]; then
  SYMBOL="$1"
  
  printf "Are you sure you want to delete the symbol '$SYMBOL'? This will remove it in ALL phrases.(Y/n) "
  read REPLY
  
  if [ "$REPLY" != "y" ] && [ "$REPLY" != "Y" ]; then
    echo "Canceled."
    exit 0
  fi

  TEMP_FILE=$(mktemp)
  sed "s/$SYMBOL//g" "$FILE" | sort | uniq > "$TEMP_FILE"
else
  TEMP_FILE=$(mktemp)
  
  sort "$FILE" | uniq > "$TEMP_FILE"
fi

mv "$TEMP_FILE" "$FILE"

echo "Done."

exit 0
