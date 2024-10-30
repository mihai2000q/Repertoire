DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR_SCRIPT}/../repertoire.server

echo "Settting up the development Environment Variables for the Server"
cp .env.dev .env

echo "Installing Server Dependencies"
go mod download

echo -e "Server Setup Finished\n"

cd .. 