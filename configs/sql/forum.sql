CREATE TABLE IF NOT EXISTS users (
    nickname text NOT NULL PRIMARY KEY,
    fullname text NOT NULL,
    about text,
    email text NOT NULL UNIQUE
);




CREATE TABLE IF NOT EXISTS forums (
    id serial PRIMARY KEY NOT NULL,

    title text NOT NULL,
    nickname text NOT NULL,
    slug text NOT NULL UNIQUE,

    -- trigger ?
    posts bigint DEFAULT 0,
    threads	integer DEFAULT 0,

    FOREIGN KEY (nickname) REFERENCES users(nickname)
);
CREATE INDEX ON forums (slug);

CREATE TABLE IF NOT EXISTS threads (
    id serial PRIMARY KEY NOT NULL,
    title text NOT NULL,
    author text NOT NULL,

    forum_title text NOT NULL,
    forum_id integer NOT NULL,

    message text NOT NULL,
    votes integer DEFAULT 0, -- trigger
    slug text, --optional
    created timestamp DEFAULT NOW(),

    FOREIGN KEY (author) REFERENCES users(nickname),
    FOREIGN KEY (forum_id) REFERENCES forums(id)
);
CREATE INDEX ON threads (forum_id, author);


CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY NOT NULL,
    path bigint[] NOT NULL,

    parent_id bigint NOT NULL, -- tree
    author text NOT NULL, -- fk

    message text NOT NULL,
    isEdited bool NOT NULL DEFAULT false,

    forum text, -- slug
    forum_id integer NOT NULL,

    thread_id integer NOT NULL, --fk
    created timestamp DEFAULT NOW(),

    FOREIGN KEY (author) REFERENCES users(nickname),
    FOREIGN KEY (forum_id) REFERENCES forums(id),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);
CREATE INDEX ON posts (forum_id, author);


CREATE TABLE IF NOT EXISTS votes (
    nickname text NOT NULL,
    voice integer NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users(nickname)
);

-- ======================

-- CREATE TABLE IF NOT EXISTS user_forum (
--     nickname text NOT NULL,
--     forum_id integer NOT NULL,
--     FOREIGN KEY (nickname) REFERENCES users(nickname),
--     FOREIGN KEY (forum_id) REFERENCES forums(id),
--     UNIQUE(nickname, forum_id)
-- );
-- CREATE INDEX ON user_forum (nickname);
--
--
-- CREATE TRIGGER insert_user_forum
--     AFTER INSERT ON threads
--     FOR EACH ROW
--     EXECUTE PROCEDURE on_change_threads();
--
-- CREATE FUNCTION on_change_threads() RETURNS trigger AS $on_change_threads$
--     BEGIN
--         IF (TG_OP = 'INSERT') THEN
--             INSERT INTO user_forum (nickname, forum_id)
--             VALUES (NEW.author, NEW.forum_id);
--         END IF;
--         RETURN NULL;
--     END;
-- $on_change_threads$ LANGUAGE plpgsql;
