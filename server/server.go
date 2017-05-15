package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
)

var (
	host = flag.String("host", ":443", "HTTPS 服务器 <ip>:<port>")
	//domainCertificateFile = flag.String("domain_certificate_file", "", "域名 HTTPS 证书文件")
	//domainKeyFile         = flag.String("domain_key_file", "", "域名 HTTPS key 文件")
	staticFile = flag.String("static_file", "", "")
	appID      = flag.String("appid", "", "你的 app 的 application-identifier，通常是 <team id>.<bundle id>")
)

func main() {
	flag.Parse()

	log.Print("启动 HTTPS 服务器")
	http.HandleFunc("/.well-known/apple-app-site-association", ulinkService)
	http.HandleFunc("/apple-app-site-association", ulinkService)
	http.HandleFunc("/", staticServer)
	//log.Fatal(http.ListenAndServeTLS(*host, *domainCertificateFile, *domainKeyFile, nil))
	log.Fatal(http.ListenAndServe(*host, nil))
}

func staticServer(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, *staticFile)
	return
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

func ulinkService(w http.ResponseWriter, req *http.Request) {
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
