ALTER TABLE ONLY public.sessions
    DROP CONSTRAINT IF EXISTS sessions_username_device_id_key,
    DROP COLUMN IF EXISTS email;

ALTER TABLE public.sessions
ADD COLUMN revoked boolean DEFAULT FALSE;


BEGIN;

-- Изменение типа данных колонки expires_at
ALTER TABLE public.sessions
ALTER COLUMN expires_at TYPE TIMESTAMPTZ 
USING expires_at::TIMESTAMPTZ;

ALTER TABLE public.sessions
ALTER COLUMN expires_at SET NOT NULL;

COMMIT;


-- Создаем уникальный индекс на колонку token
CREATE UNIQUE INDEX unique_token_idx
ON public.sessions (token);

-- Создаем обычный индекс на пару колонок (expires_at, revoked)
CREATE INDEX expires_revoked_idx
ON public.sessions (expires_at, revoked);

CREATE INDEX user_devise_idx
ON public.sessions (user_id, device_id);