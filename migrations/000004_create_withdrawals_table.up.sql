CREATE TABLE IF NOT EXISTS gophermart.withdrawals(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(255) NOT NULL,
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    sum INT NOT NULL,
    date TIMESTAMP WITH TIME ZONE
);
