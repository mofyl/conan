CREATE DATABASE platform DEFAULT CHARSET utf8;

/* 第三方平台信息表 */
CREATE TABLE `platform`.`platform_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `memo` varchar(256) NOT NULL COMMENT '简介',
  `app_key` varchar(32) NOT NULL COMMENT '唯一标识',
  `secret_key` varchar(32) NOT NULL COMMENT '私钥',
  `index_url` varchar(256) NOT NULL COMMENT '主页url',
  `channel_name` varchar(32) NOT NULL COMMENT '平台名字',
  `sign_url` varchar(256) NOT NULL COMMENT '回调url',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `deleted` int(1) NOT NULL COMMENT '删除标记 1 表示未删除  2表示删除',
  `deleted_time` bigint(20) COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`app_key`),
  KEY `idx_channel_name` (`channel_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* MM用户授权表 表示用于对该三方应用进行授权，该三方应用可以获取MM用户的基本信息  */
CREATE TABLE `platform`.`mm_user_auth_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_key` varchar(32) NOT NULL COMMENT'三方平台唯一标识',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `mm_user_id` varchar(256) NOT NULL COMMENT 'mm授权用户的ID',
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


/* 管理平台 管理员表  */
CREATE TABLE `platform`.`admin_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `login_name` varchar(32) NOT NULL COMMENT '三方平台管理员的登录名',
  `login_pwd` varchar(128) NOT NULL COMMENT '三方平台管理员的密码 这里使用bcrypt编码',
  `name` varchar(32) NOT NULL COMMENT '用户名',
  `number` varchar(32) NOT NULL COMMENT '工号',
  `created` bigint(32) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  key idx_login_name(`login_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 权限中的角色 */
CREATE TABLE `platform_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(32) NOT NULL COMMENT '角色的名字',
  `role_sign` varchar(32) NOT NULL COMMENT '角色的标识 若该角色可访问多个平台，则这些数据的 role_sign 相同',
  `platform_app_key` varchar(256) NOT NULL COMMENT '该角色可访问的平台appkey',
  `create_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`platform_app_key`),
  KEY `idx_role_sign` (`role_sign`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `mm_user_platform_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mm_user_id` varchar(32) NOT NULL COMMENT 'mm用户的标识 这里可能使用userID',
  `role_sign` varchar(256) NOT NULL COMMENT '该用户拥有的角色ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;








