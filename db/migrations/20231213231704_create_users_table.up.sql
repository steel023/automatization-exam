CREATE TABLE "users"
(
    "id"              uuid DEFAULT fn_uuid_time_ordered() PRIMARY KEY,
    "email"           varchar(255) NOT NULL UNIQUE CHECK (email <> ''),
    "password"        varchar(255) NOT NULL CHECK (password <> ''),
    "role"            int DEFAULT 0 NOT NULL CHECK (role BETWEEN 0 AND 3),
    "created_at"      timestamptz DEFAULT now()
);

CREATE TABLE "tokens"
(
    "id"          uuid DEFAULT fn_uuid_time_ordered() PRIMARY KEY,
    "token"       varchar(511) NOT NULL UNIQUE,
    "user_id"     uuid NOT NULL,
    "created_at"  timestamptz DEFAULT now(),
    "expires_at"  timestamptz DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- passwords: test123 hashed
INSERT INTO users (email, password, role) VALUES
                  ('admin@mail.ru', '$2a$10$jUOo/SnKN.kg2NgNFmZ7O.m2DPWmU9NczejYe3cfDL79ijvroum3q', 2),
                  ('moder@mail.ru', '$2a$10$jUOo/SnKN.kg2NgNFmZ7O.m2DPWmU9NczejYe3cfDL79ijvroum3q', 1),
                  ('user@mail.ru', '$2a$10$jUOo/SnKN.kg2NgNFmZ7O.m2DPWmU9NczejYe3cfDL79ijvroum3q', 0);