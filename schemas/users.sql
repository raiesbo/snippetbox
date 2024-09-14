-- Select the database (PostgreSQL does not use the 'USE' statement; instead, you should connect to the desired database)
-- \c snippetbox;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Adding a unique constraint on the email field
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
