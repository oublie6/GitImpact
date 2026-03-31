# 认证与用户体系
支持 `auth.mode=config|db|mixed`。
- config: 仅配置用户登录。
- db: 仅数据库用户登录，支持注册。
- mixed: 登录顺序**先 config_users，再 DB**；注册仅写 DB。

约束：
- 使用 JWT Bearer Token，含过期时间。
- disabled 用户禁止登录。
- 密码仅存 `password_hash`。
- `password_plain` 仅用于开发模式并需显式 `allow_dev_plain=true`。
