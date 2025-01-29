--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.2 (Debian 17.2-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.books (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text NOT NULL,
    description text,
    price text,
    attributes text,
    date text,
    image_url text
);


ALTER TABLE public.books OWNER TO postgres;

--
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.books_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.books_id_seq OWNER TO postgres;

--
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying,
    email character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
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


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: books id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (id, created_at, updated_at, deleted_at, title, description, price, attributes, date, image_url) FROM stdin;
1	2025-01-07 12:19:39.549221+00	2025-01-07 12:19:39.549221+00	\N	The Good Earth	A must-read for those looking to expand their knowledge.	12.50 USD	Softcover, 200 pages	1975	https://covers.openlibrary.org/b/id/2626823-L.jpg
2	2025-01-07 12:19:39.557608+00	2025-01-07 12:19:39.557608+00	\N	Идіотъ	An enriching experience that broadens the mind and spirit.	12.50 USD	Softcover, 200 pages	1989	https://covers.openlibrary.org/b/id/9412225-L.jpg
3	2025-01-07 12:19:39.562602+00	2025-01-07 12:19:39.562602+00	\N	The Voyage of the Dawn Treader	An insightful exploration of a fascinating subject.	12.50 USD	Softcover, 200 pages	2009	https://covers.openlibrary.org/b/id/9184719-L.jpg
4	2025-01-07 12:19:39.567294+00	2025-01-07 12:19:39.567294+00	\N	Arsène Lupin, gentleman-cambrioleur	An enriching experience that broadens the mind and spirit.	12.50 USD	Softcover, 200 pages	2021-03-01	https://covers.openlibrary.org/b/id/6125938-L.jpg
5	2025-01-07 12:19:39.57124+00	2025-01-07 12:19:39.57124+00	\N	The Story of Philosophy	A must-read for those looking to expand their knowledge.	12.50 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/5444146-L.jpg
6	2025-01-07 12:19:39.574659+00	2025-01-07 12:19:39.574659+00	\N	Guess How Much I Love You	A captivating story that engages the reader from start to finish.	8.99 USD	Softcover, 200 pages	October 1995	https://covers.openlibrary.org/b/id/13282906-L.jpg
7	2025-01-07 12:19:39.577039+00	2025-01-07 12:19:39.577039+00	\N	The Thirty-Nine Steps	An enriching experience that broadens the mind and spirit.	10.99 USD	Softcover, 200 pages	Apr 01, 2018	https://covers.openlibrary.org/b/id/93020-L.jpg
8	2025-01-07 12:19:39.578894+00	2025-01-07 12:19:39.578894+00	\N	The Hunger Games	A must-read for those looking to expand their knowledge.	12.50 USD	Softcover, 200 pages	Mar 17, 2008	https://covers.openlibrary.org/b/id/12646537-L.jpg
9	2025-01-07 12:19:39.580705+00	2025-01-07 12:19:39.580705+00	\N	Three Men in a Boat (to say nothing of the dog)	A captivating story that engages the reader from start to finish.	10.99 USD	Softcover, 200 pages	Mar 23, 2015	https://covers.openlibrary.org/b/id/8243006-L.jpg
10	2025-01-07 12:19:39.582996+00	2025-01-07 12:19:39.582996+00	\N	Le Comte de Monte Cristo	An insightful exploration of a fascinating subject.	8.99 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/14566393-L.jpg
11	2025-01-07 12:19:39.584909+00	2025-01-07 12:19:39.584909+00	\N	Uncle Tom's Cabin	A captivating story that engages the reader from start to finish.	8.99 USD	Softcover, 200 pages	September 1977	https://covers.openlibrary.org/b/id/12728198-L.jpg
12	2025-01-07 12:19:39.586624+00	2025-01-07 12:19:39.586624+00	\N	NASA/DoD aerospace knowledge diffusion research project	An insightful exploration of a fascinating subject.	15.75 USD	Softcover, 200 pages	1994	https://covers.openlibrary.org/b/id/8936636-L.jpg
13	2025-01-07 12:19:39.58924+00	2025-01-07 12:19:39.58924+00	\N	Gulliver's Travels	An enriching experience that broadens the mind and spirit.	9.99 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/12717083-L.jpg
14	2025-01-07 12:19:39.590975+00	2025-01-07 12:19:39.590975+00	\N	Hamlet	An insightful exploration of a fascinating subject.	10.99 USD	Softcover, 200 pages	1784	https://covers.openlibrary.org/b/id/8281954-L.jpg
15	2025-01-07 12:19:39.592338+00	2025-01-07 12:19:39.592338+00	\N	Autobiography of a Yogi	A timeless tale that resonates with readers of all ages.	10.99 USD	Softcover, 200 pages	Mar 31, 2001	https://covers.openlibrary.org/b/id/805448-L.jpg
16	2025-01-07 12:19:39.593777+00	2025-01-07 12:19:39.593777+00	\N	Principles of Anatomy and Physiology	An insightful exploration of a fascinating subject.	8.99 USD	Softcover, 200 pages	1996	https://covers.openlibrary.org/b/id/3810109-L.jpg
17	2025-01-07 12:19:39.595323+00	2025-01-07 12:19:39.595323+00	\N	The alchemist, 1612	A timeless tale that resonates with readers of all ages.	15.75 USD	Softcover, 200 pages	1971	https://covers.openlibrary.org/b/id/7463992-L.jpg
18	2025-01-07 12:19:39.596753+00	2025-01-07 12:19:39.596753+00	\N	The Ugly Duckling	An enriching experience that broadens the mind and spirit.	12.50 USD	Softcover, 200 pages	April 1, 2003	https://covers.openlibrary.org/b/id/446546-L.jpg
19	2025-01-07 12:19:39.598336+00	2025-01-07 12:19:39.598336+00	\N	A Midsummer Night's Dream	An enriching experience that broadens the mind and spirit.	15.75 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/7205924-L.jpg
20	2025-01-07 12:19:39.600672+00	2025-01-07 12:19:39.600672+00	\N	Преступление и наказание	An enriching experience that broadens the mind and spirit.	10.99 USD	Softcover, 200 pages	October 3, 1976	https://covers.openlibrary.org/b/id/9411873-L.jpg
21	2025-01-07 12:19:39.603706+00	2025-01-07 12:19:39.603706+00	\N	Much Ado About Nothing	A captivating story that engages the reader from start to finish.	9.99 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/8290853-L.jpg
22	2025-01-07 12:19:39.605626+00	2025-01-07 12:19:39.605626+00	\N	The Lost World	A timeless tale that resonates with readers of all ages.	8.99 USD	Softcover, 200 pages	2015-08-04	https://covers.openlibrary.org/b/id/8231444-L.jpg
23	2025-01-07 12:19:39.607122+00	2025-01-07 12:19:39.607122+00	\N	As You Like It	An enriching experience that broadens the mind and spirit.	9.99 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/7338874-L.jpg
24	2025-01-07 12:19:39.608576+00	2025-01-07 12:19:39.608576+00	\N	The Adventures of Sherlock Holmes [12 stories]	A captivating story that engages the reader from start to finish.	9.99 USD	Softcover, 200 pages	Jun 27, 2017	https://covers.openlibrary.org/b/id/6717853-L.jpg
25	2025-01-07 12:19:39.611488+00	2025-01-07 12:19:39.611488+00	\N	David Copperfield	A captivating story that engages the reader from start to finish.	15.75 USD	Softcover, 200 pages	2019-01-01	https://covers.openlibrary.org/b/id/1048892-L.jpg
26	2025-01-07 12:19:39.613724+00	2025-01-07 12:19:39.613724+00	\N	Othello	An enriching experience that broadens the mind and spirit.	8.99 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/7165018-L.jpg
27	2025-01-07 12:19:39.615173+00	2025-01-07 12:19:39.615173+00	\N	Great Expectations	A captivating story that engages the reader from start to finish.	10.99 USD	Softcover, 200 pages	2006-07-14	https://covers.openlibrary.org/b/id/13322313-L.jpg
28	2025-01-07 12:19:39.61672+00	2025-01-07 12:19:39.61672+00	\N	King Lear	An enriching experience that broadens the mind and spirit.	10.99 USD	Softcover, 200 pages	October 2005	https://covers.openlibrary.org/b/id/7420452-L.jpg
29	2025-01-07 12:19:39.618612+00	2025-01-07 12:19:39.618612+00	\N	The Merchant of Venice	A must-read for those looking to expand their knowledge.	8.99 USD	Softcover, 200 pages	1929	https://covers.openlibrary.org/b/id/7182819-L.jpg
30	2025-01-07 12:19:39.620672+00	2025-01-07 12:19:39.620672+00	\N	Robinson Crusoe	An insightful exploration of a fascinating subject.	12.50 USD	Softcover, 200 pages	1818	https://covers.openlibrary.org/b/id/8783768-L.jpg
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, created_at, updated_at) FROM stdin;
1	mans	Mansur@gmail.com	\N	\N
2	m	m@gmail.com	\N	\N
3	Rostislav	prrostik@gmail.com	\N	\N
4	mansur	Mansurserikov889@gmail.com	\N	\N
\.


--
-- Name: books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.books_id_seq', 30, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 4, true);


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- Name: books uni_books_title; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT uni_books_title UNIQUE (title);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_books_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_books_deleted_at ON public.books USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

