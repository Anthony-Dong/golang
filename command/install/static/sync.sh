#!/bin/bash

set -e

########################### 配置开始
# 远程服务器的用户名称
readonly DEV_USER="${BYTED_USERNAME}"
# 远程服务器的地址
readonly DEV_IP="${BYTED_HOST_IP}"
# 远程同步目录，推荐用户的根目录，例如可以执行 echo "$HOME" 查看用户根目录
readonly REMOTE_HOME="/home/${DEV_USER}"

# 本地脚本目录
readonly SYNC_HOME="${HOME}/go/bin/sync-devbox"
# 本地同步白名单
readonly WHITE_LIST=("${HOME}/go/src" "${HOME}/data")
########################### 配置结束
LOCAL_HOME="$HOME"

function usage() {
    echo "Usage: $(basename "$0") [dir]"
    echo ""
    echo "  Support reverse synchronization, support forward synchronization"
}

# 获取远程文件的文件类型
function get_remote_file_type() {
    local FILE="$1"
    FILE_TYPE=""
    if ssh ${DEV_USER}@${DEV_IP} [ -d "$FILE" ]; then
        FILE_TYPE="dir"
    elif ssh ${DEV_USER}@${DEV_IP} [ -f "$FILE" ]; then
        FILE_TYPE="file"
    else
        echo "> [ERROR] the remote file or directory does not exist: $FILE" >&2
        return 1
    fi
    echo "$FILE_TYPE"
}

# 获取本地文件的文件类型
function get_file_type() {
    local FILE="$1"
    FILE_TYPE=""
    if [ -d "${FILE}" ]; then
        FILE_TYPE="dir"
    elif [ -f "${FILE}" ]; then
        FILE_TYPE="file"
    else
        echo "> [ERROR] the file or directory does not exist: ${FILE}" >&2
        return 1
    fi
    echo "${FILE_TYPE}"
}

# 检测是否允许传输
function is_white_list_file() {
    local FILE="$1"
    NOT_PASS_DIR=()
    for elem in "${WHITE_LIST[@]}"; do
        if [[ "${FILE}" == ${elem}/* ]]; then
            break
        fi
        NOT_PASS_DIR+=("$elem")
    done
    if [ ${#WHITE_LIST[@]} -eq ${#NOT_PASS_DIR[@]} ]; then
        echo "> [ERROR] the current directory ${FILE} is not under the whitelist: ${NOT_PASS_DIR[*]}" >&2
        return 1
    fi
}

# 正向同步 local->remote
function file_rsync() {
    FILE=$(realpath "$1")
    FILE_TYPE=$(get_file_type "$FILE")

    echo "> file: $FILE"
    echo "> file type: ${FILE_TYPE}" # 目录同步整个目录，文件同步单个文件

    is_white_list_file "$FILE"

    REMOTE_FILE=${FILE/${LOCAL_HOME}/${REMOTE_HOME}}
    ssh ${DEV_USER}@${DEV_IP} mkdir -p "$(dirname "$REMOTE_FILE")"
    echo "> mkdir -p $(dirname "$REMOTE_FILE")"

    # 这里选择下文件目录
    IGNORE_FILE=${SYNC_HOME}/.fileignore
    if [ "$FILE_TYPE" = "dir" ] && [ -f "${FILE}/.fileignore" ]; then
        IGNORE_FILE="${FILE}/.fileignore"
    fi

    # 如果同步的是文件，则忽略排除文件
    rsync_opt=()
    if [ "$FILE_TYPE" = "dir" ] && [ -f "${IGNORE_FILE}" ]; then
        rsync_opt+=(--exclude-from="${IGNORE_FILE}")
    fi

    # 目录需要特殊处理下
    if [ "$FILE_TYPE" = "dir" ]; then
        FILE="$FILE/"
        REMOTE_FILE="$REMOTE_FILE/"
    fi

    # 用法参考: https://www.ruanyifeng.com/blog/2020/08/rsync.html
    set -x
    rsync -avz \
        --delete \
        --progress \
        --log-file="${SYNC_HOME}/sync-devbox.log" \
        --log-file-format="%t %f %b" \
        "${rsync_opt[@]}" \
        "${FILE}" "${DEV_USER}@${DEV_IP}:${REMOTE_FILE}"
    set +x
    echo "> rsync ${FILE} -> ${REMOTE_FILE}  success!"
}

# 反向同步 remote->local
function file_reverse_rsync() {
    FILE="$1"
    echo "> remote file: $FILE"
    LOCAL_FILE=${FILE/${REMOTE_HOME}/${LOCAL_HOME}}
    echo "> local file: $LOCAL_FILE"
    is_white_list_file "$LOCAL_FILE"

    FILE_TYPE=$(get_remote_file_type "$FILE")
    echo "> remote file type: $FILE_TYPE"

    # 如果允许反向同步目录，那么注释掉下面这行代码
    if [ "$FILE_TYPE" = "dir" ]; then
        echo "> [ERROR] reverse synchronization only supports synchronization of file types" >&2
        return 1
    fi

    rsync_opt=()
    if [ -f "${SYNC_HOME}/.fileignore" ]; then
        rsync_opt+=(--exclude-from="${SYNC_HOME}/.fileignore")
    fi

    mkdir -p "$(dirname "$LOCAL_FILE")"
    set -x
    rsync -avz \
        --delete \
        --progress \
        --log-file="${SYNC_HOME}/sync-devbox.log" \
        --log-file-format="%t %f %b" \
        "${rsync_opt[@]}" \
        "${DEV_USER}@${DEV_IP}:${FILE}" "${LOCAL_FILE}"
    set +x
    echo "> rsync ${FILE} -> ${LOCAL_FILE}  success!"
}

FILE="$1"
if [ -z "$FILE" ]; then
    FILE=$(pwd)
fi

# usage
if [ "$1" = '--help' ] || [ "$1" = '-help' ] || [ "$1" = 'help' ] || [ "$1" = '-h' ]; then
    usage
    exit 1
fi

# 初始化sync脚本路径
if [ ! -d "$SYNC_HOME" ]; then mkdir -p "$SYNC_HOME"; fi

# 如果在远程目录下那么就表示反向同步
if [[ $FILE == ${REMOTE_HOME}/* ]]; then
    REVERSE_RSYNC="true"
    echo "> reverse_rsync: $REVERSE_RSYNC"
fi

# 执行
if [ "$REVERSE_RSYNC" = "true" ]; then
    file_reverse_rsync "$FILE"
else
    file_rsync "$FILE"
fi