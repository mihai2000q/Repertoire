DIRECTORY_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${DIRECTORY_SCRIPT}/.env

echo -e "\nMeili Search Migrations running..."

for d in ${DIRECTORY_SCRIPT}/migrations/search/*/ ; do
  go run ${d}
done

echo -e "Meili Search Migrations done!"