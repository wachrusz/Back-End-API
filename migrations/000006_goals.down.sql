DROP TABLE IF EXISTS public.goals cascade ;
DROP TABLE IF EXISTS public.goal_transactions cascade ;

create table public.goal
(
    id            serial
        primary key,
    goal          varchar(255),
    user_id       integer
        references public.users
            on delete cascade,
    need          real       default 0,
    current_state real       default 0,
    currency      varchar(3) default 'RUB'::character varying,
    start_date    timestamp  default CURRENT_TIMESTAMP,
    end_date      timestamp  default CURRENT_TIMESTAMP
);

alter table public.goal owner to postgres;
