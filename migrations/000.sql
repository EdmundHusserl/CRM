CREATE DATABASE customers WITH 
  OWNER = postgres
  ENCODING = 'UTF8'
  LC_COLLATE = 'en_US.UTF-8'
  LC_CTYPE = 'en_US.UTF-8'
  CONNECTION LIMIT = -1;

\c customers;

-- Ensure the uuid-ossp extension is enabled for UUID generation (if needed)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(50),
    contacted BOOLEAN DEFAULT FALSE
);

INSERT INTO customers (name, role, email, phone_number, contacted)
VALUES 
('Wladston', 1, 'wladston@corp.com', '514-555-6666', false),
('Willy', 1, 'wily@corp.com', '514-555-777', false),
('Wanton', 1, 'wanton@corp.com', '514-555-8888', false);