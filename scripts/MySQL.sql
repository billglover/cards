DROP TABLE IF EXISTS `cards`;

CREATE TABLE `cards` (
   `uid` SERIAL,
   `title` VARCHAR(64) NULL DEFAULT NULL,
   PRIMARY KEY (`uid`)
);

DROP TABLE IF EXISTS `decks`;

CREATE TABLE `decks` (
   `uid` SERIAL,
   `title` VARCHAR(64) NULL DEFAULT NULL,
   PRIMARY KEY (`uid`)
);