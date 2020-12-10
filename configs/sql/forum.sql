CREATE TABLE IF NOT EXISTS users (
    description text,
    nickname text NOT NULL PRIMARY KEY,
    fullname text NOT NULL,
    about text,
    email text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forums (
    description text,
    id serial PRIMARY KEY NOT NULL,
    title text NOT NULL,
    user_nickname text NOT NULL,
    slug text NOT NULL UNIQUE,

    -- trigger ?
    posts bigint DEFAULT 0,
    threads	integer DEFAULT 0,

    FOREIGN KEY user_nickname REFERENCES users(nickname)
);

CREATE TABLE IF NOT EXISTS threads (
    description text,
    id serial PRIMARY KEY NOT NULL,
    title text NOT NULL,
    author text NOT NULL,

    forum_tittle text NOT NULL,
    forum_id integer NOT NULL,

    message text NOT NULL,
    votes integer DEFAULT 0,
    slug text, --optional
    created timestamp DEFAULT NOW(),

    FOREIGN KEY author REFERENCES users(nickname),
    FOREIGN KEY forum_id REFERENCES forums(id),
);


CREATE TABLE IF NOT EXISTS posts (
description:
Сообщение внутри ветки обсуждения на форуме.

id	number($int64)
readOnly: true
Идентификатор данного сообщения.

parent	number($int64)
Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).

author*	string($identity)
example: j.sparrow
x-isnullable: false
Автор, написавший данное сообщение.

message*	string($text)
example: We should be afraid of the Kraken.
x-isnullable: false
Собственно сообщение форума.

isEdited	boolean
readOnly: true
x-isnullable: false
Истина, если данное сообщение было изменено.

forum	string($identity)
readOnly: true
Идентификатор форума (slug) данного сообещния.

thread	number($int32)
readOnly: true
Идентификатор ветви (id) обсуждения данного сообещния.

created	string($date-time)
readOnly: true
x-isnullable: true
Дата создания сообщения на форуме.
);