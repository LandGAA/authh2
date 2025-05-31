CREATE TABLE users
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role   TEXT DEFAULT "user",
    create_at DATE NOT NULL
);

CREATE INDEX idx_users_email ON users (email);

insert into users (name, email, password, role, create_at)
values ('admin', 'admin@a.com', '$2a$10$Mvf/dtTA2U.jl31suPiQcuFHyRQszf4XHFvZe8w/HBcOa0TyysAsK', 'admin', '1')