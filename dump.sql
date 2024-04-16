--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Drop databases (except postgres and template1)
--

DROP DATABASE planner;




--
-- Drop roles
--

DROP ROLE postgres;


--
-- Roles
--

CREATE ROLE postgres;
ALTER ROLE postgres WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:OxXFeqwRM9bUEazWaz5iAg==$Z9DQqwlassZlKA+xeS7W0KqnEvYckq9IAC4k1nDu69I=:EcWqtMrTs9TOAbJSyBEsexzCpiKT35G7iLr9QjY2xqs=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 15.6 (Debian 15.6-1.pgdg120+2)
-- Dumped by pg_dump version 15.6 (Debian 15.6-1.pgdg120+2)

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

UPDATE pg_catalog.pg_database SET datistemplate = false WHERE datname = 'template1';
DROP DATABASE template1;
--
-- Name: template1; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE template1 WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE template1 OWNER TO postgres;

\connect template1

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

--
-- Name: DATABASE template1; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON DATABASE template1 IS 'default template for new databases';


--
-- Name: template1; Type: DATABASE PROPERTIES; Schema: -; Owner: postgres
--

ALTER DATABASE template1 IS_TEMPLATE = true;


\connect template1

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

--
-- Name: DATABASE template1; Type: ACL; Schema: -; Owner: postgres
--

REVOKE CONNECT,TEMPORARY ON DATABASE template1 FROM PUBLIC;
GRANT CONNECT ON DATABASE template1 TO PUBLIC;


--
-- PostgreSQL database dump complete
--

--
-- Database "planner" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 15.6 (Debian 15.6-1.pgdg120+2)
-- Dumped by pg_dump version 15.6 (Debian 15.6-1.pgdg120+2)

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

--
-- Name: planner; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE planner WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE planner OWNER TO postgres;

\connect planner

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

--
-- Name: authenticate_user(character varying, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.authenticate_user(user_login character varying, user_password character varying) RETURNS TABLE(user_id integer, user_email character varying, is_authenticated boolean)
    LANGUAGE plpgsql
    AS $$
BEGIN
RETURN QUERY
SELECT u.id, u.email,
CASE
WHEN u.password = user_password THEN TRUE
ELSE FALSE
END AS is_authenticated
FROM users u
WHERE u.login = user_login;
END;
$$;


ALTER FUNCTION public.authenticate_user(user_login character varying, user_password character varying) OWNER TO postgres;

--
-- Name: change_user_password(integer, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.change_user_password(id integer, newpassword character varying) RETURNS void
    LANGUAGE plpgsql
    AS $_$
BEGIN
UPDATE users
SET password = newPassword, last_update_time = NOW()
WHERE users.id = $1;
END;
$_$;


ALTER FUNCTION public.change_user_password(id integer, newpassword character varying) OWNER TO postgres;

--
-- Name: create_new_finance(integer, character varying, integer, date); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_new_finance(price_finance integer, currancy_finance character varying, folder integer, date_finance date DEFAULT CURRENT_DATE) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO finance(price, currancy, folder_id, date)
VALUES (price_finance, currancy_finance, folder, date_finance);
END;
$$;


ALTER FUNCTION public.create_new_finance(price_finance integer, currancy_finance character varying, folder integer, date_finance date) OWNER TO postgres;

--
-- Name: create_new_folder(character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_new_folder(name_folder character varying, type_folder character varying, image_folder character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO folder (name, type, image) VALUES (name_folder, type_folder, image_folder);
END;
$$;


ALTER FUNCTION public.create_new_folder(name_folder character varying, type_folder character varying, image_folder character varying) OWNER TO postgres;

--
-- Name: create_new_note(character varying, character varying, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_new_note(title_note character varying, content_note character varying, folder integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO note(title, content, folder_id)
VALUES (title_note, content_note, folder );
END;
$$;


ALTER FUNCTION public.create_new_note(title_note character varying, content_note character varying, folder integer) OWNER TO postgres;

--
-- Name: create_new_task(character varying, boolean, integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_new_task(text_task character varying, is_completed_task boolean, task integer, folder integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO task(text, is_completed, task_id, folder_id)
VALUES (text_task, is_completed_task, task, folder);
END;
$$;


ALTER FUNCTION public.create_new_task(text_task character varying, is_completed_task boolean, task integer, folder integer) OWNER TO postgres;

--
-- Name: create_new_user(character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_new_user(user_login character varying, user_email character varying, user_password character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO users (login, email, password, registration_time, last_update_time)
VALUES (user_login, user_email, user_password, CURRENT_TIMESTAMP, NULL);
END;
$$;


ALTER FUNCTION public.create_new_user(user_login character varying, user_email character varying, user_password character varying) OWNER TO postgres;

--
-- Name: delete_line_finance(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.delete_line_finance(id_to_del integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
DELETE FROM finance WHERE id = id_to_del;
END;
$$;


ALTER FUNCTION public.delete_line_finance(id_to_del integer) OWNER TO postgres;

--
-- Name: delete_line_folder(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.delete_line_folder(id_to_del integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
DELETE FROM folder WHERE id = id_to_del;
END;
$$;


ALTER FUNCTION public.delete_line_folder(id_to_del integer) OWNER TO postgres;

--
-- Name: delete_line_note(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.delete_line_note(id_to_del integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
DELETE FROM note WHERE id = id_to_del;
END;
$$;


ALTER FUNCTION public.delete_line_note(id_to_del integer) OWNER TO postgres;

--
-- Name: delete_line_task(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.delete_line_task(id_to_del integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
DELETE FROM task WHERE id = id_to_del;
END;
$$;


ALTER FUNCTION public.delete_line_task(id_to_del integer) OWNER TO postgres;

--
-- Name: fetch_finance(integer, integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.fetch_finance(start_idx integer, end_idx integer, mth integer) RETURNS TABLE(id integer, price integer, date date, currancy character varying, folder_id integer)
    LANGUAGE plpgsql
    AS $$
BEGIN
RETURN QUERY SELECT f.id, f.price, f.date, f.currancy, f.folder_id FROM finance f WHERE EXTRACT(MONTH FROM f.date) = mth
ORDER BY f.id
OFFSET start_idx - 1 LIMIT end_idx - start_idx + 1;
END;
$$;


ALTER FUNCTION public.fetch_finance(start_idx integer, end_idx integer, mth integer) OWNER TO postgres;

--
-- Name: fetch_folders(integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.fetch_folders(start_idx integer, end_idx integer) RETURNS TABLE(id integer, name character varying, type character varying, image character varying)
    LANGUAGE plpgsql
    AS $$
BEGIN
RETURN QUERY
SELECT f.id, f.name, f.type, f.image
FROM folder f
ORDER BY f.id
OFFSET start_idx - 1
LIMIT end_idx - start_idx + 1;
END;
$$;


ALTER FUNCTION public.fetch_folders(start_idx integer, end_idx integer) OWNER TO postgres;

--
-- Name: fetch_notes(integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.fetch_notes(start_idx integer, end_idx integer) RETURNS TABLE(id integer, title character varying, content character varying, folder_id integer)
    LANGUAGE plpgsql
    AS $$
BEGIN
RETURN QUERY SELECT n.id, n.title, n.content, n.folder_id FROM note n
ORDER BY n.id
OFFSET start_idx - 1 LIMIT end_idx - start_idx + 1;
END;
$$;


ALTER FUNCTION public.fetch_notes(start_idx integer, end_idx integer) OWNER TO postgres;

--
-- Name: fetch_task(integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.fetch_task(start_idx integer, end_idx integer) RETURNS TABLE(id integer, text character varying, is_completed boolean, task_id integer, folder_id integer)
    LANGUAGE plpgsql
    AS $$
BEGIN
RETURN QUERY SELECT t.id, t.text, t.is_completed, t.task_id, t.folder_id FROM task t
ORDER BY COALESCE(t.task_id, t.id), t.id
OFFSET start_idx - 1 LIMIT end_idx - start_idx + 1;
END;
$$;


ALTER FUNCTION public.fetch_task(start_idx integer, end_idx integer) OWNER TO postgres;

--
-- Name: update_finance(integer, integer, character varying, date); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_finance(id_finance integer, price_finance integer, currancy_finance character varying, date_finance date DEFAULT CURRENT_DATE) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE finance
SET price = COALESCE(price_finance, price), date = COALESCE(date_finance, date), currancy = COALESCE(currancy_finance, currancy) WHERE id = id_finance;
END;
$$;


ALTER FUNCTION public.update_finance(id_finance integer, price_finance integer, currancy_finance character varying, date_finance date) OWNER TO postgres;

--
-- Name: update_folder(integer, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_folder(id_folder integer, type_folder character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE folder
SET type = type_folder WHERE folder.id = id_folder;
END;
$$;


ALTER FUNCTION public.update_folder(id_folder integer, type_folder character varying) OWNER TO postgres;

--
-- Name: update_folder(integer, character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_folder(id_folder integer, name_folder character varying, type_folder character varying, image_folder character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE folder
SET name = COALESCE(name_folder, name), type = COALESCE(type_folder, type), image = COALESCE(image_folder, image) WHERE id = id_folder;
END;
$$;


ALTER FUNCTION public.update_folder(id_folder integer, name_folder character varying, type_folder character varying, image_folder character varying) OWNER TO postgres;

--
-- Name: update_note(integer, character varying, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_note(id_note integer, title_note character varying, content_note character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE note
SET title = COALESCE(title_note, title), content = COALESCE(content_note, content) WHERE id = id_note;
END;
$$;


ALTER FUNCTION public.update_note(id_note integer, title_note character varying, content_note character varying) OWNER TO postgres;

--
-- Name: update_price_finance(integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_price_finance(id_finance integer, price_finance integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE finance
SET price = price_finance WHERE id = id_finance;
END;
$$;


ALTER FUNCTION public.update_price_finance(id_finance integer, price_finance integer) OWNER TO postgres;

--
-- Name: update_task(integer, character varying, boolean); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_task(id_task integer, text_task character varying, is_completed_task boolean) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE task
SET text = COALESCE(text_task, text), is_completed = COALESCE(is_completed_task, is_completed) WHERE id = id_task;
END;
$$;


ALTER FUNCTION public.update_task(id_task integer, text_task character varying, is_completed_task boolean) OWNER TO postgres;

--
-- Name: update_user_last_update(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_user_last_update(id integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE users
SET last_update_time = NOW()
WHERE users.id = id;
END;
$$;


ALTER FUNCTION public.update_user_last_update(id integer) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: finance; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.finance (
    id integer NOT NULL,
    price integer NOT NULL,
    date date DEFAULT CURRENT_DATE NOT NULL,
    currancy character varying(255) NOT NULL,
    folder_id integer NOT NULL
);


ALTER TABLE public.finance OWNER TO postgres;

--
-- Name: finance_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.finance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.finance_id_seq OWNER TO postgres;

--
-- Name: finance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.finance_id_seq OWNED BY public.finance.id;


--
-- Name: folder; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.folder (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    type character varying(255) NOT NULL,
    image character varying(255) NOT NULL
);


ALTER TABLE public.folder OWNER TO postgres;

--
-- Name: folder_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.folder_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.folder_id_seq OWNER TO postgres;

--
-- Name: folder_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.folder_id_seq OWNED BY public.folder.id;


--
-- Name: note; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.note (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    content character varying(255) NOT NULL,
    folder_id integer NOT NULL
);


ALTER TABLE public.note OWNER TO postgres;

--
-- Name: note_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.note_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.note_id_seq OWNER TO postgres;

--
-- Name: note_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.note_id_seq OWNED BY public.note.id;


--
-- Name: task; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.task (
    id integer NOT NULL,
    text character varying(255),
    is_completed boolean,
    task_id integer,
    folder_id integer NOT NULL
);


ALTER TABLE public.task OWNER TO postgres;

--
-- Name: task_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.task_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.task_id_seq OWNER TO postgres;

--
-- Name: task_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.task_id_seq OWNED BY public.task.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    registration_time timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_update_time timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: finance id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.finance ALTER COLUMN id SET DEFAULT nextval('public.finance_id_seq'::regclass);


--
-- Name: folder id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.folder ALTER COLUMN id SET DEFAULT nextval('public.folder_id_seq'::regclass);


--
-- Name: note id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.note ALTER COLUMN id SET DEFAULT nextval('public.note_id_seq'::regclass);


--
-- Name: task id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task ALTER COLUMN id SET DEFAULT nextval('public.task_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: finance; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.finance (id, price, date, currancy, folder_id) FROM stdin;
5	159	2024-03-27	rub	3
4	158	2024-03-27	rub	3
6	-289	2024-05-29	rub	3
7	165	2024-03-27	rub	3
8	250	2024-03-27	rub	3
12	100	2024-03-27	USD	3
14	100	2024-03-27	USD	3
\.


--
-- Data for Name: folder; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.folder (id, name, type, image) FROM stdin;
3	test1	mem	govno
4	test2	mem	govno
5	test3	mem	govno
6	test4	mem	govno
7	test5	mem	govno
8	test6	mem	govno
11	testoncode	govno	govno
12	testoncode	govno	govno
9	test11 update	not1 meme	pepego11
2	Новое имя папки	Новый тип	новый_путь_к_изображению
\.


--
-- Data for Name: note; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.note (id, title, content, folder_id) FROM stdin;
2	govno1	alsogovno	2
3	govno2	alsogovno	2
5	testoncode	govno	5
6	testoncode	supergovno	5
7	testoncode	ultrasupergovno	5
8	testoncode	govno on 6 folder	6
4	not11 govno	update1 test	2
\.


--
-- Data for Name: task; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.task (id, text, is_completed, task_id, folder_id) FROM stdin;
4	govno4	f	\N	3
1	new govno1	f	\N	3
3	govno3	t	\N	3
11	new govno1.1	f	1	3
12	new govno1.2	f	1	3
14	govno4.1	f	4	3
16	new text1	t	1	3
19	text	t	\N	3
20	new text1	t	\N	3
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, login, email, password, registration_time, last_update_time) FROM stdin;
1	userLogin	userEmail@example.com	userPassword	2024-03-22 10:13:38.764183	\N
2	testNewVAl	Test.Apperscale@example.com	JzKT+Nl4bR6ZSiqKz/cGuA==	2024-03-22 10:41:42.515758	\N
\.


--
-- Name: finance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.finance_id_seq', 14, true);


--
-- Name: folder_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.folder_id_seq', 12, true);


--
-- Name: note_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.note_id_seq', 9, true);


--
-- Name: task_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.task_id_seq', 20, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: finance finance_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.finance
    ADD CONSTRAINT finance_pkey PRIMARY KEY (id);


--
-- Name: folder folder_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.folder
    ADD CONSTRAINT folder_pkey PRIMARY KEY (id);


--
-- Name: note note_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.note
    ADD CONSTRAINT note_pkey PRIMARY KEY (id);


--
-- Name: task task_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task
    ADD CONSTRAINT task_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_login_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_login_key UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: finance finance_folder_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.finance
    ADD CONSTRAINT finance_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.folder(id) ON DELETE CASCADE;


--
-- Name: note note_folder_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.note
    ADD CONSTRAINT note_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.folder(id) ON DELETE CASCADE;


--
-- Name: task task_folder_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task
    ADD CONSTRAINT task_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.folder(id) ON DELETE CASCADE;


--
-- Name: task task_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task
    ADD CONSTRAINT task_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.task(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 15.6 (Debian 15.6-1.pgdg120+2)
-- Dumped by pg_dump version 15.6 (Debian 15.6-1.pgdg120+2)

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

DROP DATABASE postgres;
--
-- Name: postgres; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO postgres;

\connect postgres

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

--
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

