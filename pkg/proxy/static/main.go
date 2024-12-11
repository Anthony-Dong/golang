package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	// 生成密钥和证书
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	rootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2048),
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{"DevTool"},
			Locality:     []string{"Beijing"},
		},
		NotBefore:             time.Now().AddDate(-1, 0, 0),
		NotAfter:              time.Now().AddDate(9, 0, 0), // 有效期设置为10年。
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, key.Public(), key)
	if err != nil {
		panic(err)
	}
	if _, err := x509.ParseCertificate(certBytes); err != nil {
		panic(err)
	}
	// 保存为PEM格式
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	// -keyout key.pem -out cert.pem
	if err := os.WriteFile("key.pem", pemKey, 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("cert.pem", pemCert, 0644); err != nil {
		panic(err)
	}
	fmt.Println("write \"key.pem\" \"cert.pem\" success. \n" +
		"openssl x509 -in cert.pem -text -noout")
}
