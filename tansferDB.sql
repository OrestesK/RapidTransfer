DROP TABLE IF EXISTS user CASCADE;
DROP TABLE IF EXISTS friends CASCADE;
DROP TABLE IF EXISTS transfer CASCADE;


CREATE TABLE user (
    id SERIAL INT PRIMARY KEY,
    name VARCHAR(100),
    friendCode VARCHAR(16) UNIQUE
)

CREATE TABLE friends(
    id SERIAL INT PRIMARY KEY,
    user_from INT REFERENCES user,
    user_to INT REFERENCES user,
)

CREATE TABLE transfer(
    uidFrom int,
    uidTo int,
    key VARCHAR(100)
)