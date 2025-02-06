CREATE TABLE `mercari`.`account` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `email` varchar(255) DEFAULT NULL,
    `buyer_id` varchar(255) DEFAULT NULL,
    `delivery_address` json DEFAULT NULL,
    `delivery_identifier` varchar(255) DEFAULT NULL,
    `access_token` varchar(255) DEFAULT NULL,
    `refresh_token` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_buyer_id` (`buyer_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci