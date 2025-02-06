ALTER TABLE public.goals ADD COLUMN is_exceeded bool default false NOT NULL;
ALTER TABLE public.goals DROP COLUMN additional_months;