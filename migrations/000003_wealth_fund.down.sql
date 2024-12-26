-- Удаление столбцов из таблицы
ALTER TABLE public.wealth_fund DROP COLUMN IF EXISTS type CASCADE, DROP COLUMN IF EXISTS is_liquid CASCADE;

-- Удаление типа ENUM
DROP TYPE IF EXISTS public.active_type CASCADE;

DROP VIEW IF EXISTS public.wealth_fund_in_rubles;
DROP VIEW IF EXISTS public.income_in_rubles;
DROP VIEW IF EXISTS public.expense_in_rubles;
