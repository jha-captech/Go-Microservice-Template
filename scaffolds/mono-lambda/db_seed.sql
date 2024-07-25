-- Drop the users table if it already exists
DROP TABLE IF EXISTS users;

-- Create the users table
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(50)                                          NOT NULL,
    last_name  VARCHAR(50)                                          NOT NULL,
    role       VARCHAR(10) CHECK (role IN ('Customer', 'Employee')) NOT NULL,
    user_id    INTEGER UNIQUE                                       NOT NULL
);

-- Insert 10 records into the users table
INSERT INTO users (first_name, last_name, role, user_id)
VALUES ('John', 'Doe', 'Customer', 1001),
       ('Jane', 'Smith', 'Employee', 1002),
       ('Robert', 'Johnson', 'Employee', 1003),
       ('Emily', 'Davis', 'Customer', 1004),
       ('Michael', 'Brown', 'Employee', 1005),
       ('Linda', 'Wilson', 'Employee', 1006),
       ('David', 'Martinez', 'Customer', 1007),
       ('Elizabeth', 'Taylor', 'Employee', 1008),
       ('Richard', 'Anderson', 'Employee', 1009),
       ('Susan', 'Thomas', 'Customer', 1010);

-- Select all records to verify the insertion
SELECT *
FROM users;
