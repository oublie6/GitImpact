# 数据库初始化
## MySQL
1. 创建库并执行：`sql/mysql/init.sql`
2. 配置 `database.type=mysql` 与 MySQL DSN

## 达梦（DM）
1. 在目标 schema 执行：`sql/dameng/init.sql`
2. 配置 `database.type=dameng`
3. 当前兼容策略：使用 MySQL 风格 DSN/语义适配（首版策略），生产建议替换为专用驱动与方言。

## 差异说明
- 自增：MySQL `AUTO_INCREMENT`，达梦 `IDENTITY`。
- 大文本：MySQL `LONGTEXT`，达梦 `CLOB`。
- 时间字段默认值语法存在差异，脚本已分别提供。
