-- Active: 1741789120696@@127.0.0.1@5432@daisy_laondry@public
CREATE USER daisy_laondry WITH PASSWORD 'password';
-- Grant select privilege;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO daisy_laondry;
-- Grant all privileges;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO daisy_laondry;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS user_addresses;
DROP TABLE IF EXISTS user_otps;
DROP TABLE IF EXISTS staff_positions;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    password_hash TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    roles (id, name, description)
VALUES (
        gen_random_uuid (),
        'super_admin',
        'Pengelola sistem tingkat global'
    ),
    (
        gen_random_uuid (),
        'admin',
        'Administrator tenant laundry'
    ),
    (
        gen_random_uuid (),
        'outlet_manager',
        'Pengelola cabang laundry'
    ),
    (
        gen_random_uuid (),
        'cashier',
        'Kasir cabang outlet'
    ),
    (
        gen_random_uuid (),
        'courier',
        'Kurir antar-jemput laundry'
    ),
    (
        gen_random_uuid (),
        'laundry_staff',
        'Karyawan produksi laundry'
    ),
    (
        gen_random_uuid (),
        'customer',
        'Pelanggan laundry'
    );

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, role_id)
);

CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    label VARCHAR(100),
    address TEXT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    otp_code VARCHAR(10) NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE staff_positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    position VARCHAR(50) NOT NULL CHECK (
        position IN (
            'washing',
            'ironing',
            'folding',
            'packing'
        )
    ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, position)
);