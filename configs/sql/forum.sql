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
    FOREIGN KEY forum_id REFERENCES forums(id)
);


CREATE TABLE IF NOT EXISTS posts (
    description text,
    id bigserial PRIMARY KEY NOT NULL,
    path bigint[] NOT NULL,

    parent_id bigint NOT NULL, -- tree
    author text NOT NULL, -- fk

    messagetext NOT NULL,
    isEdited bool NOT NULL DEFAULT false,

    forum string, -- slug
    forum_id integer NOT NULL,

    thread_id integer NOT NULL, --fk
    created timestamp DEFAULT NOW(),

    FOREIGN KEY author REFERENCES users(nickname),
    FOREIGN KEY forum_id REFERENCES forums(id),
    FOREIGN KEY thread_id REFERENCES threads(id)
);