CREATE TABLE IF NOT EXISTS gophermart.orders(
    number varchar(255) PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    status VARCHAR(55) NOT NULL,
    accrual INT,
    upload_date TIMESTAMP WITH TIME ZONE,
    UNIQUE (number, user_id)
);