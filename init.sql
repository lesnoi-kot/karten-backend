BEGIN;

CREATE SCHEMA karten;
SET search_path TO karten;

CREATE TABLE files (
  id                   uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  storage_object_id    text UNIQUE NOT NULL CHECK (length("storage_object_id") > 0),
  name                 varchar(255) NOT NULL CHECK (length("name") > 0),
  mime_type            varchar(255) NOT NULL CHECK (length("mime_type") > 0),
  size                 int NOT NULL
);

CREATE TABLE image_thumbnails (
  id          uuid PRIMARY KEY REFERENCES files ON DELETE CASCADE,
  image_id    uuid             REFERENCES files ON DELETE CASCADE
);

CREATE TABLE default_cover_images (
  id uuid PRIMARY KEY REFERENCES files ON DELETE CASCADE
);

CREATE TABLE projects (
  id           uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  short_id     varchar(12) GENERATED ALWAYS AS (substring(id::text, 25)) STORED,
  name         varchar(32) NOT NULL CHECK (length("name") > 0),
  avatar_id    uuid REFERENCES files ON DELETE SET NULL
);

CREATE TABLE boards (
  id                  uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  short_id            varchar(12) GENERATED ALWAYS AS (substring(id::text, 25)) STORED,
  project_id          uuid REFERENCES projects ON DELETE CASCADE,
  name                varchar(32) NOT NULL CHECK (length("name") > 0),
  archived            boolean DEFAULT false NOT NULL,
  favorite            boolean DEFAULT false NOT NULL,
  date_created        timestamp DEFAULT CURRENT_TIMESTAMP,
  date_last_viewed    timestamp DEFAULT CURRENT_TIMESTAMP,
  color               integer DEFAULT 0 NOT NULL,
  cover_id            uuid REFERENCES files ON DELETE SET NULL
);

CREATE TABLE task_lists (
  id              uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  board_id        uuid REFERENCES boards ON DELETE CASCADE,
  name            varchar(32) NOT NULL CHECK (length("name") > 0),
  position        bigint NOT NULL,
  archived        boolean DEFAULT false NOT NULL,
  date_created    timestamp DEFAULT CURRENT_TIMESTAMP,
  color           integer DEFAULT 0 NOT NULL
);

CREATE TABLE tasks (
  id              uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  short_id        varchar(12) GENERATED ALWAYS AS (substring(id::text, 25)) STORED,
  task_list_id    uuid REFERENCES task_lists ON DELETE CASCADE,
  name            text NOT NULL CHECK (length("name") > 0),
  text            text DEFAULT '' NOT NULL,
  position        bigint NOT NULL,
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

COMMIT;
