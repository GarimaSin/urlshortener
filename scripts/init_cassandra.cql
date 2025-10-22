CREATE KEYSPACE IF NOT EXISTS urlshort WITH replication = {'class':'SimpleStrategy','replication_factor':3};
USE urlshort;
CREATE TABLE IF NOT EXISTS urls (short_code text PRIMARY KEY, long_url text, created_at timestamp);
