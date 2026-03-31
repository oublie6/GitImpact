-- 达梦兼容初始化脚本（需在 GITIMPACT schema 下执行）
CREATE TABLE users (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  username VARCHAR(64) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  display_name VARCHAR(128),
  email VARCHAR(128),
  role VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  source VARCHAR(32) NOT NULL,
  last_login_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_role ON users(role);

CREATE TABLE repositories (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  name VARCHAR(128) NOT NULL UNIQUE,
  repo_url VARCHAR(500) NOT NULL,
  default_branch VARCHAR(128) NOT NULL,
  local_cache_dir VARCHAR(500) NOT NULL,
  auth_note VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE analysis_tasks (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  task_name VARCHAR(200) NOT NULL,
  mode VARCHAR(64) NOT NULL,
  old_repo_id BIGINT,
  old_ref VARCHAR(128) NOT NULL,
  new_repo_id BIGINT,
  new_ref VARCHAR(128) NOT NULL,
  generate_markdown SMALLINT DEFAULT 1,
  generate_structured SMALLINT DEFAULT 1,
  custom_focus CLOB,
  remark VARCHAR(500),
  status VARCHAR(32) NOT NULL,
  error_message CLOB,
  created_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_tasks_status ON analysis_tasks(status);

CREATE TABLE analysis_reports (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  task_id BIGINT NOT NULL UNIQUE,
  markdown_report CLOB,
  structured_report CLOB,
  raw_stdout CLOB,
  raw_stderr CLOB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE task_logs (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  task_id BIGINT NOT NULL,
  level VARCHAR(16) NOT NULL,
  message CLOB NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE task_artifacts (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  task_id BIGINT NOT NULL,
  artifact_key VARCHAR(128) NOT NULL,
  file_path VARCHAR(500) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE system_settings (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  key VARCHAR(128) NOT NULL UNIQUE,
  value CLOB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users(username,password_hash,display_name,email,role,status,source)
VALUES('admin','$2a$10$8aMnM3zRGswyfWCGfFTWUO/Jci8hV3dqjhkPyp6Vn7I.xXvDO4Tm2','默认管理员','admin@example.com','admin','active','db');
