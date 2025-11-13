CREATE TABLE `yahoo`.`bid_request` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `request_type` varchar(255) DEFAULT NULL,
    `order_id` varchar(255) DEFAULT NULL,
    `auction_id` varchar(255) DEFAULT NULL,
    `max_bid` bigint DEFAULT NULL,
    `quantity` int DEFAULT NULL,
    `partial` tinyint(1) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `error_message` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_auction_id` (`auction_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_updated_at` (`updated_at`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

