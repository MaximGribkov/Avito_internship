CREATE TABLE users (
    user_id INT PRIMARY KEY NOT NULL,
    balance DECIMAL NOT NULL,
    reserve DECIMAL NOT NULL,
    CONSTRAINT non_negative_balance CHECK (balance >= 0 :: DECIMAL)
);

CREATE TABLE services (
    services_id INT PRIMARY KEY NOT NULL,
    name_ser VARCHAR(25) NOT NULL
);

CREATE TABLE orders (
    order_id INTEGER PRIMARY KEY NOT NULL,
    user_id INTEGER REFERENCES users(user_id) NOT NULL,
    services_id INTEGER REFERENCES services(services_id) NOT NULL,
    price DECIMAL NOT NULL,
    time_operation TIMESTAMP WITH TIME ZONE NOT NULL,
    status_order VARCHAR(25) NOT NULL
);

CREATE TABLE operation_history (
    history_id SERIAL PRIMARY KEY NOT NULL,
    amount DECIMAL NOT NULL,
    user_id INTEGER REFERENCES users(user_id),
    time_history TIMESTAMP NOT NULL
);

INSERT INTO services(services_id, name_ser) VALUES
    ('1', 'Услуга №1'),
    ('2', 'Услуга №2'),
    ('3', 'Услуга №3');

INSERT INTO users (user_id, balance, reserve) VALUES (
    '1',
    '0',
    '0'
);