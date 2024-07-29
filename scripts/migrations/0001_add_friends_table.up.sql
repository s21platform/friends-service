CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    initiator UUID,
    user_id UUID
);