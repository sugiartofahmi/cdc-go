CREATE TABLE IF NOT EXISTS public.roles (
    id         uuid        DEFAULT uuid_generate_v4() NOT NULL,
    name       varchar(255)                           NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz,
    CONSTRAINT roles_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON public.roles (name) WHERE deleted_at IS NULL;
