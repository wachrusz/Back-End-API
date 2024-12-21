--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)

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
-- Name: delete_expired_codes(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.delete_expired_codes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM confirmation_codes
    WHERE expiration_time < NOW(); -- Удаляем записи, у которых expiration_time меньше текущего времени
    RETURN NEW; -- Возвращаем новую запись (если это операция вставки или обновления)
END;
$$;


ALTER FUNCTION public.delete_expired_codes() OWNER TO postgres;

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
    token character varying(8000),
    attempts integer DEFAULT 0,
    locked_until timestamp with time zone DEFAULT now()
);


ALTER TABLE public.confirmation_codes OWNER TO postgres;

--
-- Name: connected_accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.connected_accounts (
    id integer NOT NULL,
    user_id integer DEFAULT 0,
    bank_id integer DEFAULT 0,
    account_number character varying(20) DEFAULT '00000000000000000000'::character varying,
    account_type character varying(50) DEFAULT 'unknown'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    name character varying(50) DEFAULT 'null'::character varying,
    currency character varying(3) DEFAULT 'RUB'::character varying,
    state real DEFAULT 0
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
-- Name: currency; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.currency (
    id integer NOT NULL,
    currency_code character varying(10),
    num_code character varying(255) DEFAULT 'XXX'::character varying NOT NULL,
    nominal integer DEFAULT 1 NOT NULL,
    name character varying(255) DEFAULT 'XXX'::character varying NOT NULL,
    value double precision DEFAULT 1 NOT NULL,
    previous double precision DEFAULT 1 NOT NULL,
    cbr_id character varying(20)
);


ALTER TABLE public.currency OWNER TO postgres;

--
-- Name: currency_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.currency_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.currency_id_seq OWNER TO postgres;

--
-- Name: currency_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.currency_id_seq OWNED BY public.currency.id;


--
-- Name: exchange_rates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.exchange_rates (
    currency_code character varying(10) NOT NULL,
    rate_to_ruble numeric(10,4) NOT NULL
);


ALTER TABLE public.exchange_rates OWNER TO postgres;

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
    transaction_type character varying(255) DEFAULT 'expense'::character varying,
    currency_code character varying(3) DEFAULT 'RUB'::character varying,
    connected_account character varying(20) DEFAULT '00000000000000000000'::character varying,
    sent_to character varying(300) DEFAULT 'blank'::character varying
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
    need real DEFAULT 0,
    current_state real DEFAULT 0,
    currency character varying(3) DEFAULT 'RUB'::character varying,
    start_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    end_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
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
    description character varying(255) DEFAULT 'blank'::character varying,
    sender character varying(300) DEFAULT 'blank'::character varying,
    connected_account character varying(20) DEFAULT '00000000000000000000'::character varying,
    currency_code character varying(3) DEFAULT 'RUB'::character varying
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
-- Name: profile_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.profile_images (
    id integer NOT NULL,
    profile_id integer NOT NULL,
    image_data bytea,
    url character varying(700)
);


ALTER TABLE public.profile_images OWNER TO postgres;

--
-- Name: profile_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.profile_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.profile_images_id_seq OWNER TO postgres;

--
-- Name: profile_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.profile_images_id_seq OWNED BY public.profile_images.id;


--
-- Name: service_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_images (
    id integer NOT NULL,
    image_data bytea,
    url character varying(1000),
    service_id integer
);


ALTER TABLE public.service_images OWNER TO postgres;

--
-- Name: service_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_images_id_seq OWNER TO postgres;

--
-- Name: service_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_images_id_seq OWNED BY public.service_images.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    session_id integer NOT NULL,
    email character varying(50) NOT NULL,
    device_id character varying(128) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_activity timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    user_id integer,
    token character varying(8000) NOT NULL,
    expires_at timestamp without time zone
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
    CONSTRAINT transactions_transaction_type_check CHECK (((transaction_type)::text = ANY (ARRAY[('income'::character varying)::text, ('expense'::character varying)::text])))
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
    email character varying(114) NOT NULL,
    hashed_password character varying(150) NOT NULL,
    name character varying(50) DEFAULT 'Мы тебя не знаем...'::character varying,
    surname character varying(50) DEFAULT '...'::character varying
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
    planned public.planned,
    currency_code character varying(3) DEFAULT 'RUB'::character varying,
    connected_account character varying(100) DEFAULT '00000000000000000000'::character varying,
    category_id integer DEFAULT 6
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
-- Name: currency id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency ALTER COLUMN id SET DEFAULT nextval('public.currency_id_seq'::regclass);


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
-- Name: profile_images id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile_images ALTER COLUMN id SET DEFAULT nextval('public.profile_images_id_seq'::regclass);


--
-- Name: service_images id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_images ALTER COLUMN id SET DEFAULT nextval('public.service_images_id_seq'::regclass);


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
-- Name: banks banks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banks
    ADD CONSTRAINT banks_pkey PRIMARY KEY (id);


--
-- Name: currency cbr_id_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency
    ADD CONSTRAINT cbr_id_unique UNIQUE (cbr_id);


--
-- Name: connected_accounts connected_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_pkey PRIMARY KEY (id);


--
-- Name: currency currency_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency
    ADD CONSTRAINT currency_pkey PRIMARY KEY (id);


--
-- Name: exchange_rates exchange_rates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_rates
    ADD CONSTRAINT exchange_rates_pkey PRIMARY KEY (currency_code);


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
-- Name: profile_images profile_images_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile_images
    ADD CONSTRAINT profile_images_pkey PRIMARY KEY (id);


--
-- Name: profile_images profile_images_profile_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile_images
    ADD CONSTRAINT profile_images_profile_id_key UNIQUE (profile_id);


--
-- Name: service_images service_images_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_images
    ADD CONSTRAINT service_images_pkey PRIMARY KEY (id);


--
-- Name: service_images service_images_service_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_images
    ADD CONSTRAINT service_images_service_id_key UNIQUE (service_id);


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
-- Name: connected_accounts unique_account_number; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT unique_account_number UNIQUE (account_number);


--
-- Name: currency unique_currency_code; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency
    ADD CONSTRAINT unique_currency_code UNIQUE (currency_code);


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
-- Name: confirmation_codes trigger_delete_expired_codes; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_delete_expired_codes BEFORE INSERT OR UPDATE ON public.confirmation_codes FOR EACH ROW EXECUTE FUNCTION public.delete_expired_codes();


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
-- Name: expense expense_connected_account_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_connected_account_fkey FOREIGN KEY (connected_account) REFERENCES public.connected_accounts(account_number);


--
-- Name: expense expense_currency_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_currency_code_fkey FOREIGN KEY (currency_code) REFERENCES public.currency(currency_code);


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
-- Name: income income_connected_account_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_connected_account_fkey FOREIGN KEY (connected_account) REFERENCES public.connected_accounts(account_number);


--
-- Name: income income_currency_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_currency_code_fkey FOREIGN KEY (currency_code) REFERENCES public.currency(currency_code);


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
-- Name: profile_images profile_images_profile_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile_images
    ADD CONSTRAINT profile_images_profile_id_fkey FOREIGN KEY (profile_id) REFERENCES public.users(id) ON DELETE CASCADE;


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
-- Name: wealth_fund wealth_fund_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.investment_categories(id);


--
-- Name: wealth_fund wealth_fund_currency_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_currency_code_fkey FOREIGN KEY (currency_code) REFERENCES public.currency(currency_code);


--
-- Name: wealth_fund wealth_fund_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

