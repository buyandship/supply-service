CREATE TABLE `yahoo`.`transaction` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `bid_request_id` varchar(255) DEFAULT NULL,
    `price` bigint DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `error_message` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_bid_request_id` (`bid_request_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_updated_at` (`updated_at`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

