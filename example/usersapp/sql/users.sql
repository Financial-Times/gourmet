CREATE TABLE users
(
    uid          integer PRIMARY KEY AUTOINCREMENT,
    first_name   text NOT NULL,
    last_name    text NOT NULL,
    email        text NOT NULL,
    phone        text,
    created      integer,
    last_updated integer
);
