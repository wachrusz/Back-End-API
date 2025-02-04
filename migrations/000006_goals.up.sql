DROP TABLE IF EXISTS public.goal cascade;

CREATE TABLE public.goals (
    id serial primary key,
    amount numeric not null,
    currency_code     varchar(3)   default 'RUB'::character varying
        references public.currency (currency_code),
    user_id bigint references public.users on delete cascade,
    name varchar(255) default 'My goal',
    months integer not null default 1,
    additional_months smallint default 0 NOT NULL ,
    is_completed bool default false NOT NULL ,
    start_date date default CURRENT_DATE
);

ALTER TABLE public.goals owner TO postgres;

CREATE TABLE public.goal_transactions (
    id serial primary key,
    goal_id integer references goals(id) NOT NULL,
    amount            numeric NOT NULL,
    planned           boolean default false NOT NULL,
    transaction_type  varchar(255) default 'goal transaction'::character varying,
    currency_code     varchar(3)   default 'RUB'::character varying
        references public.currency (currency_code),
    date timestamp with time zone not null default CURRENT_TIMESTAMP,
    connected_account varchar(20)  default '00000000000000000000'::character varying
        references connected_accounts (account_number)
);

ALTER TABLE public.goal_transactions owner to postgres;
