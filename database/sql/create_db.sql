BEGIN;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
DROP TABLE IF EXISTS stages CASCADE;
DROP TABLE IF EXISTS tests CASCADE;
DROP TABLE IF EXISTS option_tests CASCADE;
DROP TABLE IF EXISTS rewrite_tests CASCADE;
DROP TABLE IF EXISTS progress CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS role_type;
DROP TYPE IF EXISTS difficulty_type;
DROP TYPE IF EXISTS test_type;
DROP TYPE IF EXISTS lang_type;

CREATE TYPE role_type AS ENUM ('user','admin');
CREATE TYPE difficulty_type AS ENUM ('elementary', 'middle', 'upper-middle', 'advance', 'professional', 'god');
CREATE TYPE test_type AS ENUM ('option', 'rewrite');
CREATE TYPE lang_type AS ENUM ('en','fr','it','ru');

CREATE TABLE IF NOT EXISTS users
(
    id                SERIAL PRIMARY KEY,
    name              VARCHAR                     NOT NULL,
    avatar_img        VARCHAR                     NOT NULL,
    login             VARCHAR                     NOT NULL,
    password          VARCHAR                     NOT NULL,
    salt              INT                         NOT NULL,
    role              ROLE_TYPE                   NOT NULL,
    email             VARCHAR                     NOT NULL,
    date_registration TIMESTAMP without time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS categories
(
    id          SMALLSERIAL PRIMARY KEY,
    parent_id   INT       NOT NULL,
    name        VARCHAR   NOT NULL,
    lang_type   lang_type NOT NULL,
    description VARCHAR   NULL,

    CONSTRAINT parent_id_fk FOREIGN KEY (parent_id) REFERENCES categories (id)
);

CREATE TABLE IF NOT EXISTS courses
(
    id          SERIAL PRIMARY KEY,
    category_id SMALLINT      NOT NULL,
    name        VARCHAR       NOT NULL,
    tags        VARCHAR ARRAY NULL,
    rating      SMALLINT      NOT NULL,
    description VARCHAR       NULL,

    CONSTRAINT category_id_fk FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE IF NOT EXISTS stages
(
    id              SERIAL PRIMARY KEY,
    course_id       INT             NOT NULL,
    name            VARCHAR         NOT NULL,
    content         TEXT            NOT NULL,
    order_number    SMALLINT        NOT NULL,
    difficulty_type difficulty_type NOT NULL,
    lemmings_count  SMALLINT        NOT NULL,

    CONSTRAINT course_id_fk FOREIGN KEY (course_id) REFERENCES courses (id)

);

CREATE TABLE IF NOT EXISTS tests
(
    id              SERIAL PRIMARY KEY,
    stage_id        INT             NOT NULL,
    name            VARCHAR         NOT NULL,
    test_type       test_type       NOT NULL,
    order_number    SMALLINT        NOT NULL,
    difficulty_type difficulty_type NOT NULL,

    CONSTRAINT stage_id_fk FOREIGN KEY (stage_id) REFERENCES stages (id)

);

CREATE TABLE IF NOT EXISTS option_tests
(
    test_id             INT PRIMARY KEY NOT NULL,
    question            TEXT            NOT NULL,
    options             VARCHAR ARRAY   NOT NULL,
    right_option_number SMALLINT        NOT NULL,

    CONSTRAINT test_id_fk FOREIGN KEY (test_id) REFERENCES tests (id)
);

CREATE TABLE IF NOT EXISTS rewrite_tests
(
    test_id      INT PRIMARY KEY NOT NULL,
    question     TEXT            NOT NULL,
    rewrite_text VARCHAR         NOT NULL,

    CONSTRAINT test_id_fk FOREIGN KEY (test_id) REFERENCES tests (id)
);


CREATE TABLE IF NOT EXISTS progress
(
    id      SERIAL PRIMARY KEY,
    test_id INT     NOT NULL,
    user_id INT     NOT NULL,
    passed  BOOLEAN NOT NULL,

    CONSTRAINT test_id_fk FOREIGN KEY (test_id) REFERENCES tests (id),
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)

);

COMMIT;
