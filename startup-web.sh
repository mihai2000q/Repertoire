WEB_DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${WEB_DIR_SCRIPT}/.scripts/base-startup.sh

# Setup Web
cd ${WEB_DIR_SCRIPT}/repertoire.web

echo "Settting up the development Environment Variables for the Web"
cp .env.dev .env

echo "Installing Web Dependencies"
npm ci

echo -e "Web Setup Finished\n"
echo "Setup Finished!"