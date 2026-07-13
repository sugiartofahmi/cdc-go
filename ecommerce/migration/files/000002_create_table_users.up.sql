CREATE TABLE IF NOT EXISTS public.users (
    id         uuid        DEFAULT uuid_generate_v4() NOT NULL,
    name       varchar(255)                           NOT NULL,
    email      varchar(255)                           NOT NULL,
    password   varchar(255)                           NOT NULL,
    status     int         DEFAULT 1                  NOT NULL,
    role_id    uuid                                     NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz,
    CONSTRAINT users_pkey     PRIMARY KEY (id),
    CONSTRAINT fk_users_role  FOREIGN KEY (role_id) REFERENCES public.roles(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON public.users (email) WHERE deleted_at IS NULL;
