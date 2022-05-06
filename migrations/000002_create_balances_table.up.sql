CREATE TABLE IF NOT EXISTS gophermart.balances(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES gophermart.users(id),
    current INT NOT NULL DEFAULT 0
);