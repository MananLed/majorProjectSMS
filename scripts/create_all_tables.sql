CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    middle_name TEXT,
    last_name TEXT NOT NULL,
    mobile_number TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'officer', 'resident')) DEFAULT 'resident'
);

CREATE TABLE IF NOT EXISTS service_requests (
    request_id UUID PRIMARY KEY,
    resident_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'inProgress', 'completed', 'cancelled')) DEFAULT 'pending',
    time_slot TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    service_type TEXT NOT NULL CHECK (service_type IN ('electrician', 'plumber'))
);

CREATE TABLE IF NOT EXISTS notices (
    id UUID PRIMARY KEY,
    date_issued TIMESTAMP NOT NULL,
    content TEXT NOT NULL,
    month INT NOT NULL CHECK (month >= 1 AND month <= 12),
    year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY,
    amount NUMERIC(10, 2) NOT NULL,
    month INT NOT NULL CHECK (month >= 1 AND month <= 12),
    year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS feedbacks (
    id UUID PRIMARY KEY,
    resident_id TEXT NOT NULL REFERENCES users(id),
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    content TEXT
);

CREATE INDEX IF NOT EXISTS idx_service_requests_resident_id
    ON service_requests (resident_id);

CREATE INDEX IF NOT EXISTS idx_feedbacks_resident_id
    ON feedbacks (resident_id);

