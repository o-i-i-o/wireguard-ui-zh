#!/usr/bin/env bash
set -e

DIR=$(dirname "$0")

# 安装 Node 模块
YARN=yarn
[ -x /usr/bin/lsb_release ] && [ -n "$(lsb_release -i | grep Debian)" ] && YARN=yarnpkg

# 修复点1：替换已弃用的 --pure-lockfile 和 --production
# 使用现代 Yarn 的 --immutable 和 --mode=production
$YARN install --immutable --mode=production

# 修复点2：确保存在工作区配置（兼容 yarn workspaces focus）
if ! grep -q '"workspaces":' "${DIR}/package.json"; then
  echo "添加工作区配置到 package.json"
  TEMP_FILE=$(mktemp)
  jq '. + {"workspaces": ["."]}' "${DIR}/package.json" > "$TEMP_FILE"
  mv "$TEMP_FILE" "${DIR}/package.json"
fi

# 复制 admin-lte 的 dist 文件
mkdir -p "${DIR}/assets/dist/js" "${DIR}/assets/dist/css" && \
  cp -r "${DIR}/node_modules/admin-lte/dist/js/adminlte.min.js" "${DIR}/assets/dist/js/adminlte.min.js" && \
  cp -r "${DIR}/node_modules/admin-lte/dist/css/adminlte.min.css" "${DIR}/assets/dist/css/adminlte.min.css"

# 复制自定义 JS 文件
cp -r "${DIR}/custom" "${DIR}/assets"

# 复制插件文件
mkdir -p "${DIR}/assets/plugins" && \
  cp -r "${DIR}/node_modules/admin-lte/plugins/jquery" \
  "${DIR}/node_modules/admin-lte/plugins/fontawesome-free" \
  "${DIR}/node_modules/admin-lte/plugins/bootstrap" \
  "${DIR}/node_modules/admin-lte/plugins/icheck-bootstrap" \
  "${DIR}/node_modules/admin-lte/plugins/toastr" \
  "${DIR}/node_modules/admin-lte/plugins/jquery-validation" \
  "${DIR}/node_modules/admin-lte/plugins/select2" \
  "${DIR}/node_modules/jquery-tags-input" \
  "${DIR}/assets/plugins/"
