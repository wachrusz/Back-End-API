CREATE TYPE public.active_type AS ENUM (
    'investment',
    'saving',
    'loan'
);

ALTER TYPE public.active_type OWNER TO postgres;

ALTER TABLE public.wealth_fund ADD COLUMN type public.active_type, ADD COLUMN is_liquid BOOLEAN;
-- меняем тип и переносим данные
BEGIN;

-- Шаг 1: Добавить временную колонку с типом boolean
ALTER TABLE public.wealth_fund
ADD COLUMN planned_temp boolean DEFAULT false NOT NULL;

-- Шаг 2: Перенести данные в новую колонку
UPDATE public.wealth_fund
SET planned_temp = CASE
    WHEN planned = '1' THEN true
    ELSE false
END;

-- Шаг 3: Удалить старую колонку
ALTER TABLE public.wealth_fund
DROP COLUMN planned;

-- Шаг 4: Переименовать временную колонку в исходное имя
ALTER TABLE public.wealth_fund
RENAME COLUMN planned_temp TO planned;

-- Шаг 5: Удалить тип ENUM 'planned'
DROP TYPE public.planned;

COMMIT;


ALTER TABLE public.expense ADD COLUMN type public.active_type;
ALTER TABLE public.income ADD COLUMN type public.active_type;


CREATE VIEW public.expense_in_rubles AS
SELECT
    *,
    CASE
        WHEN currency_code = 'RUB' THEN amount
        ELSE amount * COALESCE(
            (SELECT rate_to_ruble
             FROM exchange_rates
             WHERE exchange_rates.currency_code = expense.currency_code),
            1)
    END AS amount_in_rubles
FROM expense;

CREATE VIEW public.income_in_rubles AS
SELECT
    *,
    CASE
        WHEN currency_code = 'RUB' THEN amount
        ELSE amount * COALESCE(
            (SELECT rate_to_ruble
             FROM exchange_rates
             WHERE exchange_rates.currency_code = income.currency_code),
            1)
    END AS amount_in_rubles
FROM income;

CREATE VIEW public.wealth_fund_in_rubles AS
SELECT
    *,
    CASE
        WHEN currency_code = 'RUB' THEN amount
        ELSE amount * COALESCE(
            (SELECT rate_to_ruble
             FROM exchange_rates
             WHERE exchange_rates.currency_code = wealth_fund.currency_code),
            1)
    END AS amount_in_rubles
FROM wealth_fund;
