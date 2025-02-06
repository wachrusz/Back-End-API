ALTER TABLE goals ADD COLUMN additional_months smallint default 0 NOT NULL;
ALTER TABLE public.goals DROP COLUMN is_exceeded;
