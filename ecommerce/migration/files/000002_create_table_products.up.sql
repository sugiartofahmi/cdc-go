CREATE TABLE IF NOT EXISTS public.products (
    id          uuid          DEFAULT uuid_generate_v4() NOT NULL,
    name        varchar(255)                             NOT NULL,
    slug        varchar(255)                             NOT NULL,
    description text,
    price       numeric(15,2)                            NOT NULL,
    stock       integer       DEFAULT 0                  NOT NULL,
    category_id uuid                                     NOT NULL,
    created_at  timestamptz   DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamptz   DEFAULT CURRENT_TIMESTAMP,
    deleted_at  timestamptz,
    CONSTRAINT products_pkey PRIMARY KEY (id),
    CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_products_slug ON public.products (slug) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_category_id ON public.products (category_id);
