DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR_SCRIPT}/../repertoire.storage

echo "Settting up the development Environment Variables for the Storage"
cp .env.dev .env

echo "Installing Storage Dependencies"
go mod download

echo -e "Storage Setup Finished\n"