CREATE TABLE videos (
  id SERIAL PRIMARY KEY,

  unique_id text NOT NULL UNIQUE,
  duration bigint NOT NULL,
  used_in int NULL,

  status smallint NOT NULL,
  origin smallint NOT NULL,

  created_at bigint NOT NULL
);
