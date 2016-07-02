CREATE TABLE `users` (
    `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` VARCHAR(64) NULL,
    `password` VARCHAR(64) NULL,
    `status` BOOLEAN
);

CREATE TABLE `punches` (
    `pid` INTEGER PRIMARY KEY AUTOINCREMENT,
    `uid` INTEGER NULL,
    `in` TIMESTAMP NULL,
    `out` TIMESTAMP NULL,
    `words` INTEGER NULL,
    `description` TEXT NULL,
    FOREIGN KEY(uid) REFERENCES users(uid)
);
