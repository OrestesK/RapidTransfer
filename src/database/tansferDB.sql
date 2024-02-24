DROP TABLE IF EXISTS user CASCADE;
DROP TABLE IF EXISTS friends CASCADE;
DROP TABLE IF EXISTS transfer CASCADE;


CREATE TABLE user(
    userID int,
    name VARCHAR(100),
    key int
)

CREATE TABLE friends(
    idFrom int,
    idTo, int
)

CREATE TABLE transfer(
    uidFrom int,
    uidTo int
)