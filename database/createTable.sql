CREATE TABLE `user` (
    `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` VARCHAR(64) NULL,
    `password` VARCHAR(64) NULL,
    `business_name` VARCHAR(64) DEFAULT '',
    `created_at` DATE NULL,
    `updated_at` DATE NULL
);
