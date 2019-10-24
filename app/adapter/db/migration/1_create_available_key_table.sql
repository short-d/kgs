-- +migrate Up
CREATE TABLE available_key (
    key VARCHAR(10),
    created_at TIMESTAMP WITH TIME ZONE
);

-- +migrate Down
DROP TABLE available_key;