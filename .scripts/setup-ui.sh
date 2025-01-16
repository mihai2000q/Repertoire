DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR_SCRIPT}/../repertoire.ui

echo "Installing UI Dependencies"
npm ci

echo -e "UI Setup Finished\n"