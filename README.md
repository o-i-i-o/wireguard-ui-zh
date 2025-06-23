![](https://github.com/ngoduykhanh/wireguard-ui/workflows/wireguard-ui%20build%20release/badge.svg)

### wireguard-ui
一款用于管理 WireGuard 设置的 Web 用户界面。

### 功能特点
- **友好的用户界面**
- **身份验证功能**
- **管理额外的客户端信息**（如姓名、电子邮件等）
- **支持多种方式获取客户端配置**：通过二维码、文件、电子邮件或 Telegram 等途径。

![wireguard-ui 0.3.7](https://user-images.githubusercontent.com/37958026/177041280-e3e7ca16-d4cf-4e95-9920-68af15e780dd.png)

### 运行 WireGuard - UI
> ⚠️默认用户名和密码均为`admin`。为确保设置安全，请务必进行修改。

#### 使用二进制文件
从发布页面下载二进制文件，然后直接在主机上运行：

```
./wireguard-ui
```


#### 使用 Docker Compose
[examples/docker-compose](examples/docker-compose)文件夹中提供了示例 Docker Compose 文件。选择最适合您需求的示例，根据需要调整配置，然后按如下方式运行：

```
docker-compose up
```


### 环境变量
|变量                            |描述|默认值|
|-------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------
|`BASE_PATH`                    |                                              若在反向代理虚拟主机的子路径下运行 wireguard - ui，可设置此变量（例如 /wireguard）|无|
|`BIND_ADDRESS`                 |可访问 Web 界面的地址和端口，使用 unix:///abspath/to/file.socket 表示 Unix 域套接字|0.0.0.0:80|
|`SESSION_SECRET`               |用于加密会话 cookie 的密钥，应设置为随机值|无|
|`SESSION_SECRET_FILE`          |可选的会话密钥文件路径，若`SESSION_SECRET`留空则此设置生效|无|
|`SESSION_MAX_DURATION`         |已记住会话的最大有效天数（以天为单位）。无论此设置如何，未刷新的会话最长有效期为 7 天|90|
|`SUBNET_RANGES`                |地址细分范围列表，格式为：`SR Name:10.0.1.0/24; SR2:10.0.2.0/24,10.0.3.0/24`，每个 CIDR 必须位于服务器接口之一的范围内|无|
|`WGUI_USERNAME`                |登录页面的用户名，仅用于数据库初始化|`admin`|
|`WGUI_PASSWORD`                |登录页面用户的密码，将自动进行哈希处理，仅用于数据库初始化|`admin`|
|`WGUI_PASSWORD_FILE`           |可选的用户登录密码文件路径，将自动进行哈希处理，仅用于数据库初始化。若`WGUI_PASSWORD`留空则此设置生效|无|
|`WGUI_PASSWORD_HASH`           |登录页面用户的密码哈希值（替代`WGUI_PASSWORD`），仅用于数据库初始化|无|
|`WGUI_PASSWORD_HASH_FILE`      |可选的用户登录密码哈希值文件路径（替代`WGUI_PASSWORD_FILE`），仅用于数据库初始化。若`WGUI_PASSWORD_HASH`留空则此设置生效|无|
|`WGUI_ENDPOINT_ADDRESS`        |全局设置中客户端应连接的默认端点地址，该端点也可包含端口，当在内部使用`WGUI_SERVER_LISTEN_PORT`端口监听但通过另一个端口（如 9000）转发时非常有用。例如：myvpn.dyndns.com:9000|解析为您的公网 IP 地址|
|`WGUI_FAVICON_FILE_PATH`       |用作网站图标的文件路径|嵌入式 WireGuard 标志|
|`WGUI_DNS`                     |全局设置中使用的默认 DNS 服务器（逗号分隔列表）|`1.1.1.1`|
|`WGUI_MTU`                     |全局设置中使用的默认最大传输单元|`1450`|
|`WGUI_PERSISTENT_KEEPALIVE`    |全局设置中 WireGuard 的默认持久保持活动时间|`15`|
|`WGUI_FIREWALL_MARK`           |默认的 WireGuard 防火墙标记|`0xca6c`（51820）|
|`WGUI_TABLE`                   |默认的 WireGuard 表值设置|`auto`|
|`WGUI_CONFIG_FILE_PATH`        |全局设置中使用的默认 WireGuard 配置文件路径|`/etc/wireguard/wg0.conf`|
|`WGUI_LOG_LEVEL`               |默认日志级别，可选值：`DEBUG`、`INFO`、`WARN`、`ERROR`、`OFF`|`INFO`|
|`WG_CONF_TEMPLATE`             |自定义的`wg.conf`配置文件模板，请参考我们的[默认模板](https://github.com/o-i-i-o/wireguard-ui-zh/blob/master/templates/wg.conf)|无|
|`EMAIL_FROM_ADDRESS`           |发件人电子邮件地址|无|
|`EMAIL_FROM_NAME`              |发件人姓名|`WireGuard UI`|
|`SENDGRID_API_KEY`             |SendGrid API 密钥|无|
|`SENDGRID_API_KEY_FILE`        |可选的 SendGrid API 密钥文件路径，若`SENDGRID_API_KEY`留空则此设置生效|无|
|`SMTP_HOSTNAME`                |SMTP IP 地址或主机名|`127.0.0.1`|
|`SMTP_PORT`                    |SMTP 端口|`25`|
|`SMTP_USERNAME`                |SMTP 用户名|无|
|`SMTP_PASSWORD`                |SMTP 用户密码|无|
|`SMTP_PASSWORD_FILE`           |可选的 SMTP 用户密码文件路径，若`SMTP_PASSWORD`留空则此设置生效|无|
|`SMTP_AUTH_TYPE`               |SMTP 身份验证类型，可能的值：`PLAIN`、`LOGIN`、`NONE`|`NONE`|
|`SMTP_ENCRYPTION`              |加密方法，可能的值：`NONE`、`SSL`、`SSLTLS`、`TLS`、`STARTTLS`|`STARTTLS`|
|`SMTP_HELO`                    |用于 HELO 消息的主机名，smtp - relay.gmail.com 需要将此设置为除`localhost`之外的任何值|`localhost`|
|`TELEGRAM_TOKEN`               |用于向客户端分发配置的 Telegram 机器人令牌|无|
|`TELEGRAM_ALLOW_CONF_REQUEST`  |允许用户通过发送消息从机器人获取配置|`false`|
|`TELEGRAM_FLOOD_WAIT`          |处理下一个配置请求之前的等待时间（以分钟为单位）|`60`|

#### 服务器配置默认值
这些环境变量用于控制初始化数据库时使用的默认服务器设置。

|变量|描述|默认值|
|-----------------------------------|-----------------------------------------------------------------------------------------------|-----------------|
|`WGUI_SERVER_INTERFACE_ADDRESSES`  |WireGuard 服务器配置的默认接口地址（逗号分隔列表）|`10.252.1.0/24`|
|`WGUI_SERVER_LISTEN_PORT`          |默认的服务器监听端口|`51820`|
|`WGUI_SERVER_POST_UP_SCRIPT`       |默认的服务器启动后脚本|无|
|`WGUI_SERVER_POST_DOWN_SCRIPT`     |默认的服务器关闭后脚本|无|

#### 新客户端默认值
这些环境变量用于设置`新建客户端`对话框中的默认值。

|变量|描述|默认值|
|---------------------------------------------|-------------------------------------------------------------------------------------------------|-------------|
|`WGUI_DEFAULT_CLIENT_ALLOWED_IPS`            |`允许的 IP`字段的 CIDR 逗号分隔列表（默认）|`0.0.0.0/0`|
|`WGUI_DEFAULT_CLIENT_EXTRA_ALLOWED_IPS`      |`额外允许的 IP`字段的 CIDR 逗号分隔列表（默认为空）|无|
|`WGUI_DEFAULT_CLIENT_USE_SERVER_DNS`         |布尔值[`0`、`f`、`F`、`false`、`False`、`FALSE`、`1`、`t`、`T`、`true`、`True`、`TRUE`]|`true`|
|`WGUI_DEFAULT_CLIENT_ENABLE_AFTER_CREATION`  |布尔值[`0`、`f`、`F`、`false`、`False`、`FALSE`、`1`、`t`、`T`、`true`、`True`、`TRUE`]|`true`|

#### 仅适用于 Docker
这些环境变量仅适用于 Docker 容器。

|变量|描述|默认值|
|-----------------------|---------------------------------------------------------------|---------|
|`WGUI_MANAGE_START`    |在容器启动/停止时启动/停止 WireGuard|`false`|
|`WGUI_MANAGE_RESTART`  |在 UI 中应用配置更改时自动重启 WireGuard|`false`|

### 自动重启 WireGuard 守护进程
WireGuard - UI 仅负责配置生成。您可以使用 systemd 来监视更改并重启服务，以下是一个示例：

#### 使用 systemd
创建`/etc/systemd/system/wgui.service`
```bash
cd /etc/systemd/system/
cat << EOF > wgui.service
[Unit]
Description=Restart WireGuard
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl restart wg-quick@wg0.service

[Install]
RequiredBy=wgui.path
EOF
```

创建 `/etc/systemd/system/wgui.path`

```bash
cd /etc/systemd/system/
cat << EOF > wgui.path
[Unit]
Description=Watch /etc/wireguard/wg0.conf for changes

[Path]
PathModified=/etc/wireguard/wg0.conf

[Install]
WantedBy=multi-user.target
EOF
```

应用服务

```sh
systemctl enable wgui.{path,service}
systemctl start wgui.{path,service}
```

### 使用openrc

创建 `/usr/local/bin/wgui` 文件并赋予执行权限

```sh
cd /usr/local/bin/
cat << EOF > wgui
#!/bin/sh
wg-quick down wg0
wg-quick up wg0
EOF
chmod +x wgui
```

创建 `/etc/init.d/wgui` 文件并赋予执行权限

```sh
cd /etc/init.d/
cat << EOF > wgui
#!/sbin/openrc-run

command=/sbin/inotifyd
command_args="/usr/local/bin/wgui /etc/wireguard/wg0.conf:w"
pidfile=/run/${RC_SVCNAME}.pid
command_background=yes
EOF
chmod +x wgui
```

应用设置

```sh
rc-service wgui start
rc-update add wgui default
```

### 使用 Docker

将WGUI_MANAGE_RESTART设置为true可管理 Wireguard 接口的重启。使用WGUI_MANAGE_START=true还可以替代wg-quick@wg0服务的功能，通过将容器运行时设置为restart: unless-stopped来在启动时启动 Wireguard。这些设置还可以在容器重启后获取对 Wireguard 配置文件路径的更改。请确保在容器配置中包含--cap-add=NET_ADMIN以使此功能正常工作。

## 构建

### 构建 docker 镜像

进入项目根目录并运行以下命令：

```sh
docker build --build-arg=GIT_COMMIT=$(git rev-parse --short HEAD) -t wireguard-ui .
```

或者

```sh
docker compose build --build-arg=GIT_COMMIT=$(git rev-parse --short HEAD)
```

:information_source: （我不会用docker，所以没有我构建的镜像，原版镜像是未翻译的）您可以在Docker Hub上获取容器镜像，可通过以下命令拉取并使用： [Docker Hub](https://hub.docker.com/r/ngoduykhanh/wireguard-ui)

```
docker pull ngoduykhanh/wireguard-ui
````

### 编译二进制文件

提前生成 assets 资源
```sh
./prepare_assets.sh
```

然后构建可执行文件

```sh
go build -o wireguard-ui
```

## 许可证

许可证
MIT 许可证，详情请见LICENSE。 [LICENSE](https://github.com/ngoduykhanh/wireguard-ui/blob/master/LICENSE).


