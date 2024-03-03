package proxy

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"math/big"
	math_rand "math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

type certificateCacheHelper struct {
	lock  sync.RWMutex
	cache map[string]*tls.Certificate
	max   int
}

func (c *certificateCacheHelper) Get(hosts []string) *tls.Certificate {
	key := strings.Join(hosts, "|")
	c.lock.RLock()
	defer c.lock.RUnlock()

	result := c.cache[key]
	return result
}

func (c *certificateCacheHelper) Set(hosts []string, value *tls.Certificate) {
	key := strings.Join(hosts, "|")
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]*tls.Certificate, c.max*2)
	}

	if len(c.cache) >= c.max && c.max > 1 {
		for deleteKey := range c.cache {
			delete(c.cache, deleteKey)
			break
		}
	}

	c.cache[key] = value
}

func stripPort(host string) string {
	if index := strings.LastIndex(host, ":"); index != -1 {
		host = host[:index]
	}
	if host != "" && host[0] == '[' && host[len(host)-1] == ']' {
		return host[1 : len(host)-1]
	}
	return host
}

func NewCertificate(root *tls.Certificate, hosts []string, cache *certificateCacheHelper) (*tls.Certificate, error) {
	if cache != nil {
		if ret := cache.Get(hosts); ret != nil {
			return ret, nil
		}
	}
	rootCert, err := x509.ParseCertificate(root.Certificate[0])
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(math_rand.Int63()),
		Issuer:       rootCert.Subject,
		Subject: pkix.Name{
			Organization: rootCert.Subject.Organization,
		},
		NotBefore:             time.Now().AddDate(0, 0, -1),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
			template.Subject.CommonName = h
		}
	}

	var randReader io.Reader = rand.Reader
	// 随机生成一个新的私钥
	var newKey crypto.Signer
	switch root.PrivateKey.(type) {
	case *rsa.PrivateKey:
		if newKey, err = rsa.GenerateKey(randReader, 2048); err != nil {
			return nil, err
		}
	case *ecdsa.PrivateKey:
		if newKey, err = ecdsa.GenerateKey(elliptic.P256(), randReader); err != nil {
			return nil, err
		}
	//todo ed25519.PrivateKey:
	default:
		return nil, fmt.Errorf("unsupported key type %T", root.PrivateKey)
	}
	// 使用根证书和根证书的私钥对新证书进行签名
	newCertDER, err := x509.CreateCertificate(randReader, &template, rootCert, newKey.Public(), root.PrivateKey)
	if err != nil {
		return nil, err
	}
	ret := &tls.Certificate{
		Certificate: [][]byte{newCertDER, root.Certificate[0]},
		PrivateKey:  newKey,
	}
	if cache != nil {
		cache.Set(hosts, ret)
	}
	return ret, nil
}

// issuer: C=CN; ST=Beijing; L=Beijing; O=AnthonyDong; OU=Development; CN=DevTool; emailAddress=fanhaodong516@gmail.com

// 根证书
var CA_CERT = []byte(`-----BEGIN CERTIFICATE-----
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
xwnqqjxynPfcd8l5qJFLVv4WcW6HbZs=
-----END CERTIFICATE-----`)

// 根证书私钥
var CA_KEY = []byte(`-----BEGIN PRIVATE KEY-----
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
Uplstu6NDzGABLX/uqKVELsZtShq5g==
-----END PRIVATE KEY-----`)
