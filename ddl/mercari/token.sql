CREATE TABLE `mercari`.`token` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `access_token` varchar(255) DEFAULT NULL,
    `refresh_token` LONGTEXT DEFAULT NULL,
    `expires_in` int DEFAULT NULL,
    `token_type` varchar(255) DEFAULT NULL,
    `scope` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci