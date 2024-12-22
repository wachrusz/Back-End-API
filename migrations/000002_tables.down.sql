-- Удаляем таблицу banks, если она существует
DROP TABLE IF EXISTS public.banks CASCADE;

-- Удаляем последовательность banks_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.banks_id_seq CASCADE;

-- Удаляем таблицу confirmation_codes, если она существует
DROP TABLE IF EXISTS public.confirmation_codes CASCADE;

-- Удаляем таблицу connected_accounts, если она существует
DROP TABLE IF EXISTS public.connected_accounts CASCADE;

-- Удаляем последовательность connected_accounts_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.connected_accounts_id_seq CASCADE;

-- Удаляем таблицу currency, если она существует
DROP TABLE IF EXISTS public.currency CASCADE;

-- Удаляем последовательность currency_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.currency_id_seq CASCADE;

-- Удаляем таблицу exchange_rates, если она существует
DROP TABLE IF EXISTS public.exchange_rates CASCADE;

-- Удаляем таблицу expense, если она существует
DROP TABLE IF EXISTS public.expense CASCADE;

-- Удаляем таблицу expense_categories, если она существует
DROP TABLE IF EXISTS public.expense_categories CASCADE;

-- Удаляем последовательность expense_categories_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.expense_categories_id_seq CASCADE;

-- Удаляем последовательность expense_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.expense_id_seq CASCADE;

-- Удаляем таблицу goal, если она существует
DROP TABLE IF EXISTS public.goal CASCADE;

-- Удаляем последовательность goal_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.goal_id_seq CASCADE;

-- Удаляем таблицу income, если она существует
DROP TABLE IF EXISTS public.income CASCADE;

-- Удаляем таблицу income_categories, если она существует
DROP TABLE IF EXISTS public.income_categories CASCADE;

-- Удаляем последовательность income_categories_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.income_categories_id_seq CASCADE;

-- Удаляем последовательность income_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.income_id_seq CASCADE;

-- Удаляем таблицу investment_categories, если она существует
DROP TABLE IF EXISTS public.investment_categories CASCADE;

-- Удаляем последовательность investment_categories_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.investment_categories_id_seq CASCADE;

-- Удаляем таблицу operations, если она существует
DROP TABLE IF EXISTS public.operations CASCADE;

-- Удаляем последовательность operations_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.operations_id_seq CASCADE;

-- Удаляем таблицу profile_images, если она существует
DROP TABLE IF EXISTS public.profile_images CASCADE;

-- Удаляем последовательность profile_images_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.profile_images_id_seq CASCADE;

-- Удаляем таблицу service_images, если она существует
DROP TABLE IF EXISTS public.service_images CASCADE;

-- Удаляем последовательность service_images_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.service_images_id_seq CASCADE;

-- Удаляем таблицу sessions, если она существует
DROP TABLE IF EXISTS public.sessions CASCADE;

-- Удаляем последовательность sessions_session_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.sessions_session_id_seq CASCADE;

-- Удаляем таблицу subscriptions, если она существует
DROP TABLE IF EXISTS public.subscriptions CASCADE;

-- Удаляем последовательность subscriptions_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.subscriptions_id_seq CASCADE;

-- Удаляем таблицу tracking_state, если она существует
DROP TABLE IF EXISTS public.tracking_state CASCADE;

-- Удаляем последовательность tracking_state_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.tracking_state_id_seq CASCADE;

-- Удаляем таблицу transactions, если она существует
DROP TABLE IF EXISTS public.transactions CASCADE;

-- Удаляем последовательность transactions_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.transactions_id_seq CASCADE;

-- Удаляем таблицу users, если она существует
DROP TABLE IF EXISTS public.users CASCADE;

-- Удаляем последовательность users_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.users_id_seq CASCADE;

-- Удаляем таблицу wealth_fund, если она существует
DROP TABLE IF EXISTS public.wealth_fund CASCADE;

-- Удаляем последовательность wealth_fund_id_seq, если она существует
DROP SEQUENCE IF EXISTS public.wealth_fund_id_seq CASCADE;

-- Удаляем триггер для проверки ссылки транзакции на таблице transactions, если он существует
DROP TRIGGER IF EXISTS transaction_reference_check ON public.transactions;

-- Удаляем триггер для удаления просроченных кодов на таблице confirmation_codes, если он существует
DROP TRIGGER IF EXISTS trigger_delete_expired_codes ON public.confirmation_codes;

-- Удаляем триггер для обновления транзакций на таблице income, если он существует
DROP TRIGGER IF EXISTS update_transactions_trigger ON public.income;

-- Удаляем триггер для обновления транзакций на таблице expense, если он существует
DROP TRIGGER IF EXISTS update_transactions_trigger_expenses ON public.expense;
