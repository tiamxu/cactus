-- 创建数据库：
CREATE DATABASE cmdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `username` varchar(128) DEFAULT NULL COMMENT '''用户名''',
  `password` varchar(128) DEFAULT NULL COMMENT '''用户密码''',
  `phone` varchar(11) DEFAULT NULL COMMENT '''手机号码''',
  `email` varchar(128) DEFAULT NULL COMMENT '''邮箱''',
  `nick_name` varchar(128) DEFAULT NULL COMMENT '''用户昵称''',
  `avatar` varchar(128) DEFAULT 'http://dnsjia.com/img/avatar.png' COMMENT '''用户头像''',
  `status` tinyint(1) DEFAULT '1' COMMENT '''用户状态(正常/禁用, 默认正常)''',
  `role_id` bigint(20) unsigned DEFAULT NULL COMMENT '''角色id外键''',
  `dept_id` bigint(20) unsigned DEFAULT NULL COMMENT '''部门id外键''',
  `uid` bigint(20) DEFAULT NULL COMMENT '''用戶uid''',
  `create_by` varchar(191) DEFAULT NULL COMMENT '''创建来源''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(128) DEFAULT NULL COMMENT '''角色名称''',
  `desc` varchar(128) DEFAULT NULL COMMENT '''角色描述''',
  PRIMARY KEY (`id`),
  KEY `idx_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL COMMENT '''菜单名称''',
  `icon` varchar(64) DEFAULT NULL COMMENT '''菜单图标''',
  `path` varchar(64) DEFAULT NULL COMMENT '''菜单访问路径''',
  `sort` int(3) DEFAULT '0' COMMENT '''菜单顺序(同级菜单, 从0开始, 越小显示越靠前)''',
  `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '''父菜单编号(编号为0时表示根菜单)''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  PRIMARY KEY (`id`),
  KEY `idx_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `projects`;
CREATE TABLE `projects` (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(50) NOT NULL NOT NULL DEFAULT '' COMMENT '项目名',
    description TEXT NOT NULL DEFAULT '' COMMENT '描述',
    status      VARCHAR(50) NOT NULL DEFAULT 'active' COMMENT '状态' ,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tags        TEXT[]
);

CREATE INDEX idx_projects_name ON projects(name);
CREATE INDEX idx_projects_status ON projects(status);


CREATE TABLE IF NOT EXISTS navigations (
    -- 主键ID，自增整数，唯一标识每条导航链接记录
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键ID',
    -- 链接显示标题，必填项，用于前端展示的链接名称
    title VARCHAR(100) NOT NULL COMMENT '链接标题/名称',
    -- 链接的实际URL地址，必填项，用户点击后跳转的目标地址
    url VARCHAR(255) NOT NULL COMMENT '链接URL地址',
    -- 链接关联的图标类名，通常用于Font Awesome等图标库
    -- 示例: 'fa-home', 'fa-github'，可为空表示不使用图标
    icon VARCHAR(50) COMMENT '图标CSS类名',
    -- 链接分类，用于对导航链接进行分组管理
    -- 示例: '开发工具'、'搜索引擎'、'社交媒体'
    category VARCHAR(50) COMMENT '链接分类',
    -- 链接的详细描述信息，用于鼠标悬停提示或辅助说明
    description TEXT COMMENT '链接描述信息',
    -- 记录创建时间，自动设置为记录插入时的时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    -- 记录最后更新时间，自动更新为记录修改时的时间戳
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间' 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='快捷导航链接表';