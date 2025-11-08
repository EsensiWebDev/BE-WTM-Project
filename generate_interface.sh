#!/bin/sh

DOMAIN_DIR="internal/domain"

# Function untuk ambil inisial receiver
get_initial() {
  echo "$1" | grep -o '[A-Z][a-z]*' | awk '{ printf "%s", tolower(substr($0,1,1)) }'
}

# Scan semua file .go di domain
find "$DOMAIN_DIR" -type f -name "*.go" | while read FILE; do

  # Temukan definisi interface Usecase dan Repository
  grep -E 'type [A-Za-z0-9]+(Usecase|Repository) interface' "$FILE" | while read LINE; do
    INTERFACE_NAME=$(echo "$LINE" | sed -E 's/type ([A-Za-z0-9]+) interface.*/\1/')
    INTERFACE_LOWER=$(echo "$INTERFACE_NAME" | sed -E 's/([A-Z]+)/_\L\1/g' | sed -E 's/^_//')
    INITIAL=$(get_initial "$INTERFACE_NAME")

    # Tentukan folder target
    if echo "$INTERFACE_NAME" | grep -qi "Usecase"; then
      TARGET_FOLDER="internal/usecase/${INTERFACE_LOWER}/"
    else
      TARGET_FOLDER="internal/repository/${INTERFACE_LOWER}/"
    fi

    mkdir -p "$TARGET_FOLDER"
    PACKAGE_NAME=$(basename "$TARGET_FOLDER")

    # Ambil semua method dari interface
    awk "/type $INTERFACE_NAME interface/,/}/" "$FILE" | \
      grep '^[[:space:]]*[A-Z][A-Za-z0-9_]*(' | \
      sed -E 's/^[[:space:]]*//' | while read signature; do

        METHOD_NAME=$(echo "$signature" | sed -E 's/\(.*//')
        METHOD_FILE=$(echo "$METHOD_NAME" | sed -E 's/([A-Z]+)/_\L\1/g' | sed -E 's/^_//')
        FILENAME="${TARGET_FOLDER}${METHOD_FILE}.go"

        if [ ! -f "$FILENAME" ]; then
          echo "✅ Membuat: $FILENAME"
          echo "package $PACKAGE_NAME" > "$FILENAME"
          echo "" >> "$FILENAME"
          echo "func (${INITIAL} *$INTERFACE_NAME) $signature {" >> "$FILENAME"
          echo "    // TODO: implement" >> "$FILENAME"
          echo "}" >> "$FILENAME"
        else
          echo "⏩ File sudah ada: $FILENAME (diabaikan)"
        fi
    done
  done
done
