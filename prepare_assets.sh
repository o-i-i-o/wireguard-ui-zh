#!/usr/bin/env bash
set -e

DIR=$(dirname "$0")

# install node modules
YARN=yarn
[ -x /usr/bin/lsb_release ] && [ -n "`lsb_release -i | grep Debian`" ] && YARN=yarnpkg

# 使用 workspaces focus 替代 install --production
$YARN workspaces focus --production

# 动态查找admin-lte目录
ADMIN_LTE_UNPLUGGED_PATH=$(find "${DIR}/.yarn/unplugged" -type d -name "admin-lte" -print -quit)

if [ -z "$ADMIN_LTE_UNPLUGGED_PATH" ]; then
  echo "Error: admin-lte package not found in .yarn/unplugged directory!"
  exit 1
fi

ADMIN_LTE_PATH="${ADMIN_LTE_UNPLUGGED_PATH}/node_modules/admin-lte"

# 检查 dist 目录是否存在
if [ -d "${ADMIN_LTE_PATH}/dist" ]; then
  mkdir -p "${DIR}/assets/dist/js" "${DIR}/assets/dist/css" && \
  cp -r "${ADMIN_LTE_PATH}/dist/js/adminlte.min.js" "${DIR}/assets/dist/js/" && \
  cp -r "${ADMIN_LTE_PATH}/dist/css/adminlte.min.css" "${DIR}/assets/dist/css/"
else
  echo "Error: admin-lte/dist directory not found!"
  exit 1
fi

# 检查 plugins 目录是否存在
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
  echo "Error: admin-lte/plugins directory not found!"
  exit 1
fi

# Copy helper js
cp -r "${DIR}/custom" "${DIR}/assets"
