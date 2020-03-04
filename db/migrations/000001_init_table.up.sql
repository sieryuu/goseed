CREATE TABLE IF NOT EXISTS t_user (
    id serial PRIMARY KEY,
    username varchar(32) NOT NULL,
    first_name varchar(32) NOT NULL,
    last_name varchar(32) NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamp without time zone,
    created_by varchar(32) NOT NULL,
    last_updated_at timestamp without time zone,
    last_updated_by varchar(32) NOT NULL,
    is_active boolean,
    UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS t_tenant (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    tenancy_name varchar(255) NOT NULL,
    created_at timestamp without time zone,
    created_by varchar(32) NOT NULL,
    last_updated_at timestamp without time zone,
    last_updated_by varchar(32) NOT NULL,
    is_active boolean,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS t_article (
    id serial PRIMARY KEY,
    tenant_id integer REFERENCES t_tenant (id) ON DELETE RESTRICT,
    title varchar(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

