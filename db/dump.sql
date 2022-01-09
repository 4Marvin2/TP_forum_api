CREATE EXTENSION IF NOT EXISTS CITEXT;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE TABLE IF NOT EXISTS forums(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    username CITEXT NOT NULL,
    slug CITEXT NOT NULL UNIQUE,
    posts BIGINT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    parent BIGINT DEFAULT 0,
    path BIGINT[] NOT NULL DEFAULT '{0}',
    author CITEXT NOT NULL,
    message TEXT NOT NULL,
    isEdited BOOL DEFAULT false,
    forum CITEXT,
    thread INT,
    created TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS threads(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    author CITEXT NOT NULL,
    forum CITEXT,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT,
    created TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    nickname CITEXT NOT NULL UNIQUE,
    fullname CITEXT NOT NULL,
    about TEXT,
    email CITEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forum_users(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    forum_id BIGINT REFERENCES forums(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS votes(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    thread_id BIGINT REFERENCES threads(id) NOT NULL,
    voice SMALLINT NOT NULL
);

-- CREATE FUNCTION update_thread_votes_after_insert()
--     RETURNS TRIGGER AS '
--     BEGIN
--         UPDATE threads
--         SET
--             votes = votes + NEW.voice
--         WHERE id = NEW.thread;
--         RETURN NULL;
--     END;
-- ' LANGUAGE plpgsql;

-- CREATE TRIGGER on_vote_insert
--     AFTER INSERT ON votes
--     FOR EACH ROW EXECUTE PROCEDURE update_thread_votes_after_insert();

-- CREATE FUNCTION update_thread_votes_after_update()
--     RETURNS TRIGGER AS '
--     BEGIN
--         IF OLD.voice = NEW.voice
--         THEN
--             RETURN NULL;
--         END IF;
--         UPDATE threads
--         SET
--             votes = votes + CASE
--                 WHEN NEW.voice = -1
--                 THEN -2
--                 ELSE 2
--                 END
--         WHERE id = NEW.thread;
--         RETURN NULL;
--     END;
-- ' LANGUAGE plpgsql;

-- CREATE TRIGGER on_vote_update
--     AFTER UPDATE ON votes
--     FOR EACH ROW EXECUTE PROCEDURE update_thread_votes_after_update();
