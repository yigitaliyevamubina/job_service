-- Jobs Table
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID REFERENCES clients(id),
    price FLOAT NOT NULL DEFAULT 0.0,
    from_date VARCHAR(15) NOT NULL,
    to_date VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_jobs_owner_id ON jobs (owner_id);

-- Mock Data
INSERT INTO jobs (id, title, description, owner_id, from_date, to_date) VALUES
    ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Software Engineer', 'Developing new features for a web application', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '2024-01-12', '2024-12-12'),
    ('ac1f8087-cc21-4e62-b2f0-cb3c9a77b3e1', 'Data Analyst', 'Analyzing data trends and generating reports', 'ac1f8087-cc21-4e62-b2f0-cb3c9a77b3e1', '2024-01-12', '2024-12-12'),
    ('6af8b0c1-2e2c-4b5d-92d5-b0e8b29f65ad', 'Marketing Specialist', 'Creating marketing campaigns and strategies', '6af8b0c1-2e2c-4b5d-92d5-b0e8b29f65ad', '2024-01-12', '2024-12-12'),
    ('e963c3b1-37e8-4c5e-b15d-38c47d85e8b1', 'Graphic Designer', 'Designing logos, banners, and promotional materials', 'e963c3b1-37e8-4c5e-b15d-38c47d85e8b1', '2024-01-12', '2024-12-12'),
    ('da455e4a-7cb8-47bf-a0d2-18b74b5b5690', 'Project Manager', 'Overseeing project timelines and deliverables', 'da455e4a-7cb8-47bf-a0d2-18b74b5b5690', '2024-01-12', '2024-12-12'),
    ('47d23bf0-2050-4d8c-810e-d2b8d66fc7a0', 'Customer Support Representative', 'Providing assistance to customers via phone and email', '47d23bf0-2050-4d8c-810e-d2b8d66fc7a0', '2024-01-12', '2024-12-12'),
    ('d0f8c0ef-24d1-4f3d-ae6a-43b01ad6b587', 'Content Writer', 'Producing engaging content for blogs and social media', 'd0f8c0ef-24d1-4f3d-ae6a-43b01ad6b587', '2024-01-12', '2024-12-12'),
    ('9c4e8f6e-4d17-4854-91e7-36b63d8006aa', 'Sales Executive', 'Identifying and contacting potential clients', '9c4e8f6e-4d17-4854-91e7-36b63d8006aa', '2024-01-12', '2024-12-12'),
    ('ea95b085-702a-434d-823d-8d4059ef71e9', 'HR Coordinator', 'Managing recruitment processes and employee relations', 'ea95b085-702a-434d-823d-8d4059ef71e9', '2024-01-12', '2024-12-12'),
    ('b13db871-8e32-48dc-9255-f31205b964c0', 'Financial Analyst', 'Analyzing financial data and preparing forecasts', 'b13db871-8e32-48dc-9255-f31205b964c0', '2024-01-12', '2024-12-12');
