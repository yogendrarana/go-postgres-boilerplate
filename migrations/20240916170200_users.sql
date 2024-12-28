-- +goose Up
-- create enum type first
CREATE TYPE user_role AS ENUM ('super_admin', 'admin', 'user');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    bio TEXT,
    avatar_url TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
DROP TYPE user_role;