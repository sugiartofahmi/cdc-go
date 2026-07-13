CREATE TABLE IF NOT EXISTS public.categories (
    id         uuid        DEFAULT uuid_generate_v4() NOT NULL,
    name       varchar(255)                           NOT NULL,
    slug       varchar(255)                           NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz,
    CONSTRAINT categories_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_slug ON public.categories (slug) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON public.categories (name) WHERE deleted_at IS NULL;
