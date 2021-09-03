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
  printf "${COLOR_TEST}\nLIST TEST\n---------${COLOR_RESET}\n"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -X GET http://localhost:8010/needs"
  fi
  curl -i -X GET http://localhost:8010/needs && printf "\n"
}

fetch () {
  printf "${COLOR_TEST}\nFETCH TEST\n----------${COLOR_RESET}\n"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -X GET http://localhost:8010/need/2"
  fi
  curl -i -X GET http://localhost:8010/need/2 && printf "\n"
}

create () {
  printf "${COLOR_TEST}\nCREATE TEST\n-----------${COLOR_RESET}\n"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -H \"Content-Type: application/json\" -d '{\"name\":\"sex\", \"priority\":\"very high\"}' -X POST http://localhost:8010/need"
  fi
  curl -i -H "Content-Type: application/json" -d '{"name":"sex", "priority":"very high"}' -X POST http://localhost:8010/need && printf "\n"
}

update () {
  printf "${COLOR_TEST}\nUPDATE TEST\n-----------${COLOR_RESET}\n"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -X POST -d '{\"Name\":\"rest\",\"Priority\":\"huge\"}' http://localhost:8010/need"
  fi
  curl -i -X POST -d '{"Name":"rest","Priority":"huge"}' http://localhost:8010/need && printf "\n"
}

delete () {
  printf "${COLOR_TEST}\nDELETE TEST\n-----------${COLOR_RESET}\n"
  if [ "${log_query}" = "true" ]; then
    echo "DEBUG: curl -i -X DELETE http://localhost:8010/need/sex"
  fi
  curl -i -X DELETE http://localhost:8010/need/sex && printf "\n"
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
