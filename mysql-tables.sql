CREATE TABLE `t_source_file` (
    `data_cid` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
    `file_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `file_size` bigint DEFAULT NULL,
    `rebuild_flag` tinyint(1) DEFAULT NULL,
    `rebuild_status` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `create_at` datetime DEFAULT NULL,
    PRIMARY KEY (`data_cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `t_file_miner` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `data_cid` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `miner_id` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `status` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_t_file_miner_deleted_at` (`deleted_at`),
    KEY `fk_t_source_file_miner_ids` (`data_cid`),
    CONSTRAINT `fk_t_source_file_miner_ids` FOREIGN KEY (`data_cid`) REFERENCES `t_source_file` (`data_cid`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `t_miner_peer` (
    `miner_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `peer_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`miner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `t_miner` (
    `miner_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    PRIMARY KEY (`miner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `t_file_ipfs` (
    `data_cid` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
    `ipfs_url` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`data_cid`),
    CONSTRAINT `fk_t_source_file_ipfs_urls` FOREIGN KEY (`data_cid`) REFERENCES `t_source_file` (`data_cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

