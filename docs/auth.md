# 认证与用户体系
支持 `auth.mode=config|db|mixed`。
- config: 仅配置用户登录。
- db: 仅数据库用户登录，支持注册。
- mixed: 登录顺序**先 config_users，再 DB**；注册仅写 DB。

约束：
- 使用 JWT Bearer Token，含过期时间。
- disabled 用户禁止登录。
- 配置用户（`config_users`）使用明文 `password` 字段校验。
- 数据库用户密码仍存 `password_hash`，并通过 bcrypt 校验。
- `mixed` 模式登录顺序固定为：先 `config_users`，后 DB。
