ALTER TABLE mercari.token ADD account_id INT UNSIGNED DEFAULT 1 NOT NULL;

ALTER TABLE mercari.review ADD account_id INT UNSIGNED DEFAULT 1 NOT NULL;

ALTER TABLE mercari.transaction ADD account_id INT UNSIGNED DEFAULT 1 NOT NULL;

ALTER TABLE mercari.message ADD account_id INT UNSIGNED DEFAULT 1 NOT NULL;


ALTER TABLE mercari.account ADD status VARCHAR(100) NULL;
ALTER TABLE mercari.account ADD priority INT UNSIGNED NULL;
ALTER TABLE mercari.account ADD banned_at datetime(3) NULL;
ALTER TABLE mercari.account ADD active_at datetime(3) NULL;

ALTER TABLE mercari.account ADD CONSTRAINT account_unique UNIQUE KEY (email);
