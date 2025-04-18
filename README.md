# 简易的 DNS 管理

> 此项目仅用于小规模部署和非生产部署, 请**不要**在 HTTP 下使用此应用

## 功能概览

- DNS 管理
- TOKEN 管理
- 操作审计

## 使用说明

安装 golang 1.24 后，执行如下命令:

```bash
git clone https://github.com/console-dns/server
cd server
PROD=true make console-dns conf/server.yaml
sed -i -e 's|conf/|/var/lib/console-dns/|g' conf/server.yaml
mkdir -p /usr/lib/console-dns /usr/log/console-dns /etc/console-dns
install -Dm0644 conf/server.yaml /etc/console-dns/server.yaml
install -Dm0644 console-dns.service /etc/systemd/system/console-dns.service
install -Dm0755 console-dns /usr/local/bin/console-dns
systemctl daemon-reload
systemctl start console-dns.service
```

## TODO

- ~~支持 ldap 和 oidc 登录~~ (关键组件无需统一管理)
- 重构 web 页面的 DNS 修改/删除，改为按照内容执行操作
- 覆盖测试用例

## 已知问题

- 可能存在 xss 注入问题，仅在登录后生效
- 由于使用 `rw-lock` , 可能存在性能低下的问题

## License 

项目使用 MIT License