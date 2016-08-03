ilink
==

Go 写的 iOS universal links 转发服务

Universal Links 是 iOS 9 提供的应用间跳转服务，允许你从一个 app 跳转到另一个 app 的某个页面。关于 Universal Links 的实现细节，见 Apple [官方文档](https://developer.apple.com/library/ios/documentation/General/Conceptual/AppSearch/UniversalLinks.html)。

ilink 实现了 Universal Links 的服务端：

1、当用户安装了该 app 时，跳转到 app 中详情页

2、当用户没有安装 app 时，跳转到 app store 的安装页

## 编译

```
go get github.com/yingtu/ilink
```

进入 ilink/server 目录执行

```
go build
```

## 运行

第一次运行需要启动 --renew_certificate 参数，以便从 Let's Encrypt 得到 HTTPS 证书

```
./server --domain <你的域名，不包含 https://> --email <你的邮箱> \
  --domain_certificate_file <你的域名的 crt 文件> --domain_key_file <你的域名的 key 文件> \
  --appid <你的域名的 appid> --redirect_url <你的 app 的 app store 网址> \
  --renew_certificate
```

证书会写在 --domain_certificate_file 和 --domain_key_file 里面。第二次运行请省略 --renew_certificate，其他参数保持不变。

参数含义：

> --domain：app 跳转的中转域名，ilink 会为这个域名从 Let's Encrypt 申请 HTTPS 证书，你需要保证 server 有绑定 443 和 80 端口的权限，并且可以被外网访问到。当然，DNS 也需要将这个域名指定到这台服务器。

> --email：Let's Encrypt 注册邮箱，可以随便填写

> --domain_certificate_file 和 --domain_key_file：从 Let's Encrypt 得到的证书会保存在这个文件，HTTPS 服务也会使用这个文件启动 TLS。

> --appid：你的 app ID，通常为 team_id.bundle_id 的格式

> --redirect_url：当用户没有安装 app 时，会访问你的域名，然后被重定向到这个域名，请设置为你的 app 的 app store 页面地址

> --renew_certificate：当带这个参数时，会从 Let's Encrypt 获得证书和钥匙，并写入 --domain_certificate_file 和 --domain_key_file 两个文件。当不带这个参数时，只从这两个文件读入配置。注意 Let's Encrypt 的证书每 90 天需要更新一次，但是，请不要频繁 renew certificate，每个星期最多 renew 5 次（Let's Encrypt 的限制）。
