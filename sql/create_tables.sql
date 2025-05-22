--
-- PostgreSQL database --

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.apps (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    release character varying(255),
    path character varying(255),
    init character varying(255),
    web character varying(255),
    title character varying(255),
    created integer,
    updated integer
);

ALTER TABLE public.apps ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.apps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: companies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.companies (
    id integer NOT NULL,
    name character varying(255),
    created integer,
    updated integer
);

ALTER TABLE public.companies ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.companies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: config; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.config (
    id integer NOT NULL,
    name character varying(255),
    value character varying(255),
    created integer,
    updated integer
);

ALTER TABLE public.config ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: dbs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.dbs (
    id integer NOT NULL,
    dsn character varying(255),
    options character varying(255),
    alt_dsn character varying(255),
    created integer,
    updated integer
);

ALTER TABLE public.dbs ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.dbs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: dbsauth; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.dbsauth (
    id integer NOT NULL,
    name character varying(255),
    db_id integer,
    grants character varying(255),
    user_id integer,
    group_id integer,
    grants_options integer,
    created integer,
    updated integer
);

ALTER TABLE public.dbsauth ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.dbsauth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: groups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.groups (
    id integer NOT NULL,
    name character varying(255),
    profile_id integer,
    dbsauth_id integer,
    is_public boolean NOT NULL DEFAULT false,
    leader integer,
    expiration_date timestamp with time zone,
    lan character varying(255) NOT NULL DEFAULT 'es',
    company_id integer,
    created integer,
    updated integer
);

ALTER TABLE public.groups ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: groups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.invitations (
    id integer NOT NULL,
    used boolean NOT NULL DEFAULT false,
    code character varying(255),
    created integer,
    updated integer
);

ALTER TABLE public.invitations ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.invitations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: profiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.profiles (
    id integer NOT NULL,
    name character varying(255),
    app_ids character varying(255),
    db_ids character varying(255),
    created integer,
    updated integer
);

ALTER TABLE public.profiles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: password_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.password_requests (
    id integer NOT NULL,
    user_id integer NOT NULL,
    code character varying(255),
    used boolean NOT NULL DEFAULT false,
    date timestamp with time zone,
    created integer,
    updated integer
);

ALTER TABLE public.password_requests ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.password_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: password_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.security_image (
    random_id integer NOT NULL,
    real_text character varying(255),
    date timestamp with time zone,
    created integer,
    updated integer
);

CREATE UNIQUE INDEX security_image_idx ON public.security_image (real_text, date);

--
-- Name: public.sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id character varying(255),
    value character varying(255),
    time time,
    db character varying(255),
    created integer,
    updated integer
);

CREATE UNIQUE INDEX sessions_idx ON public.sessions (value) INCLUDE (time);

--
-- Name: pages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    code character varying(255),
    active boolean NOT NULL DEFAULT true,
    last_login integer,
    last_session character varying(255) NOT NULL DEFAULT '',
    blocked boolean DEFAULT false,
    tries integer NOT NULL DEFAULT 0,
    last_try integer,
    email character varying(255),
    profile_id integer,
    group_id integer,
    dbsauth_id integer,
    activation_time integer,
    last_action timestamp with time zone,
    last_app integer,
    last_db integer,
    lan character varying(255) NOT NULL DEFAULT 'es',
    company_id integer,
    created integer,
    updated integer
);

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);