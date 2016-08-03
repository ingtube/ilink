package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"

	"github.com/xenolf/lego/acme"
)

const (
	letsencryptEndpoint = "https://acme-v01.api.letsencrypt.org/directory"
)

// 实现 acme.User
type MyUser struct {
	Email        string
	Registration *acme.RegistrationResource
	key          crypto.PrivateKey
}

func (u MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}
func (u MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// 从 Let's Encrypt 得到证书
func getCertificate(domain string, email string, certificateFile string, keyFile string) error {
	// 生成认证用户
	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		log.Fatal(err)
	}
	myUser := MyUser{
		Email: email,
		key:   privateKey,
	}

	// 建立客户端
	client, err := acme.NewClient(letsencryptEndpoint, &myUser, acme.RSA2048)
	if err != nil {
		return err
	}

	// 注册用户
	reg, err := client.Register()
	if err != nil {
		return err
	}
	myUser.Registration = reg

	// 同意 Let's Encrypt 的服务协议
	err = client.AgreeToTOS()
	if err != nil {
		return err
	}

	// 获取认证，本机必须能够访问认证域名
	bundle := false
	certificates, failures := client.ObtainCertificate([]string{domain}, bundle, nil)
	if len(failures) > 0 {
		return errors.New("无法得到证书")
	}

	// 读入 Let's Encrypt 的证书
	intermediate, err := ioutil.ReadFile("intermediate.pem")
	if err != nil {
		return errors.New("无法读入 intermediate pem")
	}

	// 生成域名认证文件
	if err := ioutil.WriteFile(
		certificateFile,
		bytes.Join([][]byte{certificates.Certificate, intermediate}, nil),
		0600); err != nil {
		return err
	}

	// 生成域名认证秘钥
	if certificates.PrivateKey != nil {
		err = ioutil.WriteFile(keyFile, certificates.PrivateKey, 0600)
		if err != nil {
			return err
		}
	}

	return nil
}
