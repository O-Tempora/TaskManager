--
-- PostgreSQL database dump
--

-- Dumped from database version 14.7 (Ubuntu 14.7-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)

-- Started on 2023-04-08 16:54:33 MSK

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
-- TOC entry 3440 (class 1262 OID 16432)
-- Name: dip; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE dip WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'en_US.UTF-8';


ALTER DATABASE dip OWNER TO postgres;

\connect dip

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
-- TOC entry 3 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 3441 (class 0 OID 0)
-- Dependencies: 3
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 222 (class 1259 OID 16534)
-- Name: commentaries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.commentaries (
    id integer NOT NULL,
    content text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    person_id integer NOT NULL,
    task_id integer NOT NULL
);


ALTER TABLE public.commentaries OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 16533)
-- Name: commentaries_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.commentaries_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.commentaries_id_seq OWNER TO postgres;

--
-- TOC entry 3442 (class 0 OID 0)
-- Dependencies: 221
-- Name: commentaries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.commentaries_id_seq OWNED BY public.commentaries.id;


--
-- TOC entry 220 (class 1259 OID 16520)
-- Name: person_task; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.person_task (
    person_id integer NOT NULL,
    task_id integer NOT NULL
);


ALTER TABLE public.person_task OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16500)
-- Name: person_workspace; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.person_workspace (
    person_id integer NOT NULL,
    workspace_id integer NOT NULL,
    role_id integer NOT NULL
);


ALTER TABLE public.person_workspace OWNER TO postgres;

--
-- TOC entry 212 (class 1259 OID 16450)
-- Name: persons; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.persons (
    id integer NOT NULL,
    name text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    settings text NOT NULL,
    phone text NOT NULL
);


ALTER TABLE public.persons OWNER TO postgres;

--
-- TOC entry 211 (class 1259 OID 16449)
-- Name: persons_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.persons_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.persons_id_seq OWNER TO postgres;

--
-- TOC entry 3443 (class 0 OID 0)
-- Dependencies: 211
-- Name: persons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.persons_id_seq OWNED BY public.persons.id;


--
-- TOC entry 224 (class 1259 OID 16553)
-- Name: statuses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.statuses (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.statuses OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 16552)
-- Name: statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.statuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.statuses_id_seq OWNER TO postgres;

--
-- TOC entry 3444 (class 0 OID 0)
-- Dependencies: 223
-- Name: statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.statuses_id_seq OWNED BY public.statuses.id;


--
-- TOC entry 216 (class 1259 OID 16473)
-- Name: task_groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.task_groups (
    id integer NOT NULL,
    name text NOT NULL,
    color text NOT NULL,
    workspace_id integer NOT NULL
);


ALTER TABLE public.task_groups OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 16472)
-- Name: task_groups_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.task_groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.task_groups_id_seq OWNER TO postgres;

--
-- TOC entry 3445 (class 0 OID 0)
-- Dependencies: 215
-- Name: task_groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.task_groups_id_seq OWNED BY public.task_groups.id;


--
-- TOC entry 218 (class 1259 OID 16487)
-- Name: tasks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tasks (
    id integer NOT NULL,
    description text NOT NULL,
    created_at date NOT NULL,
    start_at date NOT NULL,
    finish_at date NOT NULL,
    group_id integer NOT NULL,
    status_id integer NOT NULL
);


ALTER TABLE public.tasks OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16486)
-- Name: tasks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tasks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tasks_id_seq OWNER TO postgres;

--
-- TOC entry 3446 (class 0 OID 0)
-- Dependencies: 217
-- Name: tasks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tasks_id_seq OWNED BY public.tasks.id;


--
-- TOC entry 210 (class 1259 OID 16441)
-- Name: user_role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_role (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.user_role OWNER TO postgres;

--
-- TOC entry 209 (class 1259 OID 16440)
-- Name: user_role_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_role_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_role_id_seq OWNER TO postgres;

--
-- TOC entry 3447 (class 0 OID 0)
-- Dependencies: 209
-- Name: user_role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_role_id_seq OWNED BY public.user_role.id;


--
-- TOC entry 214 (class 1259 OID 16464)
-- Name: workspaces; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.workspaces (
    id integer NOT NULL,
    name text NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.workspaces OWNER TO postgres;

--
-- TOC entry 213 (class 1259 OID 16463)
-- Name: workspaces_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.workspaces_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.workspaces_id_seq OWNER TO postgres;

--
-- TOC entry 3448 (class 0 OID 0)
-- Dependencies: 213
-- Name: workspaces_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.workspaces_id_seq OWNED BY public.workspaces.id;


--
-- TOC entry 3250 (class 2604 OID 16537)
-- Name: commentaries id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commentaries ALTER COLUMN id SET DEFAULT nextval('public.commentaries_id_seq'::regclass);


--
-- TOC entry 3246 (class 2604 OID 16453)
-- Name: persons id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.persons ALTER COLUMN id SET DEFAULT nextval('public.persons_id_seq'::regclass);


--
-- TOC entry 3251 (class 2604 OID 16556)
-- Name: statuses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.statuses ALTER COLUMN id SET DEFAULT nextval('public.statuses_id_seq'::regclass);


--
-- TOC entry 3248 (class 2604 OID 16476)
-- Name: task_groups id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_groups ALTER COLUMN id SET DEFAULT nextval('public.task_groups_id_seq'::regclass);


--
-- TOC entry 3249 (class 2604 OID 16490)
-- Name: tasks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks ALTER COLUMN id SET DEFAULT nextval('public.tasks_id_seq'::regclass);


--
-- TOC entry 3245 (class 2604 OID 16444)
-- Name: user_role id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_role ALTER COLUMN id SET DEFAULT nextval('public.user_role_id_seq'::regclass);


--
-- TOC entry 3247 (class 2604 OID 16467)
-- Name: workspaces id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workspaces ALTER COLUMN id SET DEFAULT nextval('public.workspaces_id_seq'::regclass);


--
-- TOC entry 3432 (class 0 OID 16534)
-- Dependencies: 222
-- Data for Name: commentaries; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.commentaries (id, content, created_at, person_id, task_id) FROM stdin;
\.


--
-- TOC entry 3430 (class 0 OID 16520)
-- Dependencies: 220
-- Data for Name: person_task; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.person_task (person_id, task_id) FROM stdin;
\.


--
-- TOC entry 3429 (class 0 OID 16500)
-- Dependencies: 219
-- Data for Name: person_workspace; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.person_workspace (person_id, workspace_id, role_id) FROM stdin;
1	1	1
\.


--
-- TOC entry 3422 (class 0 OID 16450)
-- Dependencies: 212
-- Data for Name: persons; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.persons (id, name, password, email, settings, phone) FROM stdin;
1	admin	$2a$04$qIqHAJTuakpQno7GI6t3D.X9WYtEzaSDzoDj4NMaS5vQunWWFSqIG	admin@gmail.com		+79290874840
2	testName	$2a$04$PZJFgGJ1S01ba.aIQdUzM.x1OBkSBdy3./HljkzSxAnJ7Oid/b0qO	test1@mail.ru		+79290874841
3	testName2	$2a$04$H1In8dsfEC//GnDRZ1A.v.tpVH7MjeVx5GGI/eH8ns/rDdZrdbkDq	test2@mail.ru		+79290874842
4	web	$2a$04$QyoG9Gvw4bNZmLt/JvC43ONcAJ15.MTQyNslcbxiGKEi8Q82LmbwS	web@gmail.com		+79290874844
6	Valik Valokordin	$2a$04$6KnaWpZaJApG72yrUmTNDuh2MSVbUL1mOcs7pvM5CLKwbg8zMZqha	Retardo12@yandex.ru		89290874840
7	vanya	$2a$04$0g/ou9LYXSZBTirvi0/W6.wSzB.DNScLo7DCj.DGy3qPRlVGBSoru	vanya@mail.ru		+79158477147
8	vanya2	$2a$04$x5m7DwXhtVxK4rMcw6Gs5.M2nx9u0itWEkIEMTkWCkKQXKh6DPu/G	vanya@yandex.ru		+79290874850
9	NewName	$2a$04$ErgqxOcW3CM3jF8Yv8v4JuevXT/VGdsvzzP2hwPGmxJYT.a0sD2G.	newemail@yandex.ru		+79290874849
\.


--
-- TOC entry 3434 (class 0 OID 16553)
-- Dependencies: 224
-- Data for Name: statuses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.statuses (id, name) FROM stdin;
1	test status
2	ToDo
3	Done
\.


--
-- TOC entry 3426 (class 0 OID 16473)
-- Dependencies: 216
-- Data for Name: task_groups; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.task_groups (id, name, color, workspace_id) FROM stdin;
1	ToDo	#00FF00	1
\.


--
-- TOC entry 3428 (class 0 OID 16487)
-- Dependencies: 218
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tasks (id, description, created_at, start_at, finish_at, group_id, status_id) FROM stdin;
1	Test task1	2023-03-23	2023-03-23	2023-03-26	1	1
2	Test task 2	2023-03-27	2023-03-28	2023-03-30	1	2
9	NewTask	2023-03-30	2023-03-01	2023-03-06	1	3
\.


--
-- TOC entry 3420 (class 0 OID 16441)
-- Dependencies: 210
-- Data for Name: user_role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_role (id, name) FROM stdin;
1	Admin
\.


--
-- TOC entry 3424 (class 0 OID 16464)
-- Dependencies: 214
-- Data for Name: workspaces; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.workspaces (id, name, description) FROM stdin;
1	TestSpace	Created for test purposes
\.


--
-- TOC entry 3449 (class 0 OID 0)
-- Dependencies: 221
-- Name: commentaries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.commentaries_id_seq', 1, false);


--
-- TOC entry 3450 (class 0 OID 0)
-- Dependencies: 211
-- Name: persons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.persons_id_seq', 9, true);


--
-- TOC entry 3451 (class 0 OID 0)
-- Dependencies: 223
-- Name: statuses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.statuses_id_seq', 3, true);


--
-- TOC entry 3452 (class 0 OID 0)
-- Dependencies: 215
-- Name: task_groups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.task_groups_id_seq', 1, true);


--
-- TOC entry 3453 (class 0 OID 0)
-- Dependencies: 217
-- Name: tasks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tasks_id_seq', 11, true);


--
-- TOC entry 3454 (class 0 OID 0)
-- Dependencies: 209
-- Name: user_role_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_role_id_seq', 1, true);


--
-- TOC entry 3455 (class 0 OID 0)
-- Dependencies: 213
-- Name: workspaces_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.workspaces_id_seq', 1, true);


--
-- TOC entry 3267 (class 2606 OID 16541)
-- Name: commentaries commentaries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commentaries
    ADD CONSTRAINT commentaries_pkey PRIMARY KEY (id);


--
-- TOC entry 3265 (class 2606 OID 16504)
-- Name: person_workspace person_workspace_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.person_workspace
    ADD CONSTRAINT person_workspace_pkey PRIMARY KEY (person_id, workspace_id, role_id);


--
-- TOC entry 3255 (class 2606 OID 16573)
-- Name: persons persons_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.persons
    ADD CONSTRAINT persons_email_key UNIQUE (email);


--
-- TOC entry 3257 (class 2606 OID 16457)
-- Name: persons persons_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.persons
    ADD CONSTRAINT persons_pkey PRIMARY KEY (id);


--
-- TOC entry 3269 (class 2606 OID 16560)
-- Name: statuses statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.statuses
    ADD CONSTRAINT statuses_pkey PRIMARY KEY (id);


--
-- TOC entry 3261 (class 2606 OID 16480)
-- Name: task_groups task_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_groups
    ADD CONSTRAINT task_groups_pkey PRIMARY KEY (id);


--
-- TOC entry 3263 (class 2606 OID 16494)
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- TOC entry 3253 (class 2606 OID 16448)
-- Name: user_role user_role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_role
    ADD CONSTRAINT user_role_pkey PRIMARY KEY (id);


--
-- TOC entry 3259 (class 2606 OID 16471)
-- Name: workspaces workspaces_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workspaces
    ADD CONSTRAINT workspaces_pkey PRIMARY KEY (id);


--
-- TOC entry 3271 (class 2606 OID 16495)
-- Name: tasks fk_group; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES public.task_groups(id);


--
-- TOC entry 3276 (class 2606 OID 16523)
-- Name: person_task fk_pers; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.person_task
    ADD CONSTRAINT fk_pers FOREIGN KEY (person_id) REFERENCES public.persons(id);


--
-- TOC entry 3274 (class 2606 OID 16510)
-- Name: person_workspace fk_person; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.person_workspace
    ADD CONSTRAINT fk_person FOREIGN KEY (person_id) REFERENCES public.persons(id);


--
-- TOC entry 3278 (class 2606 OID 16542)
-- Name: commentaries fk_person; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commentaries
    ADD CONSTRAINT fk_person FOREIGN KEY (person_id) REFERENCES public.persons(id);


--
-- TOC entry 3273 (class 2606 OID 16505)
-- Name: person_workspace fk_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.person_workspace
    ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES public.user_role(id);


--
-- TOC entry 3272 (class 2606 OID 16561)
-- Name: tasks fk_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT fk_status FOREIGN KEY (status_id) REFERENCES public.statuses(id);


--
-- TOC entry 3277 (class 2606 OID 16528)
-- Name: person_task fk_task; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.person_task
    ADD CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES public.tasks(id);


--
-- TOC entry 3279 (class 2606 OID 16547)
-- Name: commentaries fk_task; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commentaries
    ADD CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES public.tasks(id);


--
-- TOC entry 3270 (class 2606 OID 16481)
-- Name: task_groups fk_workspace; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_groups
    ADD CONSTRAINT fk_workspace FOREIGN KEY (workspace_id) REFERENCES public.workspaces(id);


--
-- TOC entry 3275 (class 2606 OID 16515)
-- Name: person_workspace fk_workspace; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.person_workspace
    ADD CONSTRAINT fk_workspace FOREIGN KEY (workspace_id) REFERENCES public.workspaces(id);


-- Completed on 2023-04-08 16:54:33 MSK

--
-- PostgreSQL database dump complete
--

