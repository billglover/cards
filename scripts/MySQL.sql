DROP TABLE IF EXISTS `cards`;

CREATE TABLE `cards` (
   `uid` SERIAL,
   `title` VARCHAR(64) NULL DEFAULT NULL,
   PRIMARY KEY (`uid`)
);

DROP TABLE IF EXISTS `links`;

CREATE TABLE `links` (
   `uid` SERIAL,
   `parent` BIGINT UNSIGNED NOT NULL,
   `child` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY (`uid`)
);
