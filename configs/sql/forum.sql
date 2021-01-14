CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    nickname CITEXT NOT NULL PRIMARY KEY,
    fullname text NOT NULL,
    about text,
    email text NOT NULL
);

CREATE UNIQUE INDEX email_unique_idx on users (LOWER(email));



CREATE TABLE IF NOT EXISTS forums (
    id serial PRIMARY KEY NOT NULL,

    title text NOT NULL,
    nickname citext NOT NULL,
    slug citext NOT NULL UNIQUE,

    -- trigger ?
    posts bigint DEFAULT 0,
    threads	integer DEFAULT 0,

    FOREIGN KEY (nickname) REFERENCES users(nickname)
);

CREATE INDEX ON forums (slug);



CREATE TABLE IF NOT EXISTS threads (
    id serial PRIMARY KEY NOT NULL,
    title text NOT NULL,
    author citext NOT NULL,

    forum_slug citext NOT NULL,
    forum_id integer NOT NULL,

    message text NOT NULL,
    votes integer DEFAULT 0, -- trigger
    slug citext unique, --optional
    created timestamp with time zone DEFAULT NOW(),

    FOREIGN KEY (author) REFERENCES users(nickname),
    FOREIGN KEY (forum_id) REFERENCES forums(id)
);
CREATE INDEX ON threads (forum_id, author);


CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    path bigint[] NOT NULL,

    parent_id bigint, --  extra?
    author citext NOT NULL, -- fk

    message text NOT NULL,
    isEdited bool NOT NULL DEFAULT false,

    forum_slug citext, -- slug
    forum_id integer NOT NULL,

    thread_id integer NOT NULL, --fk
    created timestamp with time zone DEFAULT NOW(),

    FOREIGN KEY (author) REFERENCES users(nickname),
    FOREIGN KEY (forum_id) REFERENCES forums(id),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);

CREATE INDEX ON posts (forum_id, author);


CREATE TABLE IF NOT EXISTS votes (
    nickname citext NOT NULL,
    thread_id integer NOT NULL,
    voice integer NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users(nickname),
    FOREIGN KEY (thread_id) REFERENCES threads(id),
    UNIQUE (thread_id, nickname)
);


-- CREATE TABLE IF NOT EXISTS stats (
--     usr integer NOT NULL DEFAULT 0,
--     forum integer NOT NULL DEFAULT 0,
--     thread integer NOT NULL DEFAULT 0,
--     post bigint NOT NULL DEFAULT 0
-- );

-- ======================

CREATE FUNCTION on_vote_threads() RETURNS trigger AS $on_vote_threads$
    BEGIN
        IF (TG_OP = 'INSERT') THEN
            UPDATE threads
            SET votes = CASE WHEN NEW.voice > 0 THEN votes + 1 ELSE votes - 1 END
            WHERE id = NEW.thread_id;

        ELSIF (TG_OP = 'UPDATE') THEN
            UPDATE threads
            SET votes = votes + NEW.voice - OLD.voice
            WHERE id = NEW.thread_id;

        END IF;
        RETURN NULL;
    END;
$on_vote_threads$ LANGUAGE plpgsql;

CREATE TRIGGER vote_threads
    AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE on_vote_threads();







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
