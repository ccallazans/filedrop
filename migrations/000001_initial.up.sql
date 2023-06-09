CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY,
    name VARCHAR(50) not null,
    email VARCHAR(50) not null,
    password VARCHAR(50) not null,
    role VARCHAR(50) not null
);

CREATE TABLE IF NOT EXISTS files (
	id BIGINT primary key,
	filename VARCHAR(50) not null,
	size VARCHAR(50) not null,
	location_url VARCHAR(50) not null,
	uuser_uuid UUID REFERENCES users(uuid)
);