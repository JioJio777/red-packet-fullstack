# 接口设计文档

Base URL：`http://localhost:8080/api`  
数据格式：JSON  
认证方式：JWT，需要认证的接口在 Header 中携带 `Authorization: Bearer <token>`

---

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| code | 含义 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录 / Token 无效 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 1001 | 余额不足 |
| 1002 | 红包已抢完 |
| 1003 | 红包已过期 |
| 1004 | 已领取过该红包 |

---

## 一、认证模块

### 1.1 注册

`POST /auth/register`  
无需认证

**请求体：**
```json
{
  "username": "alice",
  "password": "123456"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice"
  }
}
```

---

### 1.2 登录

`POST /auth/login`  
无需认证

**请求体：**
```json
{
  "username": "alice",
  "password": "123456"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGci...",
    "expires_in": 86400
  }
}
```

---

## 二、用户模块

### 2.1 获取个人信息

`GET /user/profile`  
需要认证

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "balance": 10000
  }
}
```

> `balance` 单位为分，10000 = 100.00 元

---

## 三、红包模块

### 3.1 发红包

`POST /red-packets`  
需要认证

**请求体：**
```json
{
  "type": 1,
  "total_amount": 1000,
  "total_count": 5
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| type | int | 1=普通红包（每人等额），2=拼手气红包（随机） |
| total_amount | int | 总金额，单位：分 |
| total_count | int | 红包个数 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 100,
    "type": 1,
    "total_amount": 1000,
    "total_count": 5,
    "expired_at": "2026-02-20T10:00:00Z"
  }
}
```

---

### 3.2 领红包

`POST /red-packets/:id/claim`  
需要认证

**路径参数：** `id` — 红包ID

**请求体：** 无

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "amount": 200
  }
}
```

> `amount` 为本次领取的金额（单位：分）

---

### 3.3 查看红包详情

`GET /red-packets/:id`  
需要认证

用于渲染红包详情页的基本信息区域，包含当前用户的领取状态。

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 100,
    "sender_id": 1,
    "sender_name": "alice",
    "type": 1,
    "total_amount": 1000,
    "total_count": 5,
    "remaining_amount": 800,
    "remaining_count": 4,
    "claimed_count": 1,
    "status": 1,
    "expired_at": "2026-02-20T10:00:00Z",
    "created_at": "2026-02-19T10:00:00Z",
    "my_claim": {
      "claimed": true,
      "amount": 200,
      "claimed_at": "2026-02-19T10:05:00Z"
    }
  }
}
```

| 字段 | 说明 |
|------|------|
| claimed_count | 已领取人数 |
| my_claim.claimed | 当前用户是否已领取 |
| my_claim.amount | 当前用户领取金额，未领取时不返回 |
| my_claim.claimed_at | 当前用户领取时间，未领取时不返回 |

---

### 3.4 查看红包领取记录（分页）

`GET /red-packets/:id/records`  
需要认证

用于渲染详情页领取列表区域，支持分页。

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 10，最大 50 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "receiver_id": 2,
        "receiver_name": "bob",
        "amount": 200,
        "claimed_at": "2026-02-19T10:05:00Z"
      }
    ]
  }
}
```

---

### 3.5 查看我发出的红包

`GET /user/red-packets/sent`  
需要认证

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 10 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 3,
    "list": [
      {
        "id": 100,
        "type": 1,
        "total_amount": 1000,
        "total_count": 5,
        "remaining_count": 4,
        "status": 1,
        "created_at": "2026-02-19T10:00:00Z"
      }
    ]
  }
}
```

---

### 3.6 查看我收到的红包

`GET /user/red-packets/received`  
需要认证

**查询参数：** 同 3.5

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "red_packet_id": 100,
        "sender_name": "alice",
        "amount": 200,
        "created_at": "2026-02-19T10:05:00Z"
      }
    ]
  }
}
```

---

## 接口汇总

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /auth/register | 注册 | 否 |
| POST | /auth/login | 登录 | 否 |
| GET | /user/profile | 获取个人信息 | 是 |
| POST | /red-packets | 发红包 | 是 |
| POST | /red-packets/:id/claim | 领红包 | 是 |
| GET | /red-packets/:id | 红包详情（含当前用户领取状态） | 是 |
| GET | /red-packets/:id/records | 领取记录（分页） | 是 |
| GET | /user/red-packets/sent | 我发出的红包 | 是 |
| GET | /user/red-packets/received | 我收到的红包 | 是 |
