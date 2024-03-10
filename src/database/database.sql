CREATE TABLE IF NOT EXISTS transfer 
(
    id SERIAL PRIMARY KEY, 
    userFrom INT NOT NULL, 
    userTo INT NOT NULL, 
    keyword VARCHAR(100), 
    address VARCHAR(100), 
    filename VARCHAR(100)
);


CREATE TABLE IF NOT EXISTS users 
(
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    keyword VARCHAR(100), 
    macaddr VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS friends
(
    orig_user INT NOT NULL, 
    friend_id INT NOT NULL, 
);