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

CREATE INDEX fk_nickname ON forums (nickname);



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

--CREATE INDEX ON threads (forum_id, author);

CREATE INDEX ON threads (created);
CREATE INDEX fk_threads_author ON threads(author);
CREATE INDEX fk_forum_id ON threads(forum_id);


CREATE TABLE IF NOT EXISTS posts (
    id bigint PRIMARY KEY,
    path bigint[] NOT NULL,
  --  path1 bigint NOT NULL,

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
CREATE SEQUENCE pidseq START 1;

CREATE INDEX ON posts (forum_id, author);
CREATE INDEX fk_thread_id ON posts (thread_id);
CREATE INDEX fk_thread_id_id ON posts (thread_id, id);   -- was only thread_id, maybe it's extra (no better)
CREATE INDEX id_thread_id_path_idx ON posts /*USING GIN*/(thread_id, path);  -- was (thread_id, path)
CREATE INDEX id_thread_id_path_root_idx ON posts /*USING GIN*/(thread_id, path) WHERE array_length(path, 1) = 1;
--CREATE INDEX posts_author ON posts (author);
--CREATE INDEX posts_author_forum_id ON posts (author, forum_id);
--CREATE INDEX ON posts (path1);


CREATE TABLE IF NOT EXISTS votes (
    nickname citext NOT NULL,
    thread_id integer NOT NULL,
    voice integer NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users(nickname),
    FOREIGN KEY (thread_id) REFERENCES threads(id),
    UNIQUE (thread_id, nickname)
);


CREATE TABLE IF NOT EXISTS forum_user(
    nickname citext NOT NULL,
    forum_id integer NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users(nickname),
    FOREIGN KEY (forum_id) REFERENCES forums(id),
    UNIQUE (forum_id, nickname)
);

--CREATE UNIQUE INDEX forum_user_id_nickname ON forum_user (forum_id, nickname);


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


CREATE FUNCTION on_threads_ins_upd() RETURNS trigger AS $on_threads_ins_upd$
    BEGIN
        IF (TG_OP = 'INSERT') THEN
            UPDATE forums
            SET threads = threads + (SELECT COUNT(id) FROM inserted)
            WHERE id = (SELECT forum_id FROM inserted LIMIT 1);

            WITH f AS (SELECT forum_id FROM inserted LIMIT 1)
            INSERT INTO forum_user (nickname, forum_id)
            SELECT author, (SELECT forum_id FROM f) FROM inserted
            ON CONFLICT DO NOTHING;

--         ELSIF (TG_OP = 'UPDATE') THEN

        END IF;
        RETURN NULL;
    END;
$on_threads_ins_upd$ LANGUAGE plpgsql;

CREATE TRIGGER threads_ins_upd
    AFTER INSERT ON threads
    REFERENCING NEW TABLE AS inserted
    FOR EACH STATEMENT
    EXECUTE PROCEDURE on_threads_ins_upd();


CREATE FUNCTION on_posts_ins_upd() RETURNS trigger AS $on_posts_ins_upd$
    BEGIN
        IF (TG_OP = 'INSERT') THEN
            UPDATE forums
            SET posts = posts + (SELECT COUNT(id) FROM inserted)
            WHERE id = (SELECT forum_id FROM inserted LIMIT 1);

            WITH f AS (SELECT forum_id FROM inserted LIMIT 1)
            INSERT INTO forum_user (nickname, forum_id)
            SELECT author, (SELECT forum_id FROM f) FROM inserted
            ON CONFLICT DO NOTHING;

--         ELSIF (TG_OP = 'UPDATE') THEN

        END IF;
        RETURN NULL;
    END;
$on_posts_ins_upd$ LANGUAGE plpgsql;

CREATE TRIGGER posts_ins_upd
    AFTER INSERT ON posts
    REFERENCING NEW TABLE AS inserted
    FOR EACH STATEMENT
    EXECUTE PROCEDURE on_posts_ins_upd();

