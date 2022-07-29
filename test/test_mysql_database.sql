
/*
main_info和sub_info为一对多关系, main_info的id字段关联sub_info的main_id字段
sub_info和dept_info为多对多关系,关联的中间表是sub_dept_rel,中间表关联的字段是sub_info和dept_info各自的id字段
dict_info是字典项信息表
*/

CREATE DATABASE IF  NOT EXISTS `test_mysql_database`;
USE `test_mysql_database`;

DROP TABLE IF EXISTS `main_info`;
CREATE TABLE `main_info` (
  `id` varchar(255) NOT NULL COMMENT '主键ID',
  `name` varchar(255) NOT NULL COMMENT '姓名',
  `age` int NOT NULL DEFAULT 0 COMMENT '年龄',
  `sex` tinyint NOT NULL DEFAULT 0 COMMENT '性别',
  `status` varchar(255) NOT NULL DEFAULT '888',
  `creator` varchar(64) NOT NULL  DEFAULT 'admin',
  `modifier` varchar(64) NOT NULL  DEFAULT 'admin',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `main_info` WRITE;
INSERT INTO `main_info` VALUES 
('419c98573158498ea318d1665303f10a','main_name_1',30,0,'888','admin','admin',now(),now(),NULL),
('d7cafc7da04148c7a3020ed57c45f5cc','main_name_2',32,1,'888','admin','admin',now(),now(),NULL),
('0c0e8f711a36446d9a012565d0304170','main_name_3',33,1,'888','admin','admin',now(),now(),NULL),
('f80a79bc356d444889a2b504acc7ee12','main_name_4',31,0,'888','admin','admin',now(),now(),NULL),
('3c7126b0b089405eac0df11c29dc0975','main_name_5',35,0,'888','admin','admin',now(),now(),NULL),
('b506cce19d60446ca43b0d3ef3cb67eb','main_name_6',31,1,'888','admin','admin',now(),now(),NULL);
UNLOCK TABLES;


DROP TABLE IF EXISTS `sub_info`;
CREATE TABLE `sub_info` (
  `id` varchar(255) NOT NULL COMMENT '主键ID',
  `main_id` varchar(255) NOT NULL COMMENT 'main表的主键ID',
  `name` varchar(255) NOT NULL COMMENT '姓名',
  `age` int NOT NULL COMMENT '年龄',
  `sex` tinyint NOT NULL DEFAULT '0' COMMENT '性别',
  `dict_field` varchar(4) NOT NULL DEFAULT '00' COMMENT '测试字典项',
  `field_1` varchar(255) NOT NULL COMMENT '字段1',
  `field_2` varchar(255) NOT NULL COMMENT '字段2',
  `field_3` varchar(255) NOT NULL COMMENT '字段3',
  `field_4` varchar(255) NOT NULL COMMENT '字段4',
  `status` varchar(255) NOT NULL DEFAULT '888',
  `creator` varchar(64) NOT NULL  DEFAULT 'admin',
  `modifier` varchar(64) NOT NULL  DEFAULT 'admin',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;


LOCK TABLES `sub_info` WRITE;
INSERT INTO `sub_info` VALUES
('13892d41d95a4061a7aa346692e8e888','419c98573158498ea318d1665303f10a','小明',12,0,'00','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('45a7c7baa6c6430a850f33f68f3f4eba','419c98573158498ea318d1665303f10a','小青',12,1,'00','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('6c7bcf4822e5437d8448b322f052ceee','0c0e8f711a36446d9a012565d0304170','小红',12,1,'01','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('a218a0a6afb3447a833533203e1cbdb4','f80a79bc356d444889a2b504acc7ee12','小萧',12,0,'01','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('b489e10af58c47898e8768d67b54a434','d7cafc7da04148c7a3020ed57c45f5cc','小微',12,1,'02','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('bd0c847f6d1143ec88bcbef375e749fe','d7cafc7da04148c7a3020ed57c45f5cc','小伟',12,0,'02','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('c5c096870b784305b79ac9aa339d5591','f80a79bc356d444889a2b504acc7ee12','Mark',12,0,'03','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('80e52c60f538495ea621a27ef3b245ed','3c7126b0b089405eac0df11c29dc0975','Judy',12,0,'03','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('127e86a82143494d9be270b3284cb783','b506cce19d60446ca43b0d3ef3cb67eb','Erin',12,1,'03','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('2ba7adb2b3cd44ea9ecf9fa10ae30b69','0c0e8f711a36446d9a012565d0304170','John',12,0,'00','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('6a938cd43f3d4d2fa31f3ca7834d57a2','3c7126b0b089405eac0df11c29dc0975','Bob',12,0,'02','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL),
('9c7d1a14713949f6b7c7039c333360d2','b506cce19d60446ca43b0d3ef3cb67eb','Tom',12,0,'01','test_field_separation_string','[\"field_separation_1\",1,\"field_separation_2\"]','{\"key_1\":1, \"key_2\":\"test_key_2\", \"key_3\":3}','[{\"Key_1\":1111,\"Key_2\":\"TEST_KEY_2\"},{\"Key_3\":\"TEST_key_3\",\"Key_4\":444},{\"Key_5\":555,\"Key_6\":\"TEST_KEY_6\"}]','888','admin','admin',now(),now(),NULL);
UNLOCK TABLES;


DROP TABLE IF EXISTS `dept_info`;
CREATE TABLE `dept_info`(
  `id` varchar(255) NOT NULL COMMENT '主键ID',
  `dept_name` varchar(255) NOT NULL COMMENT '所属部门',
  `phone` varchar(32) NOT NULL  COMMENT '电话',
  `dept_floor` int(10) NOT NULL COMMENT '部门所在楼层',
  `status` varchar(255) NOT NULL DEFAULT '888',
  `creator` varchar(64) NOT NULL  DEFAULT 'admin',
  `modifier` varchar(64) NOT NULL  DEFAULT 'admin',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `dept_info` WRITE;
INSERT INTO `dept_info` VALUES
('65648bd4926f44c3b35b426207fa6282','测试部门_1','123456',4,'888','admin','admin',now(),now(),NULL),
('87289b2334254c25b9ea6f97c213f82c','测试部门_2','123456',5,'888','admin','admin',now(),now(),NULL),
('5f3969b6b92645e3bf83e788dd9bd691','测试部门_3','123456',8,'888','admin','admin',now(),now(),NULL),
('b6051d4ae94b41e5ae66288edcb9803b','测试部门_4','123456',12,'888','admin','admin',now(),now(),NULL),
('77964202be5d4d8296eaaa7ef0bb595a','测试部门_5','123456',15,'888','admin','admin',now(),now(),NULL),
('8319e5c4385e46cbbda8d4355c9a0e55','测试部门_6','123456',9,'888','admin','admin',now(),now(),NULL),
('017a992c460b470da68f262deadc9f78','测试部门_7','123456',6,'888','admin','admin',now(),now(),NULL);
UNLOCK TABLES;


DROP TABLE IF EXISTS `sub_dept_rel`;
CREATE TABLE `sub_dept_rel`(
  `id` bigint NOT NULL COMMENT '主键ID',
  `sub_id` varchar(255) NOT NULL COMMENT 'sub_info表的主键ID',
  `dept_id` varchar(255) NOT NULL COMMENT 'dept_info表的主键ID',
  `status` varchar(255) NOT NULL DEFAULT '888',
  `creator` varchar(64) NOT NULL  DEFAULT 'admin',
  `modifier` varchar(64) NOT NULL  DEFAULT 'admin',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `sub_dept_rel` WRITE;
INSERT INTO `sub_dept_rel` VALUES(1,'13892d41d95a4061a7aa346692e8e888','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(2,'45a7c7baa6c6430a850f33f68f3f4eba','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(3,'6c7bcf4822e5437d8448b322f052ceee','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(4,'a218a0a6afb3447a833533203e1cbdb4','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(5,'b489e10af58c47898e8768d67b54a434','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(6,'bd0c847f6d1143ec88bcbef375e749fe','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(7,'c5c096870b784305b79ac9aa339d5591','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(8,'80e52c60f538495ea621a27ef3b245ed','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(9,'127e86a82143494d9be270b3284cb783','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(10,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(11,'6a938cd43f3d4d2fa31f3ca7834d57a2','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(12,'9c7d1a14713949f6b7c7039c333360d2','65648bd4926f44c3b35b426207fa6282','888','admin','admin',now(),now(),NULL),
(13,'13892d41d95a4061a7aa346692e8e888','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(14,'45a7c7baa6c6430a850f33f68f3f4eba','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(15,'6c7bcf4822e5437d8448b322f052ceee','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(16,'a218a0a6afb3447a833533203e1cbdb4','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(17,'b489e10af58c47898e8768d67b54a434','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(18,'bd0c847f6d1143ec88bcbef375e749fe','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(19,'c5c096870b784305b79ac9aa339d5591','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(20,'80e52c60f538495ea621a27ef3b245ed','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(21,'127e86a82143494d9be270b3284cb783','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(22,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(23,'6a938cd43f3d4d2fa31f3ca7834d57a2','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(24,'9c7d1a14713949f6b7c7039c333360d2','87289b2334254c25b9ea6f97c213f82c','888','admin','admin',now(),now(),NULL),
(25,'13892d41d95a4061a7aa346692e8e888','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(26,'45a7c7baa6c6430a850f33f68f3f4eba','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(27,'6c7bcf4822e5437d8448b322f052ceee','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(28,'a218a0a6afb3447a833533203e1cbdb4','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(29,'b489e10af58c47898e8768d67b54a434','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(30,'bd0c847f6d1143ec88bcbef375e749fe','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(31,'c5c096870b784305b79ac9aa339d5591','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(32,'80e52c60f538495ea621a27ef3b245ed','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(33,'127e86a82143494d9be270b3284cb783','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(34,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(35,'6a938cd43f3d4d2fa31f3ca7834d57a2','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(36,'9c7d1a14713949f6b7c7039c333360d2','5f3969b6b92645e3bf83e788dd9bd691','888','admin','admin',now(),now(),NULL),
(37,'13892d41d95a4061a7aa346692e8e888','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(38,'45a7c7baa6c6430a850f33f68f3f4eba','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(39,'6c7bcf4822e5437d8448b322f052ceee','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(40,'a218a0a6afb3447a833533203e1cbdb4','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(41,'b489e10af58c47898e8768d67b54a434','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(42,'bd0c847f6d1143ec88bcbef375e749fe','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(43,'c5c096870b784305b79ac9aa339d5591','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(44,'80e52c60f538495ea621a27ef3b245ed','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(45,'127e86a82143494d9be270b3284cb783','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(46,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(47,'6a938cd43f3d4d2fa31f3ca7834d57a2','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(48,'9c7d1a14713949f6b7c7039c333360d2','b6051d4ae94b41e5ae66288edcb9803b','888','admin','admin',now(),now(),NULL),
(49,'13892d41d95a4061a7aa346692e8e888','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(50,'45a7c7baa6c6430a850f33f68f3f4eba','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(51,'6c7bcf4822e5437d8448b322f052ceee','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(52,'a218a0a6afb3447a833533203e1cbdb4','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(53,'b489e10af58c47898e8768d67b54a434','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(54,'bd0c847f6d1143ec88bcbef375e749fe','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(55,'c5c096870b784305b79ac9aa339d5591','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(56,'80e52c60f538495ea621a27ef3b245ed','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(57,'127e86a82143494d9be270b3284cb783','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(58,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(59,'6a938cd43f3d4d2fa31f3ca7834d57a2','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(60,'9c7d1a14713949f6b7c7039c333360d2','77964202be5d4d8296eaaa7ef0bb595a','888','admin','admin',now(),now(),NULL),
(61,'13892d41d95a4061a7aa346692e8e888','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(62,'45a7c7baa6c6430a850f33f68f3f4eba','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(63,'6c7bcf4822e5437d8448b322f052ceee','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(64,'a218a0a6afb3447a833533203e1cbdb4','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(65,'b489e10af58c47898e8768d67b54a434','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(66,'bd0c847f6d1143ec88bcbef375e749fe','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(67,'c5c096870b784305b79ac9aa339d5591','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(68,'80e52c60f538495ea621a27ef3b245ed','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(69,'127e86a82143494d9be270b3284cb783','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(70,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(71,'6a938cd43f3d4d2fa31f3ca7834d57a2','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(72,'9c7d1a14713949f6b7c7039c333360d2','8319e5c4385e46cbbda8d4355c9a0e55','888','admin','admin',now(),now(),NULL),
(73,'13892d41d95a4061a7aa346692e8e888','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(74,'45a7c7baa6c6430a850f33f68f3f4eba','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(75,'6c7bcf4822e5437d8448b322f052ceee','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(76,'a218a0a6afb3447a833533203e1cbdb4','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(77,'b489e10af58c47898e8768d67b54a434','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(78,'bd0c847f6d1143ec88bcbef375e749fe','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(79,'c5c096870b784305b79ac9aa339d5591','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(80,'80e52c60f538495ea621a27ef3b245ed','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(81,'127e86a82143494d9be270b3284cb783','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(82,'2ba7adb2b3cd44ea9ecf9fa10ae30b69','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(83,'6a938cd43f3d4d2fa31f3ca7834d57a2','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL),
(84,'9c7d1a14713949f6b7c7039c333360d2','017a992c460b470da68f262deadc9f78','888','admin','admin',now(),now(),NULL);
UNLOCK TABLES;


DROP TABLE IF EXISTS `dict_info`;
CREATE TABLE `dict_info` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `dict_code` varchar(120) NOT NULL COMMENT '字典code',
  `dict_item_name` varchar(120) NOT NULL COMMENT '字典项name',
  `dict_item_value` varchar(120) NOT NULL COMMENT '字典项value',
  `comment` varchar(2000) NOT NULL DEFAULT '描述',
  `seq` int(11) DEFAULT NULL COMMENT '排列序号',
  `status` varchar(3) NOT NULL DEFAULT '888',
  `creator` varchar(64) NOT NULL  DEFAULT 'admin',
  `modifier` varchar(64) NOT NULL  DEFAULT 'admin',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `dict_info` WRITE;
INSERT INTO `dict_info` VALUES (1,'性别','男','0','性别字典项desc',1,'888','admin','admin',now(),now(),NULL),
(2,'性别','女','1','性别字典项desc',2,'888','admin','admin',now(),now(),NULL),
(3,'字典项测试','测试字典code_1','00','字典项测试desc',1,'888','admin','admin',now(),now(),NULL),
(4,'字典项测试','测试字典code_2','01','字典项测试desc',2,'888','admin','admin',now(),now(),NULL),
(5,'字典项测试','测试字典code_3','02','字典项测试desc',3,'888','admin','admin',now(),now(),NULL),
(6,'字典项测试','测试字典code_4','03','字典项测试desc',4,'888','admin','admin',now(),now(),NULL);
UNLOCK TABLES;