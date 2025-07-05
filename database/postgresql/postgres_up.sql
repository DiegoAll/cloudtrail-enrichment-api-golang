--CREATE DATABASE IF NOT EXISTS authdb;
\c authdb;

DROP TABLE IF EXISTS authdb;

CREATE TABLE
  public.products (
    id serial NOT NULL,
    --uuid character varying(255) UNIQUE NOT NULL, -- Agrega esta línea, idealmente INDEXADO
    name character varying (255) NOT NULL,
    description character varying (255) NOT NULL,
    price numeric (10,2) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
  );

CREATE TABLE public.users (
    id serial NOT NULL,
    email character varying(128) NOT NULL UNIQUE,
    first_name character varying(64) NOT NULL,
    last_name character varying(64) NOT NULL,
    password character varying(128) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);


--FK PENDING: tokens_relation_1 trevor 60 CASCADE CASCADE
CREATE TABLE public.tokens (
    id serial NOT NULL,
    user_id integer,
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    expiry timestamp with time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

INSERT INTO public.users (email, first_name, last_name, password, created_at, updated_at) 
VALUES ('diego@diego.com', 'diego', 'last_name_placeholder', '$2a$12$ZTTPWrzmiT5wSz9gl2FrJuiP4wwoXNKriRbSnZuWCvN/ZgxutLYjG', now(), now());
