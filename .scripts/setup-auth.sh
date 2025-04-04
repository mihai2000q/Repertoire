DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR_SCRIPT}/../repertoire.auth

echo "Settting up the development Environment Variables for the Auth"
cp .env.dev .env

echo "Installing Auth Dependencies"
go mod download

echo -e "Auth Setup Finished\n"