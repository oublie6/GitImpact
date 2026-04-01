# 数据库初始化
## MySQL
1. 创建库并执行：`sql/mysql/init.sql`
2. 配置 `database.type=mysql` 与 MySQL DSN
3. 后端驱动：`gorm.io/driver/mysql`

## 达梦（DM）
1. 在目标 schema 执行：`sql/dameng/init.sql`
2. 配置 `database.type=dameng`
3. 后端驱动：`github.com/godoes/gorm-dameng`（GORM Dialector）
4. 连接方式：
   - 若 `database.dsn` 非空，直接使用该值（示例：`dm://SYSDBA:SYSDBA@127.0.0.1:5236?schema=SYSDBA`）。
   - 若 `database.dsn` 为空，则按 `database.dameng.host/port/user/password/dbname` 自动拼接，
     其中 `dbname` 映射为 `schema` 参数。

## 差异说明
- 自增：MySQL `AUTO_INCREMENT`，达梦 `IDENTITY`。
- 大文本：MySQL `LONGTEXT`，达梦 `CLOB`。
- 时间字段默认值语法存在差异，脚本已分别提供。
