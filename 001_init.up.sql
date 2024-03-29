-- Connect to Postgres.
--psql -U cms

-- Create a user named `cms` with the password `cms`. If you ever
-- do this in production, please use a better password.
--CREATE USER cms WITH PASSWORD 'cms';

-- Create the database we're going to use.
--CREATE DATABASE cms;

-- Grant all privleges to our user on the DB.
--GRANT ALL PRIVILEGES ON DATABASE cms to cms;

-- Create a new table to store our pages.
BEGIN;
CREATE TABLE IF NOT EXISTS PAGES(
  id             SERIAL    PRIMARY KEY,
  title          TEXT      NOT NULL,
  content        TEXT      NOT NULL
);

-- Create a new table to store our posts.
CREATE TABLE IF NOT EXISTS POSTS(
  id             SERIAL    PRIMARY KEY,
  title          TEXT      NOT NULL,
  content        TEXT      NOT NULL,
  date_created   DATE      NOT NULL
);

-- Create a new table to store our comments.
CREATE TABLE IF NOT EXISTS COMMENTS(
  id             SERIAL    PRIMARY KEY,
  author         TEXT      NOT NULL,
  content        TEXT      NOT NULL,
  date_created   DATE      NOT NULL,
  post_id        INT       references POSTS(id)
);
COMMIT;