# 数据库表设计

数据库：MySQL  
字符集：utf8mb4  
存储引擎：InnoDB

---

## 1. 用户表 `users`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 用户ID |
| username | VARCHAR(50) | NOT NULL, UNIQUE | 用户名 |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt 加密后的密码 |
| balance | BIGINT UNSIGNED | NOT NULL, DEFAULT 0 | 账户余额（单位：分） |
| created_at | DATETIME | NOT NULL | 创建时间 |
| updated_at | DATETIME | NOT NULL | 更新时间 |

> 余额使用整数（分）存储，避免浮点数精度问题。

---

## 2. 红包表 `red_packets`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 红包ID |
| sender_id | BIGINT UNSIGNED | NOT NULL, FK → users.id | 发送者ID |
| type | TINYINT | NOT NULL | 红包类型：1=普通红包，2=拼手气红包 |
| total_amount | BIGINT UNSIGNED | NOT NULL | 红包总金额（单位：分） |
| total_count | INT UNSIGNED | NOT NULL | 红包总个数 |
| remaining_amount | BIGINT UNSIGNED | NOT NULL | 剩余金额（单位：分） |
| remaining_count | INT UNSIGNED | NOT NULL | 剩余个数 |
| status | TINYINT | NOT NULL, DEFAULT 1 | 状态：1=可领取，2=已抢完，3=已过期 |
| expired_at | DATETIME | NOT NULL | 过期时间（默认发出后 24 小时） |
| created_at | DATETIME | NOT NULL | 创建时间 |

**索引：**
- `idx_sender_id`：sender_id（查询我发出的红包）
- `idx_status_expired_at`：status, expired_at（过期扫描）

---

## 3. 领取记录表 `red_packet_records`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 记录ID |
| red_packet_id | BIGINT UNSIGNED | NOT NULL, FK → red_packets.id | 红包ID |
| receiver_id | BIGINT UNSIGNED | NOT NULL, FK → users.id | 领取者ID |
| amount | BIGINT UNSIGNED | NOT NULL | 本次领取金额（单位：分） |
| created_at | DATETIME | NOT NULL | 领取时间 |

**索引：**
- `uk_packet_receiver`：(red_packet_id, receiver_id) UNIQUE（防止同一用户重复领同一红包）
- `idx_receiver_id`：receiver_id（查询我收到的红包）

---

## 4. 流水表 `transactions`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 流水ID |
| user_id | BIGINT UNSIGNED | NOT NULL, FK → users.id | 用户ID |
| type | VARCHAR(20) | NOT NULL | 类型：recharge / send / receive / refund |
| direction | TINYINT | NOT NULL | 资金方向：1=收入，2=支出 |
| amount | BIGINT UNSIGNED | NOT NULL | 变动金额（单位：分，恒为正数） |
| balance_after | BIGINT UNSIGNED | NOT NULL | 变动后余额（单位：分） |
| related_id | BIGINT UNSIGNED | NULL | 关联红包ID（非红包场景为 NULL） |
| remark | VARCHAR(255) | NULL | 备注 |
| created_at | DATETIME | NOT NULL | 创建时间 |

| type 值 | direction | 场景 |
|---------|-----------|------|
| recharge | 1（收入） | 充值 |
| send | 2（支出） | 发红包扣款 |
| receive | 1（收入） | 领红包到账 |
| refund | 1（收入） | 红包过期退款 |

**索引：**
- `idx_user_id_created_at`：(user_id, created_at)（查询个人流水，按时间排序）
- `idx_related_id`：related_id（按红包ID反查流水）

---

## ER 关系

```
users  ──< red_packets        (一个用户可发多个红包)
users  ──< red_packet_records (一个用户可领多个红包)
users  ──< transactions       (一个用户有多条流水)
red_packets ──< red_packet_records (一个红包可被多人领取)
red_packets ──< transactions       (一个红包对应多条流水)
```

---

## 并发安全说明

抢红包是典型的高并发写场景，核心思路：

1. **Redis 预占**：用 `DECR` 原子操作扣减 Redis 中的剩余个数，抢到名额再写 MySQL
2. **MySQL 事务**：更新 `remaining_amount`、`remaining_count`，插入 `red_packet_records`，更新 `users.balance` 在同一事务内完成
3. **唯一索引兜底**：`uk_packet_receiver` 防止并发场景下重复写入
