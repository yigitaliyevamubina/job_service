-- Client Table
CREATE TABLE clients (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Mock Data
INSERT INTO clients (id, username, email, phone, address) VALUES
    ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'John Doe', 'john@example.com', '+1234567890', '123 Main Street, Anytown, USA'),
    ('ac1f8087-cc21-4e62-b2f0-cb3c9a77b3e1', 'Jane Smith', 'jane@example.com', '+0987654321', '456 Elm Street, Othertown, USA'),
    ('6af8b0c1-2e2c-4b5d-92d5-b0e8b29f65ad', 'Michael Johnson', 'michael@example.com', '+1122334455', '789 Oak Street, Somewhere, USA'),
    ('e963c3b1-37e8-4c5e-b15d-38c47d85e8b1', 'Emily Brown', 'emily@example.com', '+9988776655', '101 Pine Street, Nowhere, USA'),
    ('da455e4a-7cb8-47bf-a0d2-18b74b5b5690', 'Christopher Lee', 'chris@example.com', '+3344556677', '555 Maple Street, Anywhere, USA'),
    ('47d23bf0-2050-4d8c-810e-d2b8d66fc7a0', 'Amanda Wilson', 'amanda@example.com', '+6677889900', '777 Cedar Street, Elsewhere, USA'),
    ('d0f8c0ef-24d1-4f3d-ae6a-43b01ad6b587', 'David Martinez', 'david@example.com', '+5566778899', '888 Birch Street, Nowheretown, USA'),
    ('9c4e8f6e-4d17-4854-91e7-36b63d8006aa', 'Jessica Taylor', 'jessica@example.com', '+2233445566', '999 Walnut Street, Hometown, USA'),
    ('ea95b085-702a-434d-823d-8d4059ef71e9', 'Kevin Harris', 'kevin@example.com', '+9988776655', '111 Oak Street, Anotherplace, USA'),
    ('b13db871-8e32-48dc-9255-f31205b964c0', 'Michelle Clark', 'michelle@example.com', '+1122334455', '222 Elm Street, Yetanotherplace, USA');
