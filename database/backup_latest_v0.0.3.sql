--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1 (Debian 16.1-1)
-- Dumped by pg_dump version 16.1 (Debian 16.1-1)

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
-- Name: planned; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.planned AS ENUM (
    '0',
    '1'
);


ALTER TYPE public.planned OWNER TO postgres;

--
-- Name: update_transactions(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_transactions() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' AND TG_TABLE_NAME = 'income' THEN
        -- Предположим, что в таблице "income" есть поле "id", которое является идентификатором новой записи
        INSERT INTO transactions (user_id, amount, date, transaction_type, reference_id)
        VALUES (NEW.user_id, NEW.amount, NEW.date, 'income', NEW.id);
    ELSIF TG_OP = 'INSERT' AND TG_TABLE_NAME = 'expense' THEN
        -- Проверяем, существует ли столбец "description" в таблице "expense"
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'expense' AND column_name = 'description') THEN
            -- Используем "description" только если он существует
            INSERT INTO transactions (user_id, amount, description, date, transaction_type, reference_id)
            VALUES (NEW.user_id, NEW.amount, NEW.description, NEW.date, 'expense', NEW.id);
        ELSE
            -- Иначе, используем без "description"
            INSERT INTO transactions (user_id, amount, date, transaction_type, reference_id)
            VALUES (NEW.user_id, NEW.amount, NEW.date, 'expense', NEW.id);
        END IF;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_transactions() OWNER TO postgres;

--
-- Name: validate_transaction_reference(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.validate_transaction_reference() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.transaction_type = 'income' AND NOT EXISTS (SELECT 1 FROM income WHERE id = NEW.reference_id) THEN
        RAISE EXCEPTION 'Invalid reference_id for income transaction';
    END IF;

    -- Добавьте здесь другие проверки, если необходимо

    RETURN NEW;
END;
$$;


ALTER FUNCTION public.validate_transaction_reference() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: banks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.banks (
    id integer NOT NULL,
    name character varying(255),
    icon character varying(255)
);


ALTER TABLE public.banks OWNER TO postgres;

--
-- Name: banks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.banks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.banks_id_seq OWNER TO postgres;

--
-- Name: banks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.banks_id_seq OWNED BY public.banks.id;


--
-- Name: confirmation_codes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.confirmation_codes (
    email character varying(255),
    code character varying(10),
    expiration_time timestamp without time zone,
    token character varying(8000)
);


ALTER TABLE public.confirmation_codes OWNER TO postgres;

--
-- Name: connected_accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.connected_accounts (
    id integer NOT NULL,
    user_id integer,
    bank_id integer,
    account_number character varying(20),
    account_type character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.connected_accounts OWNER TO postgres;

--
-- Name: connected_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.connected_accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.connected_accounts_id_seq OWNER TO postgres;

--
-- Name: connected_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.connected_accounts_id_seq OWNED BY public.connected_accounts.id;


--
-- Name: expense; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expense (
    id integer NOT NULL,
    amount numeric,
    date date,
    planned boolean,
    user_id integer,
    category integer,
    transaction_type character varying(255) DEFAULT 'expense'::character varying
);


ALTER TABLE public.expense OWNER TO postgres;

--
-- Name: expense_categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expense_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255) NOT NULL,
    is_fixed boolean NOT NULL,
    user_id integer
);


ALTER TABLE public.expense_categories OWNER TO postgres;

--
-- Name: expense_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.expense_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expense_categories_id_seq OWNER TO postgres;

--
-- Name: expense_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.expense_categories_id_seq OWNED BY public.expense_categories.id;


--
-- Name: expense_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.expense_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expense_id_seq OWNER TO postgres;

--
-- Name: expense_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.expense_id_seq OWNED BY public.expense.id;


--
-- Name: goal; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.goal (
    id integer NOT NULL,
    goal character varying(255),
    user_id integer,
    need real,
    current_state real
);


ALTER TABLE public.goal OWNER TO postgres;

--
-- Name: goal_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.goal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.goal_id_seq OWNER TO postgres;

--
-- Name: goal_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.goal_id_seq OWNED BY public.goal.id;


--
-- Name: income; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income (
    id integer NOT NULL,
    amount numeric,
    date date,
    planned boolean,
    user_id integer,
    category integer,
    transaction_type character varying(255) DEFAULT 'income'::character varying,
    description character varying(255) DEFAULT 'blank'::character varying
);


ALTER TABLE public.income OWNER TO postgres;

--
-- Name: income_categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255) NOT NULL,
    is_fixed boolean NOT NULL,
    user_id integer
);


ALTER TABLE public.income_categories OWNER TO postgres;

--
-- Name: income_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_categories_id_seq OWNER TO postgres;

--
-- Name: income_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_categories_id_seq OWNED BY public.income_categories.id;


--
-- Name: income_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_id_seq OWNER TO postgres;

--
-- Name: income_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_id_seq OWNED BY public.income.id;


--
-- Name: investment_categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.investment_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255) NOT NULL,
    is_fixed boolean NOT NULL,
    user_id integer
);


ALTER TABLE public.investment_categories OWNER TO postgres;

--
-- Name: investment_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.investment_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.investment_categories_id_seq OWNER TO postgres;

--
-- Name: investment_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.investment_categories_id_seq OWNED BY public.investment_categories.id;


--
-- Name: operations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.operations (
    id integer NOT NULL,
    user_id integer,
    description character varying(255),
    amount numeric(18,2),
    date timestamp without time zone,
    category character varying(255),
    operation_type character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.operations OWNER TO postgres;

--
-- Name: operations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.operations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.operations_id_seq OWNER TO postgres;

--
-- Name: operations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.operations_id_seq OWNED BY public.operations.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    session_id integer NOT NULL,
    email character varying(50) NOT NULL,
    device_id character varying(36) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_activity timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    user_id integer,
    token character varying(8000) NOT NULL
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: sessions_session_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sessions_session_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sessions_session_id_seq OWNER TO postgres;

--
-- Name: sessions_session_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sessions_session_id_seq OWNED BY public.sessions.session_id;


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.subscriptions (
    id integer NOT NULL,
    user_id integer,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    is_active boolean NOT NULL
);


ALTER TABLE public.subscriptions OWNER TO postgres;

--
-- Name: subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.subscriptions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.subscriptions_id_seq OWNER TO postgres;

--
-- Name: subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.subscriptions_id_seq OWNED BY public.subscriptions.id;


--
-- Name: tracking_state; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tracking_state (
    id integer NOT NULL,
    state real,
    user_id integer
);


ALTER TABLE public.tracking_state OWNER TO postgres;

--
-- Name: tracking_state_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tracking_state_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tracking_state_id_seq OWNER TO postgres;

--
-- Name: tracking_state_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tracking_state_id_seq OWNED BY public.tracking_state.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    user_id integer,
    amount numeric NOT NULL,
    description character varying(255),
    date date NOT NULL,
    transaction_type character varying(20),
    reference_id integer,
    CONSTRAINT transactions_transaction_type_check CHECK (((transaction_type)::text = ANY ((ARRAY['income'::character varying, 'expense'::character varying])::text[])))
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_id_seq OWNER TO postgres;

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(50) NOT NULL,
    hashed_password character varying(8000) NOT NULL,
    name character varying(50) DEFAULT 'Мы тебя не знаем...'::character varying
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
-- Name: wealth_fund; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wealth_fund (
    id integer NOT NULL,
    amount numeric,
    date date,
    user_id integer,
    planned public.planned
);


ALTER TABLE public.wealth_fund OWNER TO postgres;

--
-- Name: wealth_fund_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.wealth_fund_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.wealth_fund_id_seq OWNER TO postgres;

--
-- Name: wealth_fund_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.wealth_fund_id_seq OWNED BY public.wealth_fund.id;


--
-- Name: banks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banks ALTER COLUMN id SET DEFAULT nextval('public.banks_id_seq'::regclass);


--
-- Name: connected_accounts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts ALTER COLUMN id SET DEFAULT nextval('public.connected_accounts_id_seq'::regclass);


--
-- Name: expense id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense ALTER COLUMN id SET DEFAULT nextval('public.expense_id_seq'::regclass);


--
-- Name: expense_categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_categories ALTER COLUMN id SET DEFAULT nextval('public.expense_categories_id_seq'::regclass);


--
-- Name: goal id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goal ALTER COLUMN id SET DEFAULT nextval('public.goal_id_seq'::regclass);


--
-- Name: income id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income ALTER COLUMN id SET DEFAULT nextval('public.income_id_seq'::regclass);


--
-- Name: income_categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_categories ALTER COLUMN id SET DEFAULT nextval('public.income_categories_id_seq'::regclass);


--
-- Name: investment_categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.investment_categories ALTER COLUMN id SET DEFAULT nextval('public.investment_categories_id_seq'::regclass);


--
-- Name: operations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operations ALTER COLUMN id SET DEFAULT nextval('public.operations_id_seq'::regclass);


--
-- Name: sessions session_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions ALTER COLUMN session_id SET DEFAULT nextval('public.sessions_session_id_seq'::regclass);


--
-- Name: subscriptions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subscriptions ALTER COLUMN id SET DEFAULT nextval('public.subscriptions_id_seq'::regclass);


--
-- Name: tracking_state id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tracking_state ALTER COLUMN id SET DEFAULT nextval('public.tracking_state_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: wealth_fund id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wealth_fund ALTER COLUMN id SET DEFAULT nextval('public.wealth_fund_id_seq'::regclass);


--
-- Data for Name: banks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.banks (id, name, icon) FROM stdin;
1	Сбербанк	https://example.com/sberbank_icon
2	ВТБ	https://example.com/vtb_icon
3	Газпромбанк	https://example.com/gazprombank_icon
4	Альфа-Банк	https://example.com/alfabank_icon
5	Тинькофф Банк	https://example.com/tinkoff_icon
6	Райффайзенбанк	https://example.com/raiffeisenbank_icon
7	ЮниКредит Банк	https://example.com/unicredit_icon
8	Связь-Банк	https://example.com/svyazbank_icon
9	Росбанк	https://example.com/rosbank_icon
10	Открытие Банк	https://example.com/otkritiebank_icon
\.


--
-- Data for Name: confirmation_codes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.confirmation_codes (email, code, expiration_time, token) FROM stdin;
\.


--
-- Data for Name: connected_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.connected_accounts (id, user_id, bank_id, account_number, account_type, created_at, updated_at) FROM stdin;
3	6	6	02109012400412010000	debit	2023-12-13 00:20:31.980632	2023-12-13 00:20:31.980632
\.


--
-- Data for Name: expense; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.expense (id, amount, date, planned, user_id, category, transaction_type) FROM stdin;
3	4000	2023-11-21	f	6	\N	expense
4	0	2023-12-13	t	6	1	expense
\.


--
-- Data for Name: expense_categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.expense_categories (id, name, icon, is_fixed, user_id) FROM stdin;
1	cat1	icon1	f	6
2	cat1	icon1	f	6
\.


--
-- Data for Name: goal; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.goal (id, goal, user_id, need, current_state) FROM stdin;
1	NEWGOAL	6	100	5.2
2	NEWGOAL	6	100	5.2
\.


--
-- Data for Name: income; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.income (id, amount, date, planned, user_id, category, transaction_type, description) FROM stdin;
3	2000	2023-11-21	f	6	\N	income	blank
4	3000	2023-11-21	f	6	\N	income	blank
6	0	2023-12-13	t	6	1	income	blank
7	0	2023-12-13	t	6	1	income	blank
8	0	2023-12-13	t	6	1	income	blank
9	0	2023-12-13	t	16	1	income	blank
10	20000	2023-12-13	t	16	1	income	blank
24	212300	2023-12-16	f	17	1	income	blank
25	2000	2023-01-01	t	17	1	income	blank
26	2000	2023-01-01	t	17	1	income	blank
27	13213212000	2023-01-01	t	17	1	income	blank
\.


--
-- Data for Name: income_categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.income_categories (id, name, icon, is_fixed, user_id) FROM stdin;
1	cat1	icon1	f	6
\.


--
-- Data for Name: investment_categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.investment_categories (id, name, icon, is_fixed, user_id) FROM stdin;
1	cat1	icon1	f	6
\.


--
-- Data for Name: operations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.operations (id, user_id, description, amount, date, category, operation_type, created_at) FROM stdin;
1	6	Доход	0.00	2023-12-13 00:00:00	1	1	2023-12-13 00:54:48.018793
2	6	Расход	0.00	2023-12-13 00:00:00	1	1	2023-12-13 00:55:00.26963
3	6	Доход	0.00	2023-12-13 00:00:00	1	1	2023-12-13 01:01:52.291373
4	6	Доход	0.00	2023-12-13 00:00:00	1	1	2023-12-13 01:02:27.189658
5	16	Доход	0.00	2023-12-13 00:00:00	1	1	2023-12-13 01:21:36.39496
6	16	Доход	20000.00	2023-12-13 00:00:00	1	1	2023-12-13 01:22:38.13441
7	17	Доход	2000.00	2023-01-01 00:00:00	1	1	2023-12-16 02:13:03.60394
8	17	Доход	13213212000.00	2023-01-01 00:00:00	1	1	2023-12-16 02:13:19.642
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (session_id, email, device_id, created_at, last_activity, user_id, token) FROM stdin;
77	newuser	::1_curl/8.4.0	2023-12-09 02:02:50.656611+03	2023-12-09 02:02:50.656611+03	6	Gos1li3EufoPNdBGzFLPZqSYnwa-gc4G0arxrkHCv16Oeq6cJde4O5br8LckpkcgJzg-lF_e8juJn3gnyCLi4GaUp6lzbeCH4D5__thdOyuKLGJ8xDhQqFpZk2731NlSCDB8cvbDmlDxS1H5L0EfkQ==
78	minka@yandex.ru	::1_curl/8.4.0	2023-12-12 23:26:41.500603+03	2023-12-12 23:26:41.500603+03	16	uKldJVbyBukueuyXAzhXSk2zncDN3lbCEhTBZURXA0MMAYlsF3_ToEB18vLczwreF74eR9dzgdpQYk0n51Bg7FRxfQEZN7RETnYPzhprbS-fzK8qpbdFgOPcfvv7M5qZmp64jFf2EF58JPNS3ep2pyk=
79	lstwrd@yandex.ru	::1_curl/8.4.0	2023-12-16 01:33:35.882856+03	2023-12-16 01:33:35.882856+03	17	ppujL3POqxBOIyW5EVLXe_Nc9mxXyD3ZTEstMQwZLUyr-IL-nat8j6nT2XRNAi-6Xh4z6_BGjahqh7YbU56I7HmDd9KDGV4Bk2v5JYp6q13_A9gnoQL8qsyNB7r5rfDzbYyCmWWP1wYqkqLg_GbMyFk=
\.


--
-- Data for Name: subscriptions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.subscriptions (id, user_id, start_date, end_date, is_active) FROM stdin;
1	6	2023-11-22 00:00:00	2023-12-22 00:00:00	t
\.


--
-- Data for Name: tracking_state; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tracking_state (id, state, user_id) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (id, user_id, amount, description, date, transaction_type, reference_id) FROM stdin;
3	17	212300	\N	2023-12-16	income	24
4	17	2000	\N	2023-01-01	income	25
5	17	2000	\N	2023-01-01	income	26
6	17	13213212000	\N	2023-01-01	income	27
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, hashed_password, name) FROM stdin;
8	'	$2a$10$u.6Fu.v0vpYhPdrjYmeXG.LvIXRl2Rrq0h4sD0kvqg7BqvRkBCCmm	qwe
9	qwe	$2a$10$MNRuvurbCOBsF9UvDKIOJuWEdEt5bMcAWiUeoodM8S3eQcT4MzAua	zalupa
10	user	$2a$10$L4nu353gvOdSjXb38hNQKOvHRKIy6NMd8fmFVZiV8XiOZGsEb6G6.	name
11	s	$2a$10$oLnrFMaWEcWB9aQ8GODcCe2pkqLMc7HU5HVVx1sOWfe5NzV/dfbTK	fq
12	usr	$2a$10$rt588s.IEFq9XmBvfW2vVeRuaQbwbqepHikQCMba0MMLXOboJMEp.	name
6	newuser	$2a$10$QEeLTLxWTwbAxlYH6kEaTe66I7pZU3oMegufaRJmgYXqleQ3g22o2	zalupa
16	minka@yandex.ru	$2a$10$RC4iyWWaMb2TsH4oGp3lE.1jfLZsplOlcjth3FF1bFhOStMhY.BlC	Мы тебя не знаем...
17	lstwrd@yandex.ru	$2a$10$I8WGOipScEASB1/EXiQFBOyxewRSMkZaO/HvhE/0IQzxXUjuC/h/u	Гнидкинс
\.


--
-- Data for Name: wealth_fund; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wealth_fund (id, amount, date, user_id, planned) FROM stdin;
\.


--
-- Name: banks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.banks_id_seq', 10, true);


--
-- Name: connected_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.connected_accounts_id_seq', 3, true);


--
-- Name: expense_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.expense_categories_id_seq', 2, true);


--
-- Name: expense_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.expense_id_seq', 4, true);


--
-- Name: goal_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.goal_id_seq', 2, true);


--
-- Name: income_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.income_categories_id_seq', 2, true);


--
-- Name: income_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.income_id_seq', 27, true);


--
-- Name: investment_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.investment_categories_id_seq', 1, true);


--
-- Name: operations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.operations_id_seq', 8, true);


--
-- Name: sessions_session_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sessions_session_id_seq', 79, true);


--
-- Name: subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.subscriptions_id_seq', 1, true);


--
-- Name: tracking_state_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tracking_state_id_seq', 1, false);


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.transactions_id_seq', 6, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 17, true);


--
-- Name: wealth_fund_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.wealth_fund_id_seq', 1, true);


--
-- Name: banks banks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banks
    ADD CONSTRAINT banks_pkey PRIMARY KEY (id);


--
-- Name: connected_accounts connected_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_pkey PRIMARY KEY (id);


--
-- Name: expense_categories expense_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_pkey PRIMARY KEY (id);


--
-- Name: expense expense_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_pkey PRIMARY KEY (id);


--
-- Name: goal goal_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goal
    ADD CONSTRAINT goal_pkey PRIMARY KEY (id);


--
-- Name: income_categories income_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_pkey PRIMARY KEY (id);


--
-- Name: income income_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_pkey PRIMARY KEY (id);


--
-- Name: investment_categories investment_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.investment_categories
    ADD CONSTRAINT investment_categories_pkey PRIMARY KEY (id);


--
-- Name: operations operations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operations
    ADD CONSTRAINT operations_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (session_id);


--
-- Name: sessions sessions_username_device_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_username_device_id_key UNIQUE (email, device_id);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);


--
-- Name: tracking_state tracking_state_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tracking_state
    ADD CONSTRAINT tracking_state_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: connected_accounts unique_user_bank_account; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT unique_user_bank_account UNIQUE (user_id, bank_id, account_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (email);


--
-- Name: wealth_fund wealth_fund_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_pkey PRIMARY KEY (id);


--
-- Name: transactions transaction_reference_check; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER transaction_reference_check BEFORE INSERT OR UPDATE ON public.transactions FOR EACH ROW EXECUTE FUNCTION public.validate_transaction_reference();


--
-- Name: income update_transactions_trigger; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_transactions_trigger AFTER INSERT OR UPDATE ON public.income FOR EACH ROW EXECUTE FUNCTION public.update_transactions();


--
-- Name: expense update_transactions_trigger_expenses; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_transactions_trigger_expenses AFTER INSERT OR UPDATE ON public.expense FOR EACH ROW EXECUTE FUNCTION public.update_transactions();


--
-- Name: connected_accounts connected_accounts_bank_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_bank_id_fkey FOREIGN KEY (bank_id) REFERENCES public.banks(id) ON DELETE CASCADE;


--
-- Name: connected_accounts connected_accounts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: expense_categories expense_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: expense expense_category_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_category_fkey FOREIGN KEY (category) REFERENCES public.expense_categories(id);


--
-- Name: expense expense_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: goal goal_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goal
    ADD CONSTRAINT goal_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: income_categories income_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: income income_category_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_category_fkey FOREIGN KEY (category) REFERENCES public.income_categories(id);


--
-- Name: income income_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: investment_categories investment_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.investment_categories
    ADD CONSTRAINT investment_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: operations operations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operations
    ADD CONSTRAINT operations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: tracking_state tracking_state_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tracking_state
    ADD CONSTRAINT tracking_state_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: transactions transactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: wealth_fund wealth_fund_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

