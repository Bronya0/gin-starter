--用户表
CREATE TABLE IF NOT EXISTS auth_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,  --'密码哈希值'
    mobile VARCHAR(64),
    email VARCHAR(64),
    avatar VARCHAR(255),
    description VARCHAR(500),
    sex INT,
    status INT DEFAULT 1,  --'状态(0:禁用;1:正常)'
    last_ip INET,
    last_login_time TIMESTAMP,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_auth_user_status ON auth_user(status);
CREATE INDEX IF NOT EXISTS idx_auth_user_mobile ON auth_user(mobile);


--会话session
CREATE TABLE IF NOT EXISTS auth_session (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(512) NOT NULL,        -- 用户登录的 Token
    device_type VARCHAR(50),                   -- 设备类型 (web, mobile, etc.)
    ip_address INET,                           -- 登录时的 IP 地址
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_time TIMESTAMP NOT NULL            -- Token 过期时间
);
CREATE INDEX IF NOT EXISTS idx_auth_session_user ON auth_session(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_session_expires_time ON auth_session(expires_time);


-- 角色表
CREATE TABLE IF NOT EXISTS auth_role (
    id SERIAL PRIMARY KEY,
    role_name varchar(50) NOT NULL,
    role_code varchar(50) NOT NULL,  --角色编码
    status INT DEFAULT 1,  --'状态(0:禁用;1:正常)'
    description VARCHAR(500),
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth_user_role (
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    UNIQUE(user_id, role_id)
);


--权限表
CREATE TABLE auth_permission (
    id SERIAL PRIMARY KEY,
    parent_id INT,
    name VARCHAR(64) NOT NULL,
    path VARCHAR(256),  --路径
    component varchar(200) DEFAULT NULL,  --组件路径
    type INT NOT NULL,  --类型(1:菜单;2:按钮;3:接口)
    permission varchar(100) DEFAULT NULL,  --权限标识
    method VARCHAR(20),  -- HTTP方法
    sort INT DEFAULT 0,  --排序
    status INT DEFAULT 1,  --'状态(0:禁用;1:正常)'
    icon varchar(100)  --图标
);
CREATE INDEX IF NOT EXISTS idx_auth_permission_parent ON auth_permission(parent_id);
CREATE INDEX IF NOT EXISTS idx_auth_permission_type ON auth_permission(type);
CREATE INDEX IF NOT EXISTS idx_auth_permission_status ON auth_permission(status);

CREATE TABLE IF NOT EXISTS auth_role_permission (
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    UNIQUE(role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS auth_user_permission (
    user_id INT NOT NULL,
    permission_id INT NOT NULL,
    UNIQUE(user_id, permission_id)
);


-- 操作日志表
CREATE TABLE IF NOT EXISTS auth_operate_log (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    action VARCHAR(100) NOT NULL,  -- 操作类型
    target_type VARCHAR(50),       -- 操作对象类型
    target_id INT,                 -- 操作对象ID
    ip_address INET,
    user_agent VARCHAR(255),       -- 用户代理
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);