#! /usr/bin/env sh

options=$(getopt -o qlfcud -l query,list,fetch,create,update,delete -- "$@")

[ $? -eq 0 ] || {
  echo "Incorrect options provided"
  exit 1
}

# shell colors
COLOR_TEST=$(tput setaf 3)
COLOR_RESET=$(tput sgr0)

# by default
log_query="false"
action="all"

eval set -- "$options"
while true; do
  case "$1" in
    -q|--log-query) log_query="true"   ;;
    -l|--list)   action="list"   ;;
    -f|--fetch)  action="fetch"  ;;
    -c|--create) action="create" ;;
    -u|--update) action="update" ;;
    -d|--delete) action="delete" ;;
    # *)           action="all"    ;;
    --) shift; break ;;
  esac
  shift
done

list () {
  echo "${COLOR_TEST}\nLIST TEST\n---------${COLOR_RESET}"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -X GET http://localhost:8011/strategies"
  fi
  curl -i -X GET http://localhost:8011/strategies
}

fetch () {
  echo "${COLOR_TEST}\nFETCH TEST\n----------${COLOR_RESET}"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: No query to display yet !"
  fi
  echo "WARNING: The \"fetch\" test has not yet been defined !"
}

create () {
  echo "${COLOR_TEST}\nCREATE TEST\n-----------${COLOR_RESET}"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -H \"Content-Type: application/json\" -d '{\"description\":\"faire une bonne sieste\", \"needId\":3}' -X POST http://localhost:8011/strategy"
  fi
  curl -i -H "Content-Type: application/json" -d '{"description":"faire une bonne sieste", "needId":3}' -X POST http://localhost:8011/strategy
}

update () {
  echo "${COLOR_TEST}\nUPDATE TEST\n-----------${COLOR_RESET}"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: No query to display yet !"
  fi
  echo "WARNING: The \"update\" test has not yet been defined !"
}

delete () {
  echo "${COLOR_TEST}\nDELETE TEST\n-----------${COLOR_RESET}"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: No query to display yet !"
  fi
  echo "WARNING: The \"delete\" test has not yet been defined !"
}

if [ "${action}" = "list" ]; then
  list
elif [ "${action}" = "fetch" ]; then
  fetch
elif [ "${action}" = "create" ]; then
  create
elif [ "${action}" = "update" ]; then
  update
elif [ "${action}" = "delete" ]; then
  delete
else
  list
  fetch
  create
  update
  delete
fi

exit 0
