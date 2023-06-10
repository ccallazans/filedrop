CREATE TABLE IF NOT EXISTS users (
    UUID UUID PRIMARY KEY,
    name VARCHAR(50) not null,
    email VARCHAR(50) unique not null,
    password VARCHAR(100) not null,
    role VARCHAR(50) not null
);

CREATE TABLE IF NOT EXISTS files (
	UUID UUID primary key,
	filename VARCHAR(50) not null,
	size VARCHAR(50) not null,
	location_url VARCHAR(50) not null,
	user_uuid UUID REFERENCES users(uuid)
);

CREATE TABLE IF NOT EXISTS access_files (
    hash VARCHAR(6) PRIMARY KEY,
    lock BOOLEAN not null,
    access_code VARCHAR(100) not null,
    file_UUID UUID REFERENCES files(UUID)
);