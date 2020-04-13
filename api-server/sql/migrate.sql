Drop table if exists user;
Drop table if exists chara;
Drop table if exists rank;
Drop table if exists user_character;
Create table user (id int NOT NULL  AUTO_INCREMENT    PRIMARY KEY , name varchar(20), token varchar(32), gacha_times int);
Create table chara (id int NOT NULL  AUTO_INCREMENT    PRIMARY KEY, name varchar(20), rank varchar(2));
Create table rank (id int NOT NULL  AUTO_INCREMENT    PRIMARY KEY, name varchar(2), percent double);
Create table user_character (id int NOT NULL AUTO_INCREMENT    PRIMARY KEY, user_id int, character_id int);