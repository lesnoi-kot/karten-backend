CREATE SCHEMA karten;
SET search_path TO karten;

CREATE TABLE projects (
  id      uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  name    varchar(32) NOT NULL CHECK (length("name") > 0)
);

CREATE TABLE boards (
  id                  uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  project_id          uuid REFERENCES projects ON DELETE CASCADE,
  name                varchar(32) NOT NULL CHECK (length("name") > 0),
  archived            boolean DEFAULT false NOT NULL,
  date_created        timestamp DEFAULT CURRENT_TIMESTAMP,
  date_last_viewed    timestamp DEFAULT CURRENT_TIMESTAMP,
  color               integer DEFAULT 0 NOT NULL,
  cover_url           varchar(64)
);

CREATE TABLE task_lists (
  id              uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  board_id        uuid REFERENCES boards ON DELETE CASCADE,
  name            varchar(32) NOT NULL CHECK (length("name") > 0),
  position        integer NOT NULL,
  archived        boolean DEFAULT false NOT NULL,
  date_created    timestamp DEFAULT CURRENT_TIMESTAMP,
  color           integer DEFAULT 0 NOT NULL
);

CREATE TABLE tasks (
  id              uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  task_list_id    uuid REFERENCES task_lists ON DELETE CASCADE,
  name            text NOT NULL CHECK (length("name") > 0),
  text            text DEFAULT '' NOT NULL,
  position        integer NOT NULL,
  archived        boolean DEFAULT false NOT NULL,
  date_created    timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  due_date        timestamp
);

CREATE TABLE comments (
  id              uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  task_id         uuid REFERENCES tasks ON DELETE CASCADE,
  text            text NOT NULL CHECK (length("text") > 0),
  author          varchar(32) DEFAULT 'User',
  date_created    timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
);
