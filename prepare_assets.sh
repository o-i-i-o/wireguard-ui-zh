#!/usr/bin/env bash
set -e

# 获取脚本所在目录的路径
DIR=$(dirname "$0")

# 安装Node.js模块
YARN=yarn
# 检查系统是否为Debian系列，如果是则使用yarnpkg命令
[ -x /usr/bin/lsb_release ] && [ -n "`lsb_release -i | grep Debian`" ] && YARN=yarnpkg

# 使用workspaces focus命令替代install --production（Yarn 4.x推荐）
$YARN workspaces focus --production --immutable

# 动态查找admin-lte包在Yarn unplugged目录中的位置
# 该目录名包含版本号和哈希值，不同环境可能不同
ADMIN_LTE_UNPLUGGED_PATH=$(find "${DIR}/.yarn/unplugged" -type d -name "admin-lte-npm-*" -print -quit)

# 检查是否找到admin-lte目录
if [ -z "$ADMIN_LTE_UNPLUGGED_PATH" ]; then
  echo "错误：在.yarn/unplugged目录中未找到admin-lte包！"
  exit 1
fi

# 构建admin-lte的完整路径
ADMIN_LTE_PATH="${ADMIN_LTE_UNPLUGGED_PATH}/node_modules/admin-lte"

# 复制admin-lte的JavaScript和CSS文件到assets目录
if [ -d "${ADMIN_LTE_PATH}/dist" ]; then
  mkdir -p "${DIR}/assets/dist/js" "${DIR}/assets/dist/css" && \
  cp -r "${ADMIN_LTE_PATH}/dist/js/adminlte.min.js" "${DIR}/assets/dist/js/" && \
  cp -r "${ADMIN_LTE_PATH}/dist/css/adminlte.min.css" "${DIR}/assets/dist/css/"
else
  echo "错误：未找到admin-lte/dist目录！"
  exit 1
fi

# 复制admin-lte的插件到assets目录
if [ -d "${ADMIN_LTE_PATH}/plugins" ]; then
  mkdir -p "${DIR}/assets/plugins" && \
  cp -r "${ADMIN_LTE_PATH}/plugins/jquery" \
        "${ADMIN_LTE_PATH}/plugins/fontawesome-free" \
        "${ADMIN_LTE_PATH}/plugins/bootstrap" \
        "${ADMIN_LTE_PATH}/plugins/icheck-bootstrap" \
        "${ADMIN_LTE_PATH}/plugins/toastr" \
        "${ADMIN_LTE_PATH}/plugins/jquery-validation" \
        "${ADMIN_LTE_PATH}/plugins/select2" \
        "${DIR}/assets/plugins/"
else
  echo "错误：未找到admin-lte/plugins目录！"
  exit 1
fi

# 复制自定义JavaScript文件到assets目录
cp -r "${DIR}/custom" "${DIR}/assets"
