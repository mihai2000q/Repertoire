DIRECTORY_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${DIRECTORY_SCRIPT}/.env

SQL_DIR="$DIRECTORY_SCRIPT/.sample-data"

file_count=$(ls $SQL_DIR/*.sql 2>/dev/null | wc -l)

for i in $(seq 1 $file_count)
do
  file=$(ls $SQL_DIR/$i.*.sql 2>/dev/null)
  filename=$(basename "$file")
  echo "Executing sql file: $filename"
  docker exec -u $DB_USER Repertoire-Database psql $DB_NAME $DB_USER -f ./.sample-data/$filename
done