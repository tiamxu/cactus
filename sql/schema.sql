CREATE DATABASE cmdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

--
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `code` varchar(50) NOT NULL,
  `type` varchar(255) NOT NULL,
  `parentId` int(11) DEFAULT NULL,
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路由路径，默认为空字符串',
  `redirect` varchar(100) NOT NULL DEFAULT '' COMMENT '重定向路径，默认为空字符串',
  `icon` varchar(50) NOT NULL DEFAULT '' COMMENT '图标类名，默认为空字符串',
  `component` varchar(100) NOT NULL DEFAULT '' COMMENT '前端组件路径，默认为空字符串',
  `layout` varchar(50) NOT NULL DEFAULT '' COMMENT '布局组件，默认为空字符串',
  `keepAlive` tinyint(4) DEFAULT NULL,
  `method` varchar(10) NOT NULL DEFAULT '' COMMENT 'HTTP方法(GET/POST等)，默认为空字符串',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '权限描述，默认为空字符串',
  `show` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否展示在页面菜单',
  `enable` tinyint(4) NOT NULL DEFAULT '1',
  `order` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_30e166e8c6359970755c5727a2` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

--
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
INSERT INTO `permission` VALUES (1,'资源管理','Resource_Mgt','MENU',2,'/pms/resource','','i-fe:list','/src/views/pms/resource/index.vue','',0,'','',1,1,1),(2,'系统管理','SysMgt','MENU',NULL,'','','i-fe:grid','','',0,'','',1,1,2),(3,'角色管理','RoleMgt','MENU',2,'/pms/role','','i-fe:user-check','/src/views/pms/role/index.vue','',0,'','',1,1,2),(4,'用户管理','UserMgt','MENU',2,'/pms/user','','i-fe:user','/src/views/pms/user/index.vue','',1,'','',1,1,3),(5,'分配用户','RoleUser','MENU',3,'/pms/role/user/:roleId','','i-fe:user-plus','/src/views/pms/role/role-user.vue','',0,'','',0,1,1),(6,'业务示例','Demo','MENU',NULL,'','','i-fe:grid','','',0,'','',1,1,1),(8,'个人资料','UserProfile','MENU',NULL,'/profile','','i-fe:user','/src/views/profile/index.vue','',0,'','',0,1,99),(9,'基础功能','Base','MENU',NULL,'/base','','i-fe:grid','','',0,'','',1,1,0),(13,'创建新用户','AddUser','BUTTON',4,'','','','','',0,'','',1,1,1),(21,'项目管理','Project','MENU',6,'','','','','',0,'','',1,1,0),(22,'测试菜单','Test','MENU',NULL,'','','','','',0,'','',1,1,0);
/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
UNLOCK TABLES;

--

--

DROP TABLE IF EXISTS `profile`;

CREATE TABLE `profile` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `gender` int(11) NOT NULL DEFAULT '0' COMMENT '性别，默认值为 0',
  `avatar` varchar(255) NOT NULL DEFAULT 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80',
  `address` varchar(255) NOT NULL DEFAULT '' COMMENT '地址，默认值为空字符串',
  `email` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱，默认值为空字符串',
  `userId` int(11) NOT NULL,
  `nickName` varchar(10) NOT NULL DEFAULT '匿名' COMMENT '昵称，默认值为 "匿名"',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_a24972ebd73b106250713dcddd` (`userId`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `profile`
--

LOCK TABLES `profile` WRITE;
/*!40000 ALTER TABLE `profile` DISABLE KEYS */;
INSERT INTO `profile` VALUES (1,1,'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80','深圳','123@qq.com',1,'Admin'),(2,1,'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80','未知','1@qq.com',3,'test');
/*!40000 ALTER TABLE `profile` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;

CREATE TABLE `role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `name` varchar(50) NOT NULL,
  `enable` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_ee999bb389d7ac0fd967172c41` (`code`) USING BTREE,
  UNIQUE KEY `IDX_ae4578dcaed5adff96595e6166` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;


--

LOCK TABLES `role` WRITE;
INSERT INTO `role` VALUES (1,'SUPER_ADMIN','超级管理员',1),(2,'ROLE_QA','质检员',1);
UNLOCK TABLES;

--
--

DROP TABLE IF EXISTS `role_permissions_permission`;

CREATE TABLE `role_permissions_permission` (
  `roleId` int(11) NOT NULL,
  `permissionId` int(11) NOT NULL,
  PRIMARY KEY (`roleId`,`permissionId`) USING BTREE,
  KEY `IDX_b36cb2e04bc353ca4ede00d87b` (`roleId`) USING BTREE,
  KEY `IDX_bfbc9e263d4cea6d7a8c9eb3ad` (`permissionId`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

--

--

LOCK TABLES `role_permissions_permission` WRITE;
INSERT INTO `role_permissions_permission` VALUES (2,1),(2,2),(2,3),(2,4),(2,5),(2,6),(2,9);
UNLOCK TABLES;
--

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `enable` tinyint(4) NOT NULL DEFAULT '1',
  `createTime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updateTime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_78a916df40e02a9deb1c4b75ed` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;


--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'admin','$2a$10$TeseeIS/.56V.os9QvXf.OOxlpsJbUeLwUoQc.YGZ.ERsT062Pe7W',1,'2023-11-18 16:18:59.150632','2025-04-23 23:08:48.305310'),(3,'test','$2a$10$A4Gs8bR.epRWHdWi3T2XbOBEozmub2EFYwY2xRfKhovdapq4dQOrC',1,'2025-04-22 20:17:20.707722','2025-04-22 20:17:20.707722');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;


--

DROP TABLE IF EXISTS `user_roles_role`;

CREATE TABLE `user_roles_role` (
  `userId` int(11) NOT NULL,
  `roleId` int(11) NOT NULL,
  PRIMARY KEY (`userId`,`roleId`) USING BTREE,
  KEY `IDX_5f9286e6c25594c6b88c108db7` (`userId`) USING BTREE,
  KEY `IDX_4be2f7adf862634f5f803d246b` (`roleId`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

--
--

LOCK TABLES `user_roles_role` WRITE;
INSERT INTO `user_roles_role` VALUES (1,1),(1,2),(3,1);
UNLOCK TABLES;

