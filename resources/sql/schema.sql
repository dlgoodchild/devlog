CREATE DATABASE `devlog`;

USE `devlog`;

CREATE TABLE `post` (
  `id` int(11) unsigned not null auto_increment,
  `guid` varchar(256) not null default '',
  `title` varchar(256) default null,
  `content` mediumtext,
  `created_at` timestamp not null default current_timestamp,
  `modified_at` timestamp not null on update current_timestamp,
  primary key (`id`),
  unique key `guid` (`guid`)
) engine=InnoDB auto_increment=2 default charset=utf8 collate=utf8_general_ci;

INSERT INTO `post` (`guid`, `title`, `content`, `created_at`) VALUES
  ('hello-world', 'Hello World', 'This is the first <strong>sample</strong> page', CURRENT_TIMESTAMP);