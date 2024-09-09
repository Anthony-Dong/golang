#!/usr/bin/env bash

set -e

function usage() {
    echo "Usage:"
    echo "  $(basename "$0") [-F req] [-p pid]"
    echo "  $(basename "$0") [-F req] [binary]"
    echo "注意: -F 默认是 /proc/sys/kernel/perf_event_max_sample_rate 的值, 默认: 100000"
    echo "注意: 必须Linux环境下执行"
}


if [ -z "$1" ] || [ "$(uname)" != "Linux" ]; then
    usage
    exit 1
fi

FlameGraph_home="${HOME}/go/src/github.com/brendangregg/FlameGraph"

if [ ! -d "$FlameGraph_home" ]; then
    mkdir -p "$(dirname "$FlameGraph_home")"
    git clone --depth 1 https://github.com/brendangregg/FlameGraph.git "$FlameGraph_home"
fi

output_file=$(date '+%Y%m%d%H%M%S')
process_output_perf="output/${output_file}.data"
process_output_svg="output/${output_file}.svg"
function parse_perf() {
    echo "process done...."
    perf script -i "$process_output_perf" | "$FlameGraph_home"/stackcollapse-perf.pl | "$FlameGraph_home"/flamegraph.pl >"${process_output_svg}"
    echo "open ${process_output_svg}"
}

# 移除文件
rm -rf "$process_output_perf"
rm -rf "$process_output_svg"

mkdir -p output

# 结束的时候
trap parse_perf EXIT

set -x

perf record -o "$process_output_perf" -g "$@"