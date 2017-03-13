package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	host                  = flag.String("host", ":443", "HTTPS 服务器 <ip>:<port>")
	domainCertificateFile = flag.String("domain_certificate_file", "", "域名 HTTPS 证书文件")
	domainKeyFile         = flag.String("domain_key_file", "", "域名 HTTPS key 文件")
	staticFile            = flag.String("static_file", "", "")
)

func main() {
	flag.Parse()

	log.Print("启动 HTTPS 服务器")
	http.HandleFunc("/", staticServer)
	log.Fatal(http.ListenAndServeTLS(*host, *domainCertificateFile, *domainKeyFile, nil))
}

func staticServer(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, *staticFile)
	return
}
