-- +migrate Up
CREATE TABLE allocated_key (
    key VARCHAR(10),
    allocated_at TIMESTAMP WITH TIME ZONE
);

-- +migrate Down
DROP TABLE allocated_key;