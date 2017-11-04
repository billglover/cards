DROP TABLE IF EXISTS `cards`;

CREATE TABLE `cards` (
   `uid` VARCHAR(24) NOT NULL,
   `title` VARCHAR(64) NULL DEFAULT NULL,
   PRIMARY KEY (`uid`)
);

DROP TABLE IF EXISTS `links`;

CREATE TABLE `links` (
   `parent` VARCHAR(24) NOT NULL,
   `child` VARCHAR(24) NOT NULL,
   PRIMARY KEY (`parent`, `child`)
);
