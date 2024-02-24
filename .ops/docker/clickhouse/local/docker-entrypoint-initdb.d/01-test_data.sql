CREATE TABLE IF NOT EXISTS test_data ON CLUSTER ch_cluster (
    `id` UUID DEFAULT generateUUIDv4() Codec(LZ4),
    sign Int8 Codec(DoubleDelta, LZ4),
    version UInt32 Codec(DoubleDelta, LZ4),
    `data` String Codec(LZ4),
    `created_at` DateTime64(9) DEFAULT NOW64() Codec(DoubleDelta, LZ4),
) Engine = ReplicatedVersionedCollapsingMergeTree('/clickhouse/tables/versioned/{shard}/test_data', '{replica}', sign, version)
PARTITION BY toYYYYMM(`created_at`)
PRIMARY KEY (`id`);
