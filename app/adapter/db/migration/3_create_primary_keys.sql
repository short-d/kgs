-- +migrate Up
ALTER TABLE available_key ADD PRIMARY KEY ("key");
ALTER TABLE allocated_key ADD PRIMARY KEY ("key");