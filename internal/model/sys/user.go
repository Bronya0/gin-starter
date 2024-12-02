package sys

import "time"

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type AuthUser struct {
	Id            string    `gorm:"column:id;type:SERIAL;primaryKey;" json:"id"`
	Username      string    `gorm:"column:username;type:VARCHAR(64);not null;" json:"username"`
	Password      string    `gorm:"column:password;type:VARCHAR(255);not null;" json:"password"`
	Mobile        string    `gorm:"column:mobile;type:VARCHAR(64);" json:"mobile"`
	Email         string    `gorm:"column:email;type:VARCHAR(64);" json:"email"`
	Avatar        string    `gorm:"column:avatar;type:VARCHAR(255);" json:"avatar"`
	Description   string    `gorm:"column:description;type:VARCHAR(500);" json:"description"`
	Sex           int32     `gorm:"column:sex;type:INT;" json:"sex"`
	Status        int32     `gorm:"column:status;type:INT;default:1;" json:"status"`
	LastIp        string    `gorm:"column:last_ip;type:INET;" json:"last_ip"`
	LastLoginTime time.Time `gorm:"column:last_login_time;type:TIMESTAMP;" json:"last_login_time"`
	CreateTime    time.Time `gorm:"column:create_time;type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`
	IsDeleted     bool      `gorm:"column:is_deleted;type:BOOLEAN;default:false;" json:"is_deleted"`
}

func (m *AuthUser) TableName() string {
	return "auth_user" // return "schema.table"
}

type AuthSession struct {
	Id          string    `gorm:"column:id;type:SERIAL;primaryKey;" json:"id"`
	UserId      int32     `gorm:"column:user_id;type:INT;not null;" json:"user_id"`
	Token       string    `gorm:"column:token;type:VARCHAR(512);not null;" json:"token"`
	DeviceType  string    `gorm:"column:device_type;type:VARCHAR(50);" json:"device_type"`
	IpAddress   string    `gorm:"column:ip_address;type:INET;" json:"ip_address"`
	CreateTime  time.Time `gorm:"column:create_time;type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`
	ExpiresTime time.Time `gorm:"column:expires_time;type:TIMESTAMP;not null;" json:"expires_time"`
}

func (m *AuthSession) TableName() string {
	return "auth_session" // return "schema.table"
}

type AuthRole struct {
	Id          string    `gorm:"column:id;type:SERIAL;primaryKey;" json:"id"`
	RoleName    string    `gorm:"column:role_name;type:varchar(50);not null;" json:"role_name"`
	RoleCode    string    `gorm:"column:role_code;type:varchar(50);not null;" json:"role_code"`
	Status      int32     `gorm:"column:status;type:INT;default:1;" json:"status"`
	Description string    `gorm:"column:description;type:VARCHAR(500);" json:"description"`
	CreateTime  time.Time `gorm:"column:create_time;type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`
}

func (m *AuthRole) TableName() string {
	return "auth_role" // return "schema.table"
}

type AuthUserRole struct {
	UserId int32 `gorm:"column:user_id;type:INT;not null;" json:"user_id"`
	RoleId int32 `gorm:"column:role_id;type:INT;not null;" json:"role_id"`
}

func (m *AuthUserRole) TableName() string {
	return "auth_user_role" // return "schema.table"
}

type AuthPermission struct {
	Id         string `gorm:"column:id;type:SERIAL;primaryKey;" json:"id"`
	ParentId   int32  `gorm:"column:parent_id;type:INT;" json:"parent_id"`
	Name       string `gorm:"column:name;type:VARCHAR(64);not null;" json:"name"`
	Path       string `gorm:"column:path;type:VARCHAR(256);" json:"path"`
	Component  string `gorm:"column:component;type:varchar(200);default:NULL;" json:"component"`
	Type       int32  `gorm:"column:type;type:INT;not null;" json:"type"`
	Permission string `gorm:"column:permission;type:varchar(100);default:NULL;" json:"permission"`
	Method     string `gorm:"column:method;type:VARCHAR(20);" json:"method"`
	Sort       int32  `gorm:"column:sort;type:INT;default:0;" json:"sort"`
	Status     int32  `gorm:"column:status;type:INT;default:1;" json:"status"`
	Icon       string `gorm:"column:icon;type:varchar(100);" json:"icon"`
}

func (m *AuthPermission) TableName() string {
	return "auth_permission" // return "schema.table"
}

type AuthRolePermission struct {
	RoleId       int32 `gorm:"column:role_id;type:INT;not null;" json:"role_id"`
	PermissionId int32 `gorm:"column:permission_id;type:INT;not null;" json:"permission_id"`
}

func (m *AuthRolePermission) TableName() string {
	return "auth_role_permission" // return "schema.table"
}

type AuthUserPermission struct {
	UserId       int32 `gorm:"column:user_id;type:INT;not null;" json:"user_id"`
	PermissionId int32 `gorm:"column:permission_id;type:INT;not null;" json:"permission_id"`
}

func (m *AuthUserPermission) TableName() string {
	return "auth_user_permission" // return "schema.table"
}

type AuthOperateLog struct {
	Id         string    `gorm:"column:id;type:SERIAL;primaryKey;" json:"id"`
	UserId     int32     `gorm:"column:user_id;type:INT;not null;" json:"user_id"`
	Action     string    `gorm:"column:action;type:VARCHAR(100);not null;" json:"action"`
	TargetType string    `gorm:"column:target_type;type:VARCHAR(50);" json:"target_type"`
	TargetId   int32     `gorm:"column:target_id;type:INT;" json:"target_id"`
	IpAddress  string    `gorm:"column:ip_address;type:INET;" json:"ip_address"`
	UserAgent  string    `gorm:"column:user_agent;type:VARCHAR(255);" json:"user_agent"`
	CreateTime time.Time `gorm:"column:create_time;type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`
}

func (m *AuthOperateLog) TableName() string {
	return "auth_operate_log" // return "schema.table"
}
