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
