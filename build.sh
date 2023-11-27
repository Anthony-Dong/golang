#!/usr/bin/env bash

set -e


function cross_go_build(){
  # CGO_ENABLED=0 GOOS=windows GOARCH=amd64
  # CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
  # CGO_ENABLED=0 GOOS=darwin GOARCH=arm
  # CGO_ENABLED=0 GOOS=linux GOARCH=amd64
  binary="$3"
  if [ "$1" == "windows" ]; then
    binary="$3.exe"
  fi
  GO111MODULE=on CGO_ENABLED=1 GOOS=$1 GOARCH=$2  go build -v -ldflags "-s -w" -o "bin/$1_$2/$binary" "cli/$3/main.go"
}


function go_build(){
  binary="$1"
  GO111MODULE=on CGO_ENABLED=1 go build -v -ldflags "-s -w" -o "bin/$binary" "cli/$1/main.go"
}


function go_install(){
  binary="$1"
  GO111MODULE=on CGO_ENABLED=1 go build -v -ldflags "-s -w" -o "$(go env GOPATH)/bin/$binary" "cli/$1/main.go"
}

function format_golang_file () {
  project_dir=$(realpath "$1")
	# shellcheck disable=SC2044
	for elem in $(find "${project_dir}" -name '*.go'); do
#	  echo "format ${elem}"
		gofmt -w "${elem}"  > /dev/null 2>&1;
		goimports -w -srcdir "${project_dir}" -local "$2" "${elem}" > /dev/null 2>&1;
	done
}

case $1 in
build)
  go_build $2
  ;;
install)
  go_install $2
  ;;
format)
  format_golang_file . "github.com/anthony-dong/golang"
  ;;
*)
  ;;
esac