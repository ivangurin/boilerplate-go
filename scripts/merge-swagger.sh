#!/bin/bash

set -e

OUTPUT_FILE="pkg/swagger/swagger.json"
# Находим все swagger файлы кроме swagger.json и сортируем по имени
SWAGGER_FILES=($(find pkg/pb -name "*.swagger.json" ! -name "swagger.json" | sort))

# Проверяем наличие jq
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Install it with: brew install jq"
    exit 1
fi

# Проверяем наличие входных файлов
for file in "${SWAGGER_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo "Error: File $file not found"
        exit 1
    fi
done

# Создаём базовую структуру объединённого swagger
cat > "$OUTPUT_FILE" <<'EOF'
{
  "swagger": "2.0",
  "info": {
    "title": "Boilerplate API",
    "version": "1.0.0"
  },
  "schemes": ["http", "https"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "securityDefinitions": {
    "x-auth": {
      "type": "apiKey",
      "in": "header",
      "name": "authorization"
    }
  },
  "paths": {},
  "definitions": {}
}
EOF

# Объединяем paths из всех файлов
MERGED_PATHS=$(jq -s 'reduce .[] as $item ({}; . * ($item.paths // {}))' "${SWAGGER_FILES[@]}")
echo "$MERGED_PATHS" | jq '.' > /tmp/paths.json
jq --slurpfile paths /tmp/paths.json '.paths = $paths[0]' "$OUTPUT_FILE" > "${OUTPUT_FILE}.tmp"
mv "${OUTPUT_FILE}.tmp" "$OUTPUT_FILE"

# Объединяем definitions из всех файлов
MERGED_DEFS=$(jq -s 'reduce .[] as $item ({}; . * ($item.definitions // {}))' "${SWAGGER_FILES[@]}")
echo "$MERGED_DEFS" | jq '.' > /tmp/definitions.json
jq --slurpfile defs /tmp/definitions.json '.definitions = $defs[0]' "$OUTPUT_FILE" > "${OUTPUT_FILE}.tmp"
mv "${OUTPUT_FILE}.tmp" "$OUTPUT_FILE"

# Добавляем теги из всех файлов
TAGS=$(jq -s '[.[].tags // [] | .[]] | unique' "${SWAGGER_FILES[@]}")
jq --argjson tags "$TAGS" '.tags = $tags' "$OUTPUT_FILE" > "${OUTPUT_FILE}.tmp"
mv "${OUTPUT_FILE}.tmp" "$OUTPUT_FILE"

# Удаляем исходные файлы
for file in "${SWAGGER_FILES[@]}"; do
    rm -f "$file"
done

# Очищаем временные файлы
rm -f /tmp/paths.json /tmp/definitions.json

echo "Swagger files merged successfully into $OUTPUT_FILE"
