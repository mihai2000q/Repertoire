BASE_DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${BASE_DIR_SCRIPT}/setup-server.sh
source ${BASE_DIR_SCRIPT}/setup-auth.sh
source ${BASE_DIR_SCRIPT}/setup-storage.sh
source ${BASE_DIR_SCRIPT}/setup-ui.sh