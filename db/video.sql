CREATE TABLE videos (
  id UUID PRIMARY KEY,

  unique_id text NOT NULL,
  used boolean NOT NULL,

  status smallint NOT NULL,
  origin smallint NOT NULL,

  created_at bigint NOT NULL
);
