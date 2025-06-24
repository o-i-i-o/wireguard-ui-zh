package main

import (
	"crypto/sha512"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ngoduykhanh/wireguard-ui/store"
	"github.com/ngoduykhanh/wireguard-ui/telegram"

	"github.com/ngoduykhanh/wireguard-ui/emailer"
	"github.com/ngoduykhanh/wireguard-ui/handler"
	"github.com/ngoduykhanh/wireguard-ui/router"
	"github.com/ngoduykhanh/wireguard-ui/store/jsondb"
	"github.com/ngoduykhanh/wireguard-ui/util"
)

var (
	// 命令行横幅信息
	appVersion = "开发版本"
	gitCommit  = "N/A"
	gitRef     = "N/A"
	buildTime  = getCSTBuildTime() // 使用函数获取CST时间
	// 配置变量
	flagDisableLogin             = false     // 禁用应用程序的身份验证，这可能存在安全风险
	flagBindAddress              = "0.0.0.0:5000" // 应用程序绑定的地址和端口
	flagSmtpHostname             = "127.0.0.1"    // SMTP服务器主机名
	flagSmtpPort                 = 25            // SMTP服务器端口
	flagSmtpUsername             string          // SMTP服务器用户名
	flagSmtpPassword             string          // SMTP服务器密码
	flagSmtpAuthType             = "NONE"        // SMTP认证类型
	flagSmtpNoTLSCheck           = false         // 禁用SMTP的TLS验证，这可能存在安全风险
	flagSmtpEncryption           = "STARTTLS"    // SMTP加密方式
	flagSmtpHelo                 = "localhost"   // SMTP HELO主机名
	flagSendgridApiKey           string          // SendGrid API密钥
	flagEmailFrom                string          // 发件人邮箱地址
	flagEmailFromName            = "WireGuard UI" // 发件人邮箱名称
	flagTelegramToken            string          // 用于向客户端分发配置的Telegram机器人令牌
	flagTelegramAllowConfRequest = false         // 允许用户通过发送消息从机器人获取配置
	flagTelegramFloodWait        = 60            // 处理下一个配置请求前的等待时间（分钟）
	flagSessionSecret            = util.RandomString(32) // 用于加密会话Cookie的密钥
	flagSessionMaxDuration       = 90            // 记住的会话刷新和有效的最大天数
	flagWgConfTemplate           string          // 自定义wg.conf模板的路径
	flagBasePath                 string          // URL的基础路径
	flagSubnetRanges             string          // 为客户端分配IP时选择的IP范围
)

const (
	defaultEmailSubject = "您的WireGuard配置" // 默认邮件主题
	defaultEmailContent = `您好,</br>
<p>在这封邮件中，您可以找到我们WireGuard服务器的个人配置。</p>

<p>祝好</p>
` // 默认邮件内容
)

// 嵌入"templates"目录
//
//go:embed templates/*
var embeddedTemplates embed.FS

// 嵌入"assets"目录
//
//go:embed assets/*
var embeddedAssets embed.FS

func init() {
	// 命令行标志和环境变量
	flag.BoolVar(&flagDisableLogin, "disable-login", util.LookupEnvOrBool("DISABLE_LOGIN", flagDisableLogin), "禁用应用程序的身份验证，这可能存在安全风险。")
	flag.StringVar(&flagBindAddress, "bind-address", util.LookupEnvOrString("BIND_ADDRESS", flagBindAddress), "应用程序将绑定的地址:端口。")
	flag.StringVar(&flagSmtpHostname, "smtp-hostname", util.LookupEnvOrString("SMTP_HOSTNAME", flagSmtpHostname), "SMTP主机名")
	flag.IntVar(&flagSmtpPort, "smtp-port", util.LookupEnvOrInt("SMTP_PORT", flagSmtpPort), "SMTP端口")
	flag.StringVar(&flagSmtpHelo, "smtp-helo", util.LookupEnvOrString("SMTP_HELO", flagSmtpHelo), "SMTP HELO主机名")
	flag.StringVar(&flagSmtpUsername, "smtp-username", util.LookupEnvOrString("SMTP_USERNAME", flagSmtpUsername), "SMTP用户名")
	flag.BoolVar(&flagSmtpNoTLSCheck, "smtp-no-tls-check", util.LookupEnvOrBool("SMTP_NO_TLS_CHECK", flagSmtpNoTLSCheck), "禁用SMTP的TLS验证，这可能存在安全风险。")
	flag.StringVar(&flagSmtpEncryption, "smtp-encryption", util.LookupEnvOrString("SMTP_ENCRYPTION", flagSmtpEncryption), "SMTP加密方式：NONE、SSL、SSLTLS、TLS或STARTTLS（默认）")
	flag.StringVar(&flagSmtpAuthType, "smtp-auth-type", util.LookupEnvOrString("SMTP_AUTH_TYPE", flagSmtpAuthType), "SMTP认证类型：PLAIN、LOGIN或NONE。")
	flag.StringVar(&flagEmailFrom, "email-from", util.LookupEnvOrString("EMAIL_FROM_ADDRESS", flagEmailFrom), "发件人邮箱地址。")
	flag.StringVar(&flagEmailFromName, "email-from-name", util.LookupEnvOrString("EMAIL_FROM_NAME", flagEmailFromName), "发件人邮箱名称。")
	flag.StringVar(&flagTelegramToken, "telegram-token", util.LookupEnvOrString("TELEGRAM_TOKEN", flagTelegramToken), "用于向客户端分发配置的Telegram机器人令牌。")
	flag.BoolVar(&flagTelegramAllowConfRequest, "telegram-allow-conf-request", util.LookupEnvOrBool("TELEGRAM_ALLOW_CONF_REQUEST", flagTelegramAllowConfRequest), "允许用户通过发送消息从机器人获取配置。")
	flag.IntVar(&flagTelegramFloodWait, "telegram-flood-wait", util.LookupEnvOrInt("TELEGRAM_FLOOD_WAIT", flagTelegramFloodWait), "处理下一个配置请求前的等待时间（分钟）。")
	flag.StringVar(&flagWgConfTemplate, "wg-conf-template", util.LookupEnvOrString("WG_CONF_TEMPLATE", flagWgConfTemplate), "自定义wg.conf模板的路径。")
	flag.StringVar(&flagBasePath, "base-path", util.LookupEnvOrString("BASE_PATH", flagBasePath), "URL的基础路径")
	flag.StringVar(&flagSubnetRanges, "subnet-ranges", util.LookupEnvOrString("SUBNET_RANGES", flagSubnetRanges), "为客户端分配IP时选择的IP范围。")
	flag.IntVar(&flagSessionMaxDuration, "session-max-duration", util.LookupEnvOrInt("SESSION_MAX_DURATION", flagSessionMaxDuration), "记住的会话刷新和有效的最大天数。")

	var (
		smtpPasswordLookup   = util.LookupEnvOrString("SMTP_PASSWORD", flagSmtpPassword)
		sendgridApiKeyLookup = util.LookupEnvOrString("SENDGRID_API_KEY", flagSendgridApiKey)
		sessionSecretLookup  = util.LookupEnvOrString("SESSION_SECRET", flagSessionSecret)
	)

	// 检查空的smtpPassword环境变量
	if smtpPasswordLookup != "" {
		flag.StringVar(&flagSmtpPassword, "smtp-password", smtpPasswordLookup, "SMTP密码")
	} else {
		flag.StringVar(&flagSmtpPassword, "smtp-password", util.LookupEnvOrFile("SMTP_PASSWORD_FILE", flagSmtpPassword), "SMTP密码文件")
	}

	// 检查空的sendgridApiKey环境变量
	if sendgridApiKeyLookup != "" {
		flag.StringVar(&flagSendgridApiKey, "sendgrid-api-key", sendgridApiKeyLookup, "您的SendGrid API密钥。")
	} else {
		flag.StringVar(&flagSendgridApiKey, "sendgrid-api-key", util.LookupEnvOrFile("SENDGRID_API_KEY_FILE", flagSendgridApiKey), "包含您的SendGrid API密钥的文件。")
	}

	// 检查空的sessionSecret环境变量
	if sessionSecretLookup != "" {
		flag.StringVar(&flagSessionSecret, "session-secret", sessionSecretLookup, "用于加密会话Cookie的密钥。")
	} else {
		flag.StringVar(&flagSessionSecret, "session-secret", util.LookupEnvOrFile("SESSION_SECRET_FILE", flagSessionSecret), "包含用于加密会话Cookie的密钥的文件。")
	}

	flag.Parse()

	// 更新运行时配置
	util.DisableLogin = flagDisableLogin
	util.BindAddress = flagBindAddress
	util.SmtpHostname = flagSmtpHostname
	util.SmtpPort = flagSmtpPort
	util.SmtpHelo = flagSmtpHelo
	util.SmtpUsername = flagSmtpUsername
	util.SmtpPassword = flagSmtpPassword
	util.SmtpAuthType = flagSmtpAuthType
	util.SmtpNoTLSCheck = flagSmtpNoTLSCheck
	util.SmtpEncryption = flagSmtpEncryption
	util.SendgridApiKey = flagSendgridApiKey
	util.EmailFrom = flagEmailFrom
	util.EmailFromName = flagEmailFromName
	util.SessionSecret = sha512.Sum512([]byte(flagSessionSecret))
	util.SessionMaxDuration = int64(flagSessionMaxDuration) * 86_400 // 以秒为单位存储
	util.WgConfTemplate = flagWgConfTemplate
	util.BasePath = util.ParseBasePath(flagBasePath)
	util.SubnetRanges = util.ParseSubnetRanges(flagSubnetRanges)

	lvl, _ := util.ParseLogLevel(util.LookupEnvOrString(util.LogLevel, "INFO"))

	telegram.Token = flagTelegramToken
	telegram.AllowConfRequest = flagTelegramAllowConfRequest
	telegram.FloodWait = flagTelegramFloodWait
	telegram.LogLevel = lvl

	// 仅在日志级别为INFO或更低时打印
	if lvl <= log.INFO {
		// 打印应用程序信息
		fmt.Println("Wireguard UI")
		fmt.Println("应用程序版本\t:", appVersion)
		fmt.Println("Git提交\t:", gitCommit)
		fmt.Println("Git引用\t\t:", gitRef)
		fmt.Println("构建时间\t:", buildTime)
		fmt.Println("Git仓库\t:", "https://github.com/ngoduykhanh/wireguard-ui")
		fmt.Println("身份验证\t:", !util.DisableLogin)
		fmt.Println("绑定地址\t:", util.BindAddress)
		//fmt.Println("Sendgrid密钥\t:", util.SendgridApiKey)
		fmt.Println("发件人邮箱\t:", util.EmailFrom)
		fmt.Println("发件人邮箱名称\t:", util.EmailFromName)
		//fmt.Println("会话密钥\t:", util.SessionSecret)
		fmt.Println("自定义wg.conf\t:", util.WgConfTemplate)
		fmt.Println("基础路径\t:", util.BasePath+"/")
		fmt.Println("子网范围\t:", util.GetSubnetRangesString())
	}
}

func main() {
	db, err := jsondb.New("./db")
	if err != nil {
		panic(err)
	}
	if err := db.Init(); err != nil {
		panic(err)
	}
	// 设置应用程序额外数据
	extraData := make(map[string]interface{})
	extraData["appVersion"] = appVersion
	extraData["gitCommit"] = gitCommit
	extraData["basePath"] = util.BasePath
	extraData["loginDisabled"] = flagDisableLogin

	// 从嵌入式目录中去除"templates/"前缀，以便可以通过直接名称读取文件（例如"base.html"而不是"templates/base.html"）
	tmplDir, _ := fs.Sub(fs.FS(embeddedTemplates), "templates")

	// 启动时创建WireGuard配置（如果不存在）
	initServerConfig(db, tmplDir)

	// 检查子网范围是否对服务器配置有效
	// 移除任何无效的CIDR
	if err := util.ValidateAndFixSubnetRanges(db); err != nil {
		panic(err)
	}

	// 打印有效范围
	if lvl, _ := util.ParseLogLevel(util.LookupEnvOrString(util.LogLevel, "INFO")); lvl <= log.INFO {
		fmt.Println("有效子网范围:", util.GetSubnetRangesString())
	}

	// 注册路由
	app := router.New(tmplDir, extraData, util.SessionSecret)

	app.GET(util.BasePath, handler.WireGuardClients(db), handler.ValidSession, handler.RefreshSession)

	// 重要：确保所有非GET路由都使用handler.ContentTypeJson检查请求内容类型，以减轻CSRF攻击。这是有效的，因为浏览器不允许在跨域请求上设置Content-Type标头。

	if !util.DisableLogin {
		app.GET(util.BasePath+"/login", handler.LoginPage())
		app.POST(util.BasePath+"/login", handler.Login(db), handler.ContentTypeJson)
		app.GET(util.BasePath+"/logout", handler.Logout(), handler.ValidSession)
		app.GET(util.BasePath+"/profile", handler.LoadProfile(), handler.ValidSession, handler.RefreshSession)
		app.GET(util.BasePath+"/users-settings", handler.UsersSettings(), handler.ValidSession, handler.RefreshSession, handler.NeedsAdmin)
		app.POST(util.BasePath+"/update-user", handler.UpdateUser(db), handler.ValidSession, handler.ContentTypeJson)
		app.POST(util.BasePath+"/create-user", handler.CreateUser(db), handler.ValidSession, handler.ContentTypeJson, handler.NeedsAdmin)
		app.POST(util.BasePath+"/remove-user", handler.RemoveUser(db), handler.ValidSession, handler.ContentTypeJson, handler.NeedsAdmin)
		app.GET(util.BasePath+"/get-users", handler.GetUsers(db), handler.ValidSession, handler.NeedsAdmin)
		app.GET(util.BasePath+"/api/user/:username", handler.GetUser(db), handler.ValidSession)
	}

	var sendmail emailer.Emailer
	if util.SendgridApiKey != "" {
		sendmail = emailer.NewSendgridApiMail(util.SendgridApiKey, util.EmailFromName, util.EmailFrom)
	} else {
		sendmail = emailer.NewSmtpMail(util.SmtpHostname, util.SmtpPort, util.SmtpUsername, util.SmtpPassword, util.SmtpHelo, util.SmtpNoTLSCheck, util.SmtpAuthType, util.EmailFromName, util.EmailFrom, util.SmtpEncryption)
	}

	app.GET(util.BasePath+"/test-hash", handler.GetHashesChanges(db), handler.ValidSession)
	app.GET(util.BasePath+"/about", handler.AboutPage())
	app.GET(util.BasePath+"/_health", handler.Health())
	app.GET(util.BasePath+"/favicon", handler.Favicon())
	app.POST(util.BasePath+"/new-client", handler.NewClient(db), handler.ValidSession, handler.ContentTypeJson)
	app.POST(util.BasePath+"/update-client", handler.UpdateClient(db), handler.ValidSession, handler.ContentTypeJson)
	app.POST(util.BasePath+"/email-client", handler.EmailClient(db, sendmail, defaultEmailSubject, defaultEmailContent), handler.ValidSession, handler.ContentTypeJson)
	app.POST(util.BasePath+"/send-telegram-client", handler.SendTelegramClient(db), handler.ValidSession, handler.ContentTypeJson)
	app.POST(util.BasePath+"/client/set-status", handler.SetClientStatus(db), handler.ValidSession, handler.ContentTypeJson)
	app.POST(util.BasePath+"/remove-client", handler.RemoveClient(db), handler.ValidSession, handler.ContentTypeJson)
	app.GET(util.BasePath+"/download", handler.DownloadClient(db), handler.ValidSession)
	app.GET(util.BasePath+"/wg-server", handler.WireGuardServer(db), handler.ValidSession, handler.RefreshSession, handler.NeedsAdmin)
	app.POST(util.BasePath+"/wg-server/interfaces", handler.WireGuardServerInterfaces(db), handler.ValidSession, handler.ContentTypeJson, handler.NeedsAdmin)
	app.POST(util.BasePath+"/wg-server/keypair", handler.WireGuardServerKeyPair(db), handler.ValidSession, handler.ContentTypeJson, handler.NeedsAdmin)
	app.GET(util.BasePath+"/global-settings", handler.GlobalSettings(db), handler.ValidSession, handler.RefreshSession, handler.NeedsAdmin)
	app.POST(util.BasePath+"/global-settings", handler.GlobalSettingSubmit(db), handler.ValidSession, handler.ContentTypeJson, handler.NeedsAdmin)
	app.GET(util.BasePath+"/status", handler.Status(db), handler.ValidSession, handler.RefreshSession)
	app.GET(util.BasePath+"/api/clients", handler.GetClients(db), handler.ValidSession)
	app.GET(util.BasePath+"/api/client/:id", handler.GetClient(db), handler.ValidSession)
	app.GET(util.BasePath+"/api/machine-ips", handler.MachineIPAddresses(), handler.ValidSession)
	app.GET(util.BasePath+"/api/subnet-ranges", handler.GetOrderedSubnetRanges(), handler.ValidSession)
	app.GET(util.BasePath+"/api/suggest-client-ips", handler.SuggestIPAllocation(db), handler.ValidSession)
	app.POST(util.BasePath+"/api/apply-wg-config", handler.ApplyServerConfig(db, tmplDir), handler.ValidSession, handler.ContentTypeJson)
	app.GET(util.BasePath+"/wake_on_lan_hosts", handler.GetWakeOnLanHosts(db), handler.ValidSession, handler.RefreshSession)
	app.POST(util.BasePath+"/wake_on_lan_host", handler.SaveWakeOnLanHost(db), handler.ValidSession, handler.ContentTypeJson)
	app.DELETE(util.BasePath+"/wake_on_lan_host/:mac_address", handler.DeleteWakeOnHost(db), handler.ValidSession, handler.ContentTypeJson)
	app.PUT(util.BasePath+"/wake_on_lan_host/:mac_address", handler.WakeOnHost(db), handler.ValidSession, handler.ContentTypeJson)

	// 从嵌入式目录中去除"assets/"前缀，以便可以直接调用文件而无需"assets/"前缀
	assetsDir, _ := fs.Sub(fs.FS(embeddedAssets), "assets")
	assetHandler := http.FileServer(http.FS(assetsDir))
	// 提供其他静态文件
	app.GET(util.BasePath+"/static/*", echo.WrapHandler(http.StripPrefix(util.BasePath+"/static/", assetHandler)))

	initDeps := telegram.TgBotInitDependencies{
		DB:                             db,
		SendRequestedConfigsToTelegram: util.SendRequestedConfigsToTelegram,
	}

	initTelegram(initDeps)

	if strings.HasPrefix(util.BindAddress, "unix://") {
		// 监听UNIX域套接字
		// https://github.com/labstack/echo/issues/830
		err := syscall.Unlink(util.BindAddress[6:])
		if err != nil {
			app.Logger.Fatalf("无法取消链接UNIX套接字：错误：%v", err)
		}
		l, err := net.Listen("unix", util.BindAddress[6:])
		if err != nil {
			app.Logger.Fatalf("无法创建UNIX套接字。错误：%v", err)
		}
		app.Listener = l
		app.Logger.Fatal(app.Start(""))
	} else {
		// 监听TCP套接字
		app.Logger.Fatal(app.Start(util.BindAddress))
	}
}

// 初始化服务器配置
func initServerConfig(db store.IStore, tmplDir fs.FS) {
	settings, err := db.GetGlobalSettings()
	if err != nil {
		log.Fatalf("无法获取全局设置：%v", err)
	}

	if _, err := os.Stat(settings.ConfigFilePath); err == nil {
		// 文件存在，不隐式覆盖
		return
	}

	server, err := db.GetServer()
	if err != nil {
		log.Fatalf("无法获取服务器配置：%v", err)
	}

	clients, err := db.GetClients(false)
	if err != nil {
		log.Fatalf("无法获取客户端配置：%v", err)
	}

	users, err := db.GetUsers()
	if err != nil {
		log.Fatalf("无法获取用户配置：%v", err)
	}

	// 写入配置文件
	err = util.WriteWireGuardServerConfig(tmplDir, server, clients, users, settings)
	if err != nil {
		log.Fatalf("无法创建服务器配置：%v", err)
	}
}

// 初始化Telegram机器人
func initTelegram(initDeps telegram.TgBotInitDependencies) {
	go func() {
		for {
			err := telegram.Start(initDeps)
			if err == nil {
				break
			}
		}
	}()
}

//获取CST时间
func getCSTBuildTime() string {
	// 加载中国标准时区（UTC+8）
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 时区加载失败时返回UTC时间
		return time.Now().Format("01-02-2006 15:04:05")
	}
	// 转换为CST时间并格式化
	return time.Now().In(loc).Format("01-02-2006 15:04:05")
}
