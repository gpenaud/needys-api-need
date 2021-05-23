#! /usr/bin/env bash

options=$(getopt -o X:H:d: -l method:,header,data -- "$@")

[ $? -eq 0 ] || {
  echo "Incorrect options provided"
  exit 1
}

eval set -- "$options"
while true; do
  case "$1" in
    -X|--method) shift; method=${1} ;;
    -H|--header) shift; header=${1} ;;
    -d|--data)   shift; data=${1} ;;
    --) shift; break ;;
  esac
  shift
done

ss -lntp | grep -i  -E '*.8010' | gawk '{
  while (match($0, /pid=([0-9]+)/, ary)) {
    print 'killing', ary[1];
    system("sudo kill -9 " ary[1]);
    $0 = substr($0, RSTART + RLENGTH);
  }
}'

# start service and wait
go run main.go &
sleep 1

# register curl query
query="curl -X ${method} -H \"${header}\" -d '${data}' ${@}"

# execute curl query
eval ${query}
