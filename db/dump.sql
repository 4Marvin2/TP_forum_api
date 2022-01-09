CREATE EXTENSION IF NOT EXISTS CITEXT;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;
DROP FUNCTION IF EXISTS update_thread_votes_after_insert();
DROP FUNCTION IF EXISTS update_thread_votes_after_update();
DROP FUNCTION IF EXISTS insert_forum_users();

CREATE TABLE IF NOT EXISTS forums(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
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
    title TEXT NOT NULL,
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
    voice INT NOT NULL
);

CREATE FUNCTION update_thread_votes_after_insert()
    RETURNS TRIGGER AS '
    BEGIN
        UPDATE threads
        SET
            votes = votes + NEW.voice
        WHERE id = NEW.thread_id;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_insert
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_thread_votes_after_insert();

CREATE FUNCTION update_thread_votes_after_update()
    RETURNS TRIGGER AS '
    BEGIN
        IF OLD.voice = NEW.voice
        THEN
            RETURN NULL;
        END IF;
        UPDATE threads
        SET
            votes = votes + CASE
                WHEN NEW.voice = -1
                THEN -2
                ELSE 2
                END
        WHERE id = NEW.thread_id;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_update
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_thread_votes_after_update();

CREATE FUNCTION insert_forum_users()
    RETURNS TRIGGER AS '
    BEGIN
        INSERT INTO forum_users (user_id, forum_id) VALUES ((SELECT id FROM users WHERE NEW.author = nickname), (SELECT id FROM forums WHERE NEW.forum = slug));
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_thread_insert
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE insert_forum_users();

CREATE TRIGGER on_posts_insert
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE insert_forum_users();
