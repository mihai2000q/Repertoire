DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Setup Server
cd ${DIR_SCRIPT}/repertoire.server

echo "Settting up the development Environment Variables for the Server"
cp .env.dev .env

echo "Installing Server Dependencies"
go mod download

echo -e "Server Setup Finished\n"

cd .. 

# Setup UI
cd ${DIR_SCRIPT}/repertoire.ui

echo "Installing UI Dependencies"
npm ci

echo -e "UI Setup Finished\n"

cd ..

# Setup Web
cd ${DIR_SCRIPT}/repertoire.web

echo "Settting up the development Environment Variables for the Web"
cp .env.dev .env

echo "Installing Web Dependencies"
npm ci

echo -e "Web Setup Finished\n"
echo "Setup Finished!"