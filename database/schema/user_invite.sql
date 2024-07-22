-- CREATE TABLE user_invite (
--     id SERIAL PRIMARY KEY,
--     initiator UUID NOT NULL,
--     invited VARCHAR(255) NOT NULL,
--     is_closed BOOLEAN DEFAULT FALSE,
--     update_time TIMESTAMP
-- );

INSERT INTO user_invite (id, initiator, invited, is_closed, update_time)
VALUES (12346, gen_random_uuid(), 'a@a.ru', false, '2024-07-22 10:15:30');
--
select * from user_invite;