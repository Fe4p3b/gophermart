CREATE TABLE IF NOT EXISTS gophermart.users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    login varchar(255) NOT NULL UNIQUE,
    password varchar(100) NOT NULL
);