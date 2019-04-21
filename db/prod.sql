CREATE TABLE prods (
  id SERIAL PRIMARY KEY,

  duration bigint NOT NULL,
  used boolean NOT NULL,

  created_at bigint NOT NULL
);
