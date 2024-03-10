CREATE TABLE IF NOT EXISTS transfer 
(
    id SERIAL PRIMARY KEY, 
    from_user INTEGER NOT NULL, 
    to_user INTEGER NOT NULL, 
    keyword VARCHAR(100), 
    address VARCHAR(100), 
    filename VARCHAR(100)
);


CREATE TABLE IF NOT EXISTS users 
(
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    friend_code VARCHAR(100), 
    mac_address VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS friends
(
    user_one INTEGER REFERENCES users(id), 
    user_two INTEGER REFERENCES users(id)
);