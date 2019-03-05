
-- Info about chat bots.
CREATE TABLE Bots (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    twitch_id BIGINT UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    oauth_token VARCHAR(255),
    is_enabled TINYINT(1) DEFAULT 1
);

-- Info about streamers.
CREATE TABLE Users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    twitch_id BIGINT UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    oauth_token VARCHAR(255)
);

-- Which streamer uses which bots
CREATE TABLE UserBots (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    bot_id INTEGER NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES Users(id),
    CONSTRAINT fk_bot_id FOREIGN KEY (bot_id) REFERENCES Bots(id),
    CONSTRAINT uk_user_id_bot_id UNIQUE (user_id, bot_id)
);

-- Command responses table.
CREATE TABLE Responses (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255),  -- optional response name, in case of alias.
    response TEXT NOT NULL
)

-- Commands table.
CREATE TABLE Commands (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL INDEX,
    user_id INTEGER NOT NULL,  -- foreign key to user id
    response_id INTEGER NOT NULL,  -- foreign key to response id, because of aliases.
    cooldown_second INTEGER,
    is_enabled TINYINT(1) DEFAULT 1,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES Users(id)
    CONSTRAINT fk_response_id FOREIGN KEY (response_id) REFERENCES Responses(id),
    CONSTRAINT uk_name_user_id UNIQUE (name, user_id)
);

