CREATE TABLE urls (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    url TEXT NOT NULL,
    short_url VARCHAR(15) NOT NULL UNIQUE,
    click_count INT UNSIGNED DEFAULT 0,
    created_at INT UNSIGNED DEFAULT 0,
    user_id VARCHAR(30) NOT NULL,
    INDEX (short_url, user_id)
);