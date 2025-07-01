-- Seed script to create initial admin user
-- Run this script directly in your PostgreSQL database

-- First, make sure the role column exists
ALTER TABLE staff ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user';

-- Insert initial admin user
-- Password: admin123 (hashed with bcrypt)
INSERT INTO staff (
    first_name, 
    last_name, 
    address_id, 
    email, 
    store_id, 
    active, 
    username, 
    password, 
    role, 
    last_update
) VALUES (
    'System',
    'Administrator',
    1,
    'admin@filmrental.com',
    'store1',
    true,
    'admin',
    '$2a$10$yGKnYPElWUlAEi5NxSX/De.OH3yTDII1EghnHK1VCjyNz4sH71UIG', -- admin123
    'admin',
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- Verify the admin was created
SELECT staff_id, username, role, active FROM staff WHERE username = 'admin'; 