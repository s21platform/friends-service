CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    initiator UUID NOT NULL,
    user UUID NOT NULL
);