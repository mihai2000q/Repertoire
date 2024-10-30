DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${DIR_SCRIPT}/.scripts/base-startup.sh

# Setup Web
cd ${DIR_SCRIPT}/repertoire.desktop

echo "Settting up the development Environment Variables for the Desktop"
cp .env.dev .env

echo "Installing Desktop Dependencies"
npm ci

echo -e "Desktop Setup Finished\n"
echo "Setup Finished!"