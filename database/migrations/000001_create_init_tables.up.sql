-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="Asia/Jakarta";

-- Create users table
CREATE TABLE users (
    id VARCHAR (36) DEFAULT uuid_generate_v4 () PRIMARY KEY,
    username VARCHAR (255) NOT NULL UNIQUE,
    email VARCHAR (255) NOT NULL UNIQUE,
    password VARCHAR (255) NOT NULL,
    role_id VARCHAR (36) NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    created_by VARCHAR (100),
    updated_at TIMESTAMP NULL,
    updated_by VARCHAR (100),
    is_deleted boolean DEFAULT false
);

-- Create books authors
CREATE TABLE authors (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    created_by VARCHAR (100),
    updated_at TIMESTAMP NULL,
    updated_by VARCHAR (100),
    is_deleted boolean DEFAULT false
);

INSERT INTO users(
	id, username, email, password, role_id, status, created_at, updated_at)
	VALUES (uuid_generate_v4(), 'admin', 'admin@gmail.com', '$2a$10$3lQxgep/NIdg.ibH5Ydeo.yt3P6MSLATdAbvgoiz37KGGchVMY6.G', 'admin', 1, now(), now());
