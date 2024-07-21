CREATE TABLE user_invite (
    id SERIAL PRIMARY KEY,
    initiator UUID NOT NULL,
    invited VARCHAR(255) NOT NULL,
    is_closed BOOLEAN DEFAULT FALSE,
    update_time TIMESTAMP
);