CREATE TABLE `mercari`.`message` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `trx_id` varchar(255) DEFAULT NULL,
    `message` LONGTEXT DEFAULT NULL,
    `buyer_id` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_trx_id` (`trx_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci