-- Удаляем функцию validate_transaction_reference, если она существует
DROP FUNCTION IF EXISTS public.validate_transaction_reference CASCADE;

-- Удаляем функцию update_transactions, если она существует
DROP FUNCTION IF EXISTS public.update_transactions CASCADE;

-- Удаляем функцию delete_expired_codes, если она существует
DROP FUNCTION IF EXISTS public.delete_expired_codes CASCADE;

-- Удаляем тип planned, если он существует
DROP TYPE IF EXISTS public.planned CASCADE;
