#!/usr/bin/env bash

set -e

export GO111MODULE=on

build_flag=("-v" "-ldflags" "-s -w")

# cross_go_build windows amd64
function cross_go_build(){
  binary="$3"
  if [ "$1" == "windows" ]; then
    binary="$3.exe"
  fi
  CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build "${build_flag[@]}" -o "bin/$1_$2/$binary" "cli/$3/main.go"
}

# go_build xxx
function go_build(){
  binary="$1"
  if [ "$IS_SUBMOD" = "1" ]; then
    cd "cli/$1" && go build "${build_flag[@]}" -o "../../bin/$binary" "main.go" && cd -
  else
    go build "${build_flag[@]}" -o "bin/$binary" "cli/$1/main.go"
  fi
}

function go_install(){
  binary="$1"
  output="$(go env GOPATH)/bin/$binary"
  if [ "$IS_SUBMOD" = "1" ]; then
    cd "cli/$1" && go build "${build_flag[@]}" -o "$output" "main.go" && cd -
  else
    go build "${build_flag[@]}" -o "$output" "cli/$1/main.go"
  fi
}

function format_golang_file () {
  project_dir=$(realpath "$1")
	# shellcheck disable=SC2044
	for elem in $(find "${project_dir}" -name '*.go' | grep -v 'example/'); do
		gofmt -w "${elem}" 2>&1;
		goimports -w -srcdir "${project_dir}" -local "$2" "${elem}" 2>&1;
	done
}

case $1 in
build)
  go_build "$2"
  ;;
install)
  go_install "$2"
  ;;
cors)
  cross_go_build windows amd64 "$2"
  cross_go_build darwin amd64 "$2"
  cross_go_build linux amd64 "$2"
  cross_go_build darwin arm64 "$2"
  ;;
format)
  format_golang_file . "github.com/anthony-dong/golang"
  ;;
*)
  ;;
esac