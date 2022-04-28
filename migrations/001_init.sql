CREATE SCHEMA IF NOT EXISTS gophermart;

CREATE TABLE IF NOT EXISTS gophermart.users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    login varchar(255) NOT NULL UNIQUE,
    password varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS gophermart.orders(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    number varchar(55) NOT NULL UNIQUE,
    status varchar(55) NOT NULL,
    accrual INT,
    upload_date TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS gophermart.balances(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    current INT NOT NULL
);

CREATE TABLE IF NOT EXISTS gophermart.withdrawals(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES gophermart.orders(id),
    sum INT NOT NULL,
    date TIMESTAMP WITH TIME ZONE
);
