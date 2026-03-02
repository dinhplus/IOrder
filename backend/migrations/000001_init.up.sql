CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    logo_url    TEXT,
    timezone    VARCHAR(50) NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    currency    CHAR(3) NOT NULL DEFAULT 'VND',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    settings    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE staff (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    cognito_sub TEXT,
    full_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255),
    role        VARCHAR(30) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE floor_plans (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        VARCHAR(100) NOT NULL,
    floor_level INT NOT NULL DEFAULT 1,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE restaurant_tables (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    floor_plan_id UUID NOT NULL REFERENCES floor_plans(id),
    name          VARCHAR(50) NOT NULL,
    capacity      INT NOT NULL DEFAULT 4,
    pos_x         FLOAT NOT NULL DEFAULT 0,
    pos_y         FLOAT NOT NULL DEFAULT 0,
    shape         VARCHAR(20) NOT NULL DEFAULT 'rectangle',
    status        VARCHAR(20) NOT NULL DEFAULT 'available',
    qr_token      TEXT,
    qr_expires_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE menu_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20) NOT NULL DEFAULT 'food' CHECK (type IN ('food', 'beverage')),
    description TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE menu_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id),
    category_id  UUID NOT NULL REFERENCES menu_categories(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    price        BIGINT NOT NULL,
    image_url    TEXT,
    tags         TEXT[] NOT NULL DEFAULT '{}',
    is_available BOOLEAN NOT NULL DEFAULT true,
    sort_order   INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE customers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id),
    cognito_sub    TEXT,
    phone          VARCHAR(20),
    email          VARCHAR(255),
    full_name      VARCHAR(255),
    tier           VARCHAR(20) NOT NULL DEFAULT 'bronze',
    points_balance INT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, phone),
    UNIQUE(tenant_id, email)
);

CREATE TABLE orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    table_id        UUID NOT NULL REFERENCES restaurant_tables(id),
    session_id      UUID NOT NULL,
    status          VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    customer_id     UUID REFERENCES customers(id),
    subtotal        BIGINT NOT NULL DEFAULT 0,
    discount_amount BIGINT NOT NULL DEFAULT 0,
    total           BIGINT NOT NULL DEFAULT 0,
    notes           TEXT,
    placed_at       TIMESTAMPTZ,
    confirmed_at    TIMESTAMPTZ,
    ready_at        TIMESTAMPTZ,
    served_at       TIMESTAMPTZ,
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID NOT NULL REFERENCES orders(id),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    item_id     UUID NOT NULL REFERENCES menu_items(id),
    item_name   VARCHAR(255) NOT NULL,
    item_price  BIGINT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    modifiers   JSONB NOT NULL DEFAULT '[]',
    notes       TEXT,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    routed_to   VARCHAR(20) NOT NULL DEFAULT 'kitchen' CHECK (routed_to IN ('kitchen', 'bar'))
);

CREATE TABLE order_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID NOT NULL REFERENCES orders(id),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    from_status VARCHAR(30),
    to_status   VARCHAR(30) NOT NULL,
    actor_id    UUID,
    actor_type  VARCHAR(20),
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
