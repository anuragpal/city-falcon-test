CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "first_name" character varying(20) NOT NULL,
  "last_name" character varying(20) NOT NULL,
  "middle_name" character varying(20) NULL,
  "email" character varying(50) NOT NULL,
  "password" text NOT NULL,
  "added_date" timestamptz NOT NULL,
  "modified_date" timestamptz NOT NULL
);