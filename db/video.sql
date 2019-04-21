CREATE TABLE videos (
  id UUID PRIMARY KEY,

  unique_id text NOT NULL UNIQUE,
  duration bigint NOT NULL,
  used boolean NOT NULL,

  status smallint NOT NULL,
  origin smallint NOT NULL,

  created_at bigint NOT NULL
);

CREATE TABLE prods (
  id UUID PRIMARY KEY,

  unique_id text NOT NULL UNIQUE,
  duration bigint NOT NULL,

  created_at bigint NOT NULL
);
