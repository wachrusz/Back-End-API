-- Удаление столбцов из таблицы
ALTER TABLE public.wealth_fund DROP COLUMN IF EXISTS type CASCADE, DROP COLUMN IF EXISTS is_liquid CASCADE;
ALTER TABLE public.expense DROP COLUMN IF EXISTS type CASCADE;
ALTER TABLE public.income DROP COLUMN IF EXISTS type CASCADE;

-- Удаление типа ENUM
DROP TYPE IF EXISTS public.active_type CASCADE;

DROP VIEW IF EXISTS public.wealth_fund_in_rubles;
DROP VIEW IF EXISTS public.income_in_rubles;
DROP VIEW IF EXISTS public.expense_in_rubles;

BEGIN;

-- Шаг 1: Восстановить тип ENUM 'planned'
CREATE TYPE public.planned AS ENUM ('0', '1');
ALTER TYPE public.planned OWNER TO postgres;

-- Шаг 2: Добавить колонку planned с типом ENUM 'planned'
ALTER TABLE public.wealth_fund
ADD COLUMN planned_temp public.planned;

-- Шаг 3: Перенести данные из boolean обратно в ENUM
UPDATE public.wealth_fund
SET planned_temp = CASE
    WHEN planned = true THEN '1'::public.planned
    WHEN planned = false THEN '0'::public.planned
END;

-- Шаг 4: Удалить колонку boolean
ALTER TABLE public.wealth_fund
DROP COLUMN planned;

-- Шаг 5: Переименовать временную колонку в исходное имя
ALTER TABLE public.wealth_fund
RENAME COLUMN planned_temp TO planned;

COMMIT;
