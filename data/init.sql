set time zone 'UTC';
create extension pgcrypto;

CREATE TABLE resources (
    id VARCHAR (255) PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    type VARCHAR (255) NOT NULL,
    available BOOLEAN
);

CREATE TABLE users (
    id VARCHAR (255) PRIMARY KEY, 
    username VARCHAR (255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    deleted_at TIMESTAMP
);

CREATE TABLE lambdas (
    id VARCHAR (255) PRIMARY KEY,
    user_id VARCHAR (255),
    name VARCHAR (255),
    concurrent_limit INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE virtual_machines (
    id VARCHAR (255) PRIMARY KEY,
    user_id VARCHAR (255),
    name VARCHAR (255),
    cpus INT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE sql_databases (
    id VARCHAR (255) PRIMARY KEY,
    user_id VARCHAR (255),
    name VARCHAR (255),
    username VARCHAR (255) NOT NULL,
    password VARCHAR (255) NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE nosql_databases (
    id VARCHAR (255) PRIMARY KEY,
    user_id VARCHAR (255),
    name VARCHAR (255),
    shards INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

INSERT INTO resources (id, name, type, available) VALUES ('resource-001', 'Serverless Lambda Function', 'lambda', TRUE);
INSERT INTO resources (id, name, type, available) VALUES ('resource-002', 'Virtual Machine', 'virtual_machine', TRUE);
INSERT INTO resources (id, name, type, available) VALUES ('resource-003', 'SQL Database', 'sql_database', TRUE);
INSERT INTO resources (id, name, type, available) VALUES ('resource-004', 'No-SQL Database', 'nosql_database', TRUE);

INSERT INTO users (id, username, password, created_at, updated_at) VALUES ('demo-user-001', 'password', CURRENT_DATE, CURRENT_DATE);

INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-001', 'demo-user-001', 'My preset lambda 1', 1, CURRENT_DATE, CURRENT_DATE);
INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-002', 'demo-user-001', 'My preset lambda 2', 10, CURRENT_DATE, CURRENT_DATE);
INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-003', 'demo-user-001', 'My preset lambda 3', 25, CURRENT_DATE, CURRENT_DATE);
INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-004', 'demo-user-001', 'My preset lambda 4', 43, CURRENT_DATE, CURRENT_DATE);
INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-005', 'demo-user-001', 'My preset lambda 5', 50, CURRENT_DATE, CURRENT_DATE);
INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-006', 'demo-user-001', 'My preset lambda 6', 76, CURRENT_DATE, CURRENT_DATE);
INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at) VALUES ('preset-lambda-007', 'demo-user-001', 'My preset lambda 7', 100, CURRENT_DATE, CURRENT_DATE);

INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-001', 'demo-user-001', 'My preset virtual machine 1', 1, 1, CURRENT_DATE, CURRENT_DATE);
INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-002', 'demo-user-001', 'My preset virtual machine 2', 4, 2, CURRENT_DATE, CURRENT_DATE);
INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-003', 'demo-user-001', 'My preset virtual machine 3', 2, 3, CURRENT_DATE, CURRENT_DATE);
INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-004', 'demo-user-001', 'My preset virtual machine 4', 1, 4, CURRENT_DATE, CURRENT_DATE);
INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-005', 'demo-user-001', 'My preset virtual machine 5', 2, 1, CURRENT_DATE, CURRENT_DATE);
INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-006', 'demo-user-001', 'My preset virtual machine 6', 8, 2, CURRENT_DATE, CURRENT_DATE);
INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at) VALUES ('preset-vm-007', 'demo-user-001', 'My preset virtual machine 7', 2, 3, CURRENT_DATE, CURRENT_DATE);

INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-001', 'demo-user-001', 'My preset sql database 1', 'db-admin', 'dbpassword1', 1, CURRENT_DATE, CURRENT_DATE);
INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-002', 'demo-user-001', 'My preset sql database 2', 'db-admin', 'dbpassword1', 3, CURRENT_DATE, CURRENT_DATE);
INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-003', 'demo-user-001', 'My preset sql database 3', 'db-admin', 'dbpassword1', 2, CURRENT_DATE, CURRENT_DATE);
INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-004', 'demo-user-001', 'My preset sql database 4', 'db-admin', 'dbpassword1', 2, CURRENT_DATE, CURRENT_DATE);
INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-005', 'demo-user-001', 'My preset sql database 5', 'db-admin', 'dbpassword1', 2, CURRENT_DATE, CURRENT_DATE);
INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-006', 'demo-user-001', 'My preset sql database 6', 'db-admin', 'dbpassword1', 3, CURRENT_DATE, CURRENT_DATE);
INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at) VALUES ('preset-sql-db-007', 'demo-user-001', 'My preset sql database 7', 'db-admin', 'dbpassword1', 2, CURRENT_DATE, CURRENT_DATE);

INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-001', 'demo-user-001', 'My preset No SQL database 1', 10, CURRENT_DATE, CURRENT_DATE);
INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-002', 'demo-user-001', 'My preset No SQL database 2', 20, CURRENT_DATE, CURRENT_DATE);
INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-003', 'demo-user-001', 'My preset No SQL database 3', 30, CURRENT_DATE, CURRENT_DATE);
INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-004', 'demo-user-001', 'My preset No SQL database 4', 40, CURRENT_DATE, CURRENT_DATE);
INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-005', 'demo-user-001', 'My preset No SQL database 5', 10, CURRENT_DATE, CURRENT_DATE);
INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-006', 'demo-user-001', 'My preset No SQL database 6', 20, CURRENT_DATE, CURRENT_DATE);
INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at) VALUES ('preset-nosql-db-007', 'demo-user-001', 'My preset No SQL database 7', 30, CURRENT_DATE, CURRENT_DATE);
