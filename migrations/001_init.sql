CREATE SCHEMA IF NOT EXISTS gophermart;

CREATE TABLE IF NOT EXISTS gophermart.users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    login varchar(255) NOT NULL UNIQUE,
    password varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS gophermart.orders(
    number varchar(255) PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    status INT NOT NULL,
    accrual INT,
    upload_date TIMESTAMP WITH TIME ZONE,
    UNIQUE (number, user_id)
);

CREATE TABLE IF NOT EXISTS gophermart.balances(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    current INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS gophermart.withdrawals(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(255) NOT NULL REFERENCES gophermart.orders(number),
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    sum INT NOT NULL,
    date TIMESTAMP WITH TIME ZONE
);
