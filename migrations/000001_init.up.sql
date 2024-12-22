SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: planned; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.planned AS ENUM (
    '0',
    '1'
);


ALTER TYPE public.planned OWNER TO postgres;

--
-- Name: delete_expired_codes(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.delete_expired_codes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM confirmation_codes
    WHERE expiration_time < NOW(); -- Удаляем записи, у которых expiration_time меньше текущего времени
    RETURN NEW; -- Возвращаем новую запись (если это операция вставки или обновления)
END;
$$;


ALTER FUNCTION public.delete_expired_codes() OWNER TO postgres;

--
-- Name: update_transactions(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_transactions() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' AND TG_TABLE_NAME = 'income' THEN
        -- Предположим, что в таблице "income" есть поле "id", которое является идентификатором новой записи
        INSERT INTO transactions (user_id, amount, date, transaction_type, reference_id)
        VALUES (NEW.user_id, NEW.amount, NEW.date, 'income', NEW.id);
    ELSIF TG_OP = 'INSERT' AND TG_TABLE_NAME = 'expense' THEN
        -- Проверяем, существует ли столбец "description" в таблице "expense"
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'expense' AND column_name = 'description') THEN
            -- Используем "description" только если он существует
            INSERT INTO transactions (user_id, amount, description, date, transaction_type, reference_id)
            VALUES (NEW.user_id, NEW.amount, NEW.description, NEW.date, 'expense', NEW.id);
        ELSE
            -- Иначе, используем без "description"
            INSERT INTO transactions (user_id, amount, date, transaction_type, reference_id)
            VALUES (NEW.user_id, NEW.amount, NEW.date, 'expense', NEW.id);
        END IF;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_transactions() OWNER TO postgres;

--
-- Name: validate_transaction_reference(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.validate_transaction_reference() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.transaction_type = 'income' AND NOT EXISTS (SELECT 1 FROM income WHERE id = NEW.reference_id) THEN
        RAISE EXCEPTION 'Invalid reference_id for income transaction';
    END IF;

    -- Добавьте здесь другие проверки, если необходимо

    RETURN NEW;
END;
$$;


ALTER FUNCTION public.validate_transaction_reference() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;
