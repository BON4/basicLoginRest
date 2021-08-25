CREATE TYPE role as enum ('admin', 'user', 'viewer');

DROP TABLE IF EXISTS users;
CREATE TABLE users (id SERIAL, username varchar NOT NULL UNIQUE, email varchar NOT NULL UNIQUE, role role not null, password bytea);
SELECT * FROM users offset 1 limit 10;
UPDATE users set username = 'glad' where id = 2 returning *;

SELECT * FROM users WHERE username LIKE 'v%' AND email LIKE '%gmail.com';


DROP TABLE IF EXISTS usersTest;
CREATE TABLE usersTest (id SERIAL, username varchar NOT NULL UNIQUE, email varchar NOT NULL UNIQUE, role role not null, password bytea);
SELECT * FROM usersTest;
DELETE FROM usersTest where id = 9;


INSERT INTO users (username, email, role, password) values ('vlad', 'vlad@gmail.com', 'admin',decode('f646f00b070d2d12ab45d0bc119217f22a3280299db7c2283d3815ed18bef3e8', 'hex'));

INSERT INTO users (username, email, role, password) values ('abcd', 'bcda@gmail.com', 'admin',decode('f646f00b070d2d12ab45d0bc119217f22a3280299db7c2283d3815ed18bef3e8', 'hex'));
INSERT INTO users (username, email, role, password) values ('bcda', 'abcd@gmail.com', 'admin',decode('f646f00b070d2d12ab45d0bc119217f22a3280299db7c2283d3815ed18bef3e8', 'hex'));
select count(*) from users where (username = 'abcd' or email = 'abcd@gmail.com');

select EXISTS(select * from users where (username = 'abc' or email = 'abc@gmail.com'));