-- Удаляем индекс на пару колонок (expires_at, revoked)
DROP INDEX IF EXISTS expires_revoked_idx;

-- Удаляем уникальный индекс на колонку token
DROP INDEX IF EXISTS unique_token_idx;

DROP INDEX IF EXISTS user_devise_idx;

-- Удаляем колонку revoked
ALTER TABLE public.sessions
DROP COLUMN IF EXISTS revoked;

-- Добавляем колонку email с типом character varying(50) и NOT NULL
ALTER TABLE public.sessions
ADD COLUMN email character varying(50);

BEGIN;

ALTER TABLE public.sessions
ALTER COLUMN expires_at DROP NOT NULL;

-- Изменение типа данных колонки expires_at обратно на timestamp без time zone
ALTER TABLE public.sessions
ALTER COLUMN expires_at TYPE TIMESTAMP
USING expires_at::TIMESTAMP;

COMMIT;


-- Добавляем уникальное ограничение на пару (email, device_id)
ALTER TABLE public.sessions
ADD CONSTRAINT sessions_email_device_id_key UNIQUE (email, device_id);
