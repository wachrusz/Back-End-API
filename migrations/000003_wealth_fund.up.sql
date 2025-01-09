CREATE TYPE public.active_type AS ENUM (
    'investment',
    'saving',
    'loan'
);

ALTER TYPE public.active_type OWNER TO postgres;

ALTER TABLE public.wealth_fund ADD COLUMN type public.active_type, ADD COLUMN is_liquid BOOLEAN;
ALTER TABLE public.expense ADD COLUMN type public.active_type, ADD COLUMN is_liquid BOOLEAN;
ALTER TABLE public.income ADD COLUMN type public.active_type, ADD COLUMN is_liquid BOOLEAN;


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
             WHERE exchange_rates.currency_code = expense.currency_code),
            1)
    END AS amount_in_rubles
FROM expense;
