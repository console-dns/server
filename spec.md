# Console DNS RESTful API 规范 (v1)

本文档定义了 Console DNS 系统的高级 RESTful API 规范。本设计优先考虑并发安全性、操作原子性以及大规模数据处理能力。

---

## 1. 通用约定

### 1.1 基础信息
- **Base URL**: `/api/v1`
- **内容协商**: `Content-Type: application/json`
- **认证**: `Authorization: Bearer <token>`

### 1.2 标准响应格式
- **成功**: 直接返回对象或数组。
- **错误**: 
  ```json
  {
    "error": {
      "code": "INVALID_PARAMETER",
      "message": "字段校验失败",
      "details": { "ttl": "必须在 1-86400 之间" }
    }
  }
  ```

### 1.3 分页与过滤 (通用参数)
- `page`: 分页页码 (从 1 开始)
- `per_page`: 每页数量 (默认 20，最大 100)
- `q`: 关键词模糊搜索
- `sort`: 排序字段 (如 `created_at:desc`)

---

## 2. 区域管理 (Zones)

### 2.1 获取区域列表
- **GET** `/zones`
- **支持过滤**: `?q=example.com`

### 2.2 创建区域
- **POST** `/zones`
- **请求体**: `{ "name": "example.com" }`

### 2.3 获取区域详情
- **GET** `/zones/{zone}`
- **说明**: 返回区域元数据及解析统计概要。

### 2.4 删除区域
- **DELETE** `/zones/{zone}`

---

## 3. 记录集管理 (Record Sets - RRSet)

**设计理念**: 放弃不稳定的数组索引，采用“记录集”模式。一个记录集由 `域名 + 类型` 唯一确定。

### 3.1 获取区域内所有记录
- **GET** `/zones/{zone}/records`
- **支持过滤**: `?name=www&type=A`
- **分页**: 支持标准分页参数。

### 3.2 覆盖更新记录集 (原子操作)
- **PUT** `/zones/{zone}/records/{name}/{type}`
- **说明**: 此操作会**完全替换**该名称和类型下的所有现有值。这是保证 DNS 状态一致性的推荐做法。
- **请求体**:
  ```json
  {
    "ttl": 3600,
    "records": [
      { "content": "1.1.1.1" },
      { "content": "1.1.1.2" }
    ]
  }
  ```

### 3.3 局部修改记录 (PATCH)
- **PATCH** `/zones/{zone}/records/{name}/{type}`
- **请求体**:
  ```json
  {
    "action": "add", // 或 "remove"
    "record": { "content": "1.1.1.3" }
  }
  ```

---

## 4. 客户端与安全 (Clients)

### 4.1 获取客户端状态
- **GET** `/clients`
- **返回**: 包含客户端配置、最近活动 IP 及 Token 状态。

### 4.2 创建/重置客户端
- **POST** `/clients`
- **POST** `/clients/{client_id}/token/reset` (重置并获取新 Secret)

### 4.3 访问策略 (Policies)
- **PUT** `/clients/{client_id}/policies/{policy_id}`
- **说明**: 定义该客户端可操作的域名匹配模式和权限等级。
- **参数**: 
  - `pattern`: 正则表达式（如 `.*\.example\.com`）
  - `access`: `ro` (只读), `ru` (读/更新), `rw` (完全控制)

---

## 5. 审计日志 (Audit Logs)

### 5.1 查询审计记录
- **GET** `/logs`
- **核心过滤器**:
  - `author`: 执行者 UID
  - `zone`: 关联的区域
  - `action`: 操作类型 (create, update, delete)
- **响应**: 包含标准分页包装。

---

## 附录：DNS 数据模型与验证规范

所有记录项（Record Item）必须包含基础字段 `ttl`，具体的业务数据根据 `type` 的不同定义在 `content` 中（在 PUT 覆盖操作中，`content` 为数组项的内容）。

### 通用约束
- **TTL**: 32 位无符号整数，建议范围 `60` 至 `86400`。
- **主机名/域名**: 必须符合 RFC 1035 规范，仅限 `a-z`, `0-9`, `-`, `.`，且不能以连字符开头或结尾。末尾的 `.` 表示绝对域名（FQDN）。

---

### 1. A 记录 (IPv4 地址)
- **说明**: 将域名指向一个 IPv4 地址。
- **字段**:
  - `ip` (String): 合法的 IPv4 地址（如 `192.168.1.1`）。
- **示例**:
  ```json
  { "ip": "1.2.3.4", "ttl": 600 }
  ```

### 2. AAAA 记录 (IPv6 地址)
- **说明**: 将域名指向一个 IPv6 地址。
- **字段**:
  - `ip` (String): 合法的 IPv6 地址（如 `2001:db8::1`）。
- **示例**:
  ```json
  { "ip": "2400:3200::1", "ttl": 600 }
  ```

### 3. CNAME / NS 记录 (别名与名称服务器)
- **说明**: CNAME 指向另一个域名；NS 指定该域名的授权名称服务器。
- **字段**:
  - `host` (String): 目标主机名或域名。
- **示例**:
  ```json
  { "host": "cname.dnsv1.com.", "ttl": 3600 }
  ```

### 4. MX 记录 (邮件交换)
- **说明**: 指定接收该域名电子邮件的服务器。
- **字段**:
  - `host` (String): 邮件服务器地址。
  - `preference` (Uint16): 优先级 (0-65535)，数值越小优先级越高。
- **示例**:
  ```json
  { "host": "mx1.example.com.", "preference": 10, "ttl": 600 }
  ```

### 5. TXT 记录 (文本)
- **说明**: 通常用于 SPF 记录、DKIM 签名或域名所有权验证。
- **字段**:
  - `text` (String): 文本内容。单条文本建议不超过 255 字符，系统支持更长文本但需符合 DNS 协议。
- **示例**:
  ```json
  { "text": "v=spf1 include:_spf.google.com ~all", "ttl": 3600 }
  ```

### 6. SRV 记录 (服务资源)
- **说明**: 用于定义特定服务的服务器位置（如 LDAP、SIP）。
- **字段**:
  - `priority` (Uint16): 优先级 (0-65535)。
  - `weight` (Uint16): 权重 (0-65535)。
  - `port` (Uint16): 服务端口 (1-65535)。
  - `target` (String): 目标主机名。
- **示例**:
  ```json
  {
    "priority": 10,
    "weight": 5,
    "port": 5060,
    "target": "sipserver.example.com.",
    "ttl": 3600
  }
  ```

### 7. CAA 记录 (证书颁发机构授权)
- **说明**: 限制哪些 CA 可以为该域名签发证书。
- **字段**:
  - `flag` (Uint8): 标志，通常为 `0` 或 `128` (Critical)。
  - `tag` (String): 标签，可选值：`issue`, `issuewild`, `iodef`, `contactphone`, `contactemail`。
  - `value` (String): 对应的值。
- **示例**:
  ```json
  { "flag": 0, "tag": "issue", "value": "letsencrypt.org", "ttl": 3600 }
  ```

### 8. SOA 记录 (起始授权机构)
- **说明**: 每个区域（Zone）有且仅能有一个有效的记录集包含一条 SOA 记录，通常位于区域根部（@）。
- **字段**:
  - `mname` (String): 该区域的主域名服务器。
  - `rname` (String): 该区域管理员的邮箱（将 `@` 替换为 `.`）。
  - `serial` (Uint32): 序列号，每次更改后递增（常用格式 YYYYMMDDNN）。
  - `refresh` (Uint32): 辅助服务器请求更新的时间间隔（秒）。
  - `retry` (Uint32): 更新失败后的重试间隔（秒）。
  - `expire` (Uint32): 辅助服务器认定区域数据失效的时间（秒）。
  - `minimum` (Uint32): 负缓存 TTL。
- **示例**:
  ```json
  {
    "mname": "ns1.console-dns.com.",
    "rname": "admin.example.com.",
    "serial": 2026013001,
    "refresh": 7200,
    "retry": 3600,
    "expire": 1209600,
    "minimum": 3600,
    "ttl": 3600
  }
  ```

---

## 错误代码参考

| 代码 | 描述 |
| :--- | :--- |
| `ZONE_NOT_FOUND` | 目标区域不存在 |
| `DUPLICATE_RECORD` | 记录已存在，无法重复创建 |
| `INVALID_FQDN` | 域名格式不符合 RFC 1035 规范 |
| `IP_WHITELIST_DENY` | 客户端 IP 不在白名单范围内 |
| `INSUFFICIENT_PERMISSIONS` | 权限不足以执行此写操作 |