
1. 生成 key.pem 和 cert.pem (例如 openssl )
```shell
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 10000 -nodes
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:Beijing
Locality Name (eg, city) []:Beijing
Organization Name (eg, company) [Internet Widgits Pty Ltd]:AnthonyDong
Organizational Unit Name (eg, section) []:Development
Common Name (e.g. server FQDN or YOUR name) []:DevTool
Email Address []:fanhaodong516@gmail.com
```



2. 查看 cert.pem

```shell
openssl x509 -in cert.pem -text -noout
```


3. 配置

```shell
sudo cp cert.pem /usr/local/share/ca-certificates/devtool.crt

sudo update-ca-certificates
```


