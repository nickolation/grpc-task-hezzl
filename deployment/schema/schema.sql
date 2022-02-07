CREATE TABLE clicklogs (
    postgres_id Int32 Codec(DoubleDelta, LZ4),
    time DateTime Codec(DoubleDelta, LZ4),
    start_date ALIAS toDate(time),
    age Int32,
    username String,
    gender String,
    description String
) Engine = MergeTree
PARTITION BY toYYYYMM(time)
ORDER BY (postgres_id, time); 

CREATE TABLE clicklogs_queue (
    postgres_id Int32,
    time DateTime Codec(DoubleDelta, LZ4),
    age Int32,
    username String,
    gender String,
    description String
)
ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka-1:19092, kafka-2:29092, kafka-3:39092',
       kafka_topic_list = 'logs',
       kafka_group_name = 'readings_consumer_group1',
       kafka_format = 'JSONEachRow',
       kafka_max_block_size = 1048576; 

CREATE MATERIALIZED VIEW clicklogs_queue_mv TO clicklogs AS
SELECT postgres_id, time, username, age, gender, description
FROM clicklogs_queue;  

SHOW
TABLES();

SELECT *
FROM clicklogs;