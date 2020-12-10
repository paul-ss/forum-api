CREATE USER forum_user WITH password '662f2710-4e08';

create database forum_db
	with owner forum_user
	encoding 'utf8'
	LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
--     TABLESPACE = forum_default
;

GRANT ALL PRIVILEGES ON database forum_db TO forum_user;