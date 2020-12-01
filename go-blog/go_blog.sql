--
-- PostgreSQL database dump
--

-- Dumped from database version 13.0 (Ubuntu 13.0-1.pgdg20.04+1)
-- Dumped by pg_dump version 13.0 (Ubuntu 13.0-1.pgdg20.04+1)

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
-- Name: authors; Type: TABLE; Schema: public; Owner: promise
--

CREATE TABLE public.authors (
    name text,
    email text,
    id text NOT NULL,
    username text,
    password text NOT NULL
);


ALTER TABLE public.authors OWNER TO promise;

--
-- Name: authors_id_seq; Type: SEQUENCE; Schema: public; Owner: promise
--

CREATE SEQUENCE public.authors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authors_id_seq OWNER TO promise;

--
-- Name: authors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: promise
--

ALTER SEQUENCE public.authors_id_seq OWNED BY public.authors.id;


--
-- Name: posts; Type: TABLE; Schema: public; Owner: promise
--

CREATE TABLE public.posts (
    uuid text,
    title text,
    content text,
    author_id character varying(64)
);


ALTER TABLE public.posts OWNER TO promise;

--
-- Name: authors id; Type: DEFAULT; Schema: public; Owner: promise
--

ALTER TABLE ONLY public.authors ALTER COLUMN id SET DEFAULT nextval('public.authors_id_seq'::regclass);


--
-- Name: authors authors_pkey; Type: CONSTRAINT; Schema: public; Owner: promise
--

ALTER TABLE ONLY public.authors
    ADD CONSTRAINT authors_pkey PRIMARY KEY (id);


--
-- Name: posts fk_posts_author; Type: FK CONSTRAINT; Schema: public; Owner: promise
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT fk_posts_author FOREIGN KEY (author_id) REFERENCES public.authors(id);


--
-- PostgreSQL database dump complete
--

