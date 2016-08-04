package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

var (
	host                  = flag.String("host", ":443", "HTTPS 服务器 <ip>:<port>")
	domain                = flag.String("domain", "", "HTTPS 服务域名")
	email                 = flag.String("email", "", "Let's encrypt 注册邮箱")
	domainCertificateFile = flag.String("domain_certificate_file", "", "域名 HTTPS 证书文件")
	domainKeyFile         = flag.String("domain_key_file", "", "域名 HTTPS key 文件")
	renewCertificate      = flag.Bool("renew_certificate", false, "是否更新证书")
	appID                 = flag.String("appid", "", "你的 app 的 application-identifier，通常是 <team id>.<bundle id>")
	redirectURL           = flag.String("redirect_url", "", "当用户没有安装 app，跳转到这个网址")
)

func main() {
	flag.Parse()

	// 检查 flag 正确性
	if *appID == "" {
		log.Fatal("--appID 参数不能为空")
	}
	if *redirectURL == "" {
		log.Fatal("--redirect_url 参数不能为空")
	}

	// 更新证书
	// Let's Encrypt 的证书每 90 天失效一次，只需要在失效前更新即可
	// 更新的频率每个星期不超过 5 次（Let's Encrypt 的限制）
	if *renewCertificate {
		if *domain == "" {
			log.Fatal("--domain 参数不能为空")
		}
		if *email == "" {
			log.Fatal("--email 参数不能为空")
		}
		if err := getCertificate(*domain, *email, *domainCertificateFile, *domainKeyFile); err != nil {
			log.Fatal(err)
		}
	}

	// 启动服务
	r := http.NewServeMux()
	// 用户请求 app-site 关联文件时，返回 json 配置
	r.Handle("/.well-known/apple-app-site-association", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(ULinkService)))
	r.Handle("/apple-app-site-association", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(ULinkService)))
	// 否则（当用户手机上没有安装该 app），跳转到指定页面（通常是 app store 页面）
	r.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(RedirectService)))
	log.Print("启动 HTTPS 服务器")
	log.Fatal(http.ListenAndServeTLS(*host, *domainCertificateFile, *domainKeyFile, handlers.CompressHandler(r)))
}

func ULinkService(w http.ResponseWriter, req *http.Request) {
	saf := SiteAssociationFile{
		Applinks: Applinks{
			Apps: []string{},
			Details: []Detail{
				Detail{
					AppID: *appID,
					Paths: []string{"*"},
				},
			},
		},
	}
	fileStr, err := json.Marshal(saf)
	if err != nil {
		return
	}

	io.WriteString(w, string(fileStr))
}

type SiteAssociationFile struct {
	Applinks Applinks `json:"applinks"`
}

type Applinks struct {
	Apps    []string `json:"apps"`
	Details []Detail `json:"details"`
}

type Detail struct {
	AppID string   `json:"appID"`
	Paths []string `json:"paths"`
}

func RedirectService(w http.ResponseWriter, req *http.Request) {
	// 检测是否是从其他页面跳转而来（非直接打开）
	hasReferer := (req.Referer() != "")

	// 提取 redirection 参数
	query := req.URL.Query()
	redirection := query.Get("redirection")
	var err error
	redirection, err = url.QueryUnescape(redirection)
	if err != nil {
		log.Printf("redirection 解码错误 %s", err)
		return
	}

	// 判断是否为微信打开
	inWechat := strings.Contains(strings.ToLower(req.UserAgent()), "micromessenger")

	// 当不在微信里打开，跳转到 --redirect_url
	if !inWechat {
		log.Printf("跳转到 %s", *redirectURL)
		http.Redirect(w, req, *redirectURL, 301)
		return
	}

	// 当在微信里打开，且 refer 不为空时，不做跳转
	if inWechat && hasReferer {
		log.Printf("不跳转")
		return
	}

	// 当在微信里打开，且 refer 为空，且 redirection 不为空时，跳到 redirection 地址
	if inWechat && !hasReferer && redirection != "" {
		log.Printf("跳转到 %s", redirection)
		http.Redirect(w, req, redirection, 301)
		return
	}
}
