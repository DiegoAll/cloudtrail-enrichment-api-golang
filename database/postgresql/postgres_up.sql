--CREATE DATABASE IF NOT EXISTS authdb;
\c authdb;

DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;

CREATE TABLE public.users (
    id serial NOT NULL PRIMARY KEY,
    -- Añadir la columna uuid
    uuid character varying(36) NOT NULL UNIQUE, -- UUIDs son típicamente 36 caracteres (e.g., xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx)
    email character varying(128) NOT NULL UNIQUE,
    -- Cambiar 'password' a 'password_hash' para que coincida con el código Go
    password_hash character varying(128) NOT NULL, 
    -- Añadir la columna 'role'
    role character varying(32) NOT NULL DEFAULT 'user', -- Puedes ajustar el tamaño y el valor por defecto
    first_name character varying(64), -- Hacer opcional si no lo usas en el registro
    last_name character varying(64),  -- Hacer opcional si no lo usas en el registro
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

--FK PENDING: tokens_relation_1 trevor 60 CASCADE CASCADE
-- Crear la tabla 'tokens'
CREATE TABLE public.tokens (
    id serial NOT NULL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES public.users(id) ON DELETE CASCADE, -- Asegúrate de que user_id sea NOT NULL
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL UNIQUE, -- El token de texto plano puede ser único
    token_hash bytea NOT NULL UNIQUE, -- El hash del token debería ser único
    expiry timestamp with time zone NOT NULL,
    -- Añadir la columna 'role' a la tabla de tokens si también se guarda allí
    role character varying(32) NOT NULL DEFAULT 'user',
    created_at timestamp without time zone NOT NULL DEFAULT now(), -- Usar now() para default
    updated_at timestamp without time zone NOT NULL DEFAULT now()  -- Usar now() para default
);


INSERT INTO public.users (uuid, email, first_name, last_name, password_hash, role, created_at, updated_at) 
VALUES (
    'a1b2c3d4-e5f6-7890-1234-567890abcdef', -- UUID de ejemplo
    'diego@diego.com', 
    'diego', 
    'last_name_placeholder', 
    '$2a$12$ZTTPWrzmiT5wSz9gl2FrJuiP4wwoXNKriRbSnZuWCvN/ZgxutLYjG', -- Hash de "123123123"
    'admin', -- Rol de ejemplo
    now(), 
    now()
);
