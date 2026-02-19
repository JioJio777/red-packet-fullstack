# 前端开发计划 · 登录 & 注册模块

## 一、技术选型 & 依赖

| 类目 | 选型 | 说明 |
|------|------|------|
| 框架 | React 18 + JavaScript | 轻量上手，无需类型注解 |
| 构建工具 | Vite | 比 CRA 启动更快，配置简单 |
| 路由 | React Router v6 | 声明式路由，支持嵌套路由和守卫 |
| HTTP 客户端 | Axios | 统一拦截器处理 token 注入 & 错误 |
| UI 组件库 | Ant Design 5 | 现成表单、校验、消息提示组件，开发快 |
| 状态管理 | React Context + useReducer | 项目体量不大，无需 Redux |
| 表单校验 | Ant Design Form 内置 | 配合 AntD Form 组件即可，不额外引入 |

---

## 二、目录结构

```
frontend/
├── public/
├── src/
│   ├── api/
│   │   ├── client.js          # Axios 实例，含拦截器
│   │   └── auth.js            # 登录、注册接口封装
│   ├── context/
│   │   └── AuthContext.jsx    # 全局认证状态（token、user）
│   ├── hooks/
│   │   └── useAuth.js         # 读取 AuthContext 的快捷 hook
│   ├── pages/
│   │   ├── LoginPage.jsx      # 登录页
│   │   └── RegisterPage.jsx   # 注册页
│   ├── components/
│   │   └── ProtectedRoute.jsx # 路由守卫（未登录跳转 /login）
│   ├── router/
│   │   └── index.jsx          # 路由配置
│   ├── App.jsx
│   └── main.jsx
├── index.html
└── package.json
```

---

## 三、API 层设计

### 3.1 Axios 实例（`api/client.js`）

- 创建带 `baseURL: http://localhost:8080/api` 的 axios 实例
- **请求拦截器**：从 localStorage 读取 token，自动附加 `Authorization: Bearer <token>` header
- **响应拦截器**：
  - 统一解包 `{ code, message, data }` 格式
  - code 不为 0 时抛出带 message 的错误
  - code 为 401 时清除本地 token 并跳转到 `/login`

### 3.2 认证接口封装（`api/auth.js`）

```js
// 登录：POST /auth/login → 返回 { token, expires_in }
login(username, password)

// 注册：POST /auth/register → 返回 { id, username }
register(username, password)
```

---

## 四、状态管理设计（`context/AuthContext.jsx`）

### 状态结构
```js
{
  token: null,          // string 或 null
  isAuthenticated: false
}
```

### 行为
| Action | 说明 |
|--------|------|
| `LOGIN` | 保存 token 到 state + localStorage |
| `LOGOUT` | 清除 state + localStorage |

### 初始化
- 应用启动时从 localStorage 读取 token，若存在则初始化为已登录状态。

---

## 五、页面设计

### 5.1 登录页（`/login`）

**布局：** 居中卡片，宽度约 400px

**字段：**
- 用户名（必填，长度 1-50）
- 密码（必填，长度 6-50，密码类型输入框）

**交互逻辑：**
1. 提交时调用 `api/auth.login`
2. 成功 → dispatch LOGIN action → 跳转到首页（`/`）
3. 失败 → 展示后端返回的 `message`（用 AntD message.error）
4. 已登录用户访问 `/login` → 直接跳转到 `/`

**辅助链接：** 页面底部有「还没有账号？去注册」跳转到 `/register`

---

### 5.2 注册页（`/register`）

**布局：** 与登录页同款居中卡片

**字段：**
- 用户名（必填，长度 1-50）
- 密码（必填，长度 6-50）
- 确认密码（必填，需与密码一致）

**交互逻辑：**
1. 提交时调用 `api/auth.register`
2. 成功 → 展示成功提示 → 自动跳转到 `/login`
3. 失败 → 展示后端返回的 `message`（例如「用户名已存在」）

**辅助链接：** 底部「已有账号？去登录」

---

## 六、路由配置（`router/index.jsx`）

| 路径 | 组件 | 是否受保护 |
|------|------|-----------|
| `/login` | LoginPage | 否（已登录则重定向到 `/`） |
| `/register` | RegisterPage | 否 |
| `/` | 首页（暂时占位） | 是（未登录跳转 `/login`） |

**路由守卫（`ProtectedRoute`）：** 读取 AuthContext，若 `isAuthenticated` 为 false，直接 `<Navigate to="/login" replace />` 。

---

## 七、开发顺序

1. **搭脚手架** — Vite + React + JavaScript 初始化，安装依赖
2. **API 层** — 实现 axios 实例 + auth 接口封装
3. **AuthContext** — 实现全局认证状态
4. **路由 & 守卫** — 配置路由，实现 ProtectedRoute
5. **登录页** — 实现 UI + 接口对接
6. **注册页** — 实现 UI + 接口对接
7. **联调验证** — 启动前后端，跑通注册→登录→跳转首页全流程

---

## 八、待确认事项

- [x] 首页（`/`）此阶段是展示一个占位页还是已有设计？→ 占位页即可
- [x] 登录成功后的跳转目标是否为 `/` ？→ 是，跳转到占位首页
- [x] Token 过期策略：是否需要前端自动检测过期并提示重登？→ 不加，依赖 401 拦截器处理即可
- [x] UI 风格：AntD 默认蓝色主题是否可接受？→ 可以
