CREATE DATABASE IF NOT EXISTS explainer;

CREATE TABLE IF NOT EXISTS explainer.explanations
(
    incident_id String,
    user_id String,
    decision LowCardinality(String),
    score Float64,
    features_json String,
    created_at DateTime('UTC'),
    ingested_at DateTime('UTC') DEFAULT now()
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (incident_id, created_at)
TTL created_at + INTERVAL 1 MONTH
SETTINGS index_granularity = 8192;
