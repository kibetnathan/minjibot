-- +goose Up
-- Per-guild birthday configuration: the channel where automated birthday
-- notices are posted and the temporary role awarded on a member's birthday.
-- (`birthday channel` and `birthday role`.)
CREATE TABLE guild_birthday_settings (
    guild_id   VARCHAR(20) PRIMARY KEY REFERENCES guilds(id) ON DELETE CASCADE,
    channel_id VARCHAR(20),
    role_id    VARCHAR(20),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Per-member birthdays. One row per guild+member so the same user can have
-- different birthdays (or opt out) per server. The month/day drive the
-- automated celebration and the upcoming list; the year is optional.
-- (`birthday add`, `birthday list`, `birthday celebrate`.)
CREATE TABLE birthdays (
    id         BIGSERIAL PRIMARY KEY,
    guild_id   VARCHAR(20) NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    user_id    VARCHAR(20) NOT NULL,
    birthday   DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_guild_user_birthday UNIQUE (guild_id, user_id)
);

-- Used to quickly find members whose birthday falls on a given month/day and
-- to produce the upcoming-birthdays list sorted by date.
CREATE INDEX idx_birthdays_upcoming ON birthdays (guild_id, birthday);

-- +goose Down
DROP TABLE IF EXISTS birthdays;
DROP TABLE IF EXISTS guild_birthday_settings;
