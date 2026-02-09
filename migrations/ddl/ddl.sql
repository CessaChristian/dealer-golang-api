CREATE TABLE public.users (
	user_id bigserial NOT NULL,
	"name" varchar(150) NOT NULL,
	email varchar(150) NOT NULL,
	"password" text NOT NULL,
	"role" varchar(50) DEFAULT 'customer'::character varying NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (user_id)
);

CREATE TABLE public.favorites (
	user_id int4 NOT NULL,
	vehicle_id int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT favorites_pkey PRIMARY KEY (user_id, vehicle_id),
	CONSTRAINT fk_favorites_user FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_favorites_vehicle FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(vehicle_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE public.vehicles (
	vehicle_id bigserial NOT NULL,
	type_id int8 NOT NULL,
	brand_id int8 NOT NULL,
	"name" varchar(150) NOT NULL,
	fuel_type varchar(50) NULL,
	transmission varchar(50) NULL,
	price numeric(12, 2) NOT NULL,
	stock int4 DEFAULT 0 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT vehicles_pkey PRIMARY KEY (vehicle_id),
	CONSTRAINT fk_vehicle_brand FOREIGN KEY (brand_id) REFERENCES public.brands(brand_id) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT fk_vehicle_type FOREIGN KEY (type_id) REFERENCES public.vehicle_types(type_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE public.vehicle_types (
	type_id bigserial NOT NULL,
	type_name varchar(100) NOT NULL,
	CONSTRAINT vehicle_types_name_unique UNIQUE (type_name),
	CONSTRAINT vehicle_types_pkey PRIMARY KEY (type_id)
);

CREATE TABLE public.brands (
	brand_id bigserial NOT NULL,
	brand_name varchar(100) NOT NULL,
	CONSTRAINT brands_name_unique UNIQUE (brand_name),
	CONSTRAINT brands_pkey PRIMARY KEY (brand_id)
);

CREATE TABLE transactions (
    transaction_id      BIGSERIAL PRIMARY KEY,
    order_id            VARCHAR(50) UNIQUE NOT NULL,
    user_id             BIGINT NOT NULL REFERENCES users(user_id),
    total_amount        NUMERIC(12,2) NOT NULL,
    payment_method      VARCHAR(50) NOT NULL,
    bank                VARCHAR(20),
    va_number           VARCHAR(100),
    status              VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transaction_items (
    detail_id       BIGSERIAL PRIMARY KEY,
    transaction_id  BIGINT NOT NULL REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    vehicle_id      BIGINT NOT NULL REFERENCES vehicles(vehicle_id),
    quantity        INT NOT NULL,
    price           NUMERIC(12,2) NOT NULL
);

CREATE TABLE payment_detail (
    payment_id              BIGSERIAL PRIMARY KEY,
    transaction_id          BIGINT NOT NULL REFERENCES transactions(transaction_id),
    midtrans_transaction_id TEXT,
    payment_type            VARCHAR(50),
    transaction_status      VARCHAR(50),
    fraud_status            VARCHAR(50),
    gross_amount            NUMERIC(12,2),
    paid_at                 TIMESTAMP NULL,
    note                    TEXT NULL
);

--  INDEXES 
-- vehicles: FK columns + low-stock query
CREATE INDEX idx_vehicles_brand_id ON vehicles(brand_id);
CREATE INDEX idx_vehicles_type_id ON vehicles(type_id);
CREATE INDEX idx_vehicles_stock ON vehicles(stock) WHERE stock <= 5;

-- favorites: query by vehicle_id (user_id sudah ter-cover oleh composite PK)
CREATE INDEX idx_favorites_vehicle_id ON favorites(vehicle_id);

-- transactions: FK + filter by status
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_status ON transactions(status);

-- transaction_items: FK columns
CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);
CREATE INDEX idx_transaction_items_vehicle_id ON transaction_items(vehicle_id);

-- payment_detail: FK column
CREATE INDEX idx_payment_detail_transaction_id ON payment_detail(transaction_id);

