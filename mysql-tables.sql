CREATE TABLE `t_source_file` (
     `payload_cid` varchar(100) NOT NULL,
     `file_name` varchar(50) DEFAULT NULL,
     `file_size` bigint DEFAULT NULL,
     `upload_id` bigint DEFAULT NULL,
     `status` varchar(30) DEFAULT NULL,
     `address` varchar(100) DEFAULT NULL,
     `create_at` datetime DEFAULT NULL,
     `index_status` varchar(30) DEFAULT NULL,
     PRIMARY KEY (`payload_cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `t_miner_deal` (
    `deal_cid` varchar(100) NOT NULL,
    `deal_id` bigint DEFAULT NULL,
    `miner_id` varchar(30) DEFAULT NULL,
    `status` varchar(50) DEFAULT NULL,
    `payload_cid` varchar(100) DEFAULT NULL,
    `piece_cid` varchar(100) DEFAULT NULL,
    `operate_status` varchar(30) DEFAULT NULL,
    PRIMARY KEY (`deal_cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `t_miner_peer` (
    `miner_id` varchar(30) DEFAULT NULL,
    `peer_id` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `t_failed_source` (
    `payload_cid` varchar(100) NOT NULL,
    PRIMARY KEY (`payload_cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `t_file_ipfs` (
    `data_cid` varchar(100) NOT NULL,
    `ipfs_url` varchar(256) NOT NULL,
    `ipfs_hash` varchar(100) DEFAULT NULL,
    PRIMARY KEY (`data_cid`,`ipfs_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

