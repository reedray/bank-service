CREATE TABLE IF NOT EXISTS Customer
(
    id        uuid NOT NULL PRIMARY KEY,
    username  text NOT NULL,
    password  text NOT NULL,
    role      text,
    createdAt timestamp without time zone,
    deletedAt timestamp without time zone,
    status    text,
    balance   json
);


CREATE TABLE IF NOT EXISTS Loan
(
    "id"         uuid NOT NULL PRIMARY KEY,
    amount       numeric,
    currencyCode text,
    createdAt    timestamp without time zone,
    expiresAt    timestamp with time zone,
    customerID   uuid references Customer (id)
);

CREATE TABLE IF NOT EXISTS TransactionType
(
    id   integer NOT NULL PRIMARY KEY,
    name text
);


CREATE TABLE IF NOT EXISTS Transaction
(
    transactionID     uuid NOT NULL PRIMARY KEY,
    fromAccountID     uuid,
    toAccountID       uuid,
    amount            numeric,
    currencyCode      text,
    transactionTypeID integer REFERENCES TransactionType (id)
);


INSERT INTO TransactionType(id, name)
VALUES (1, 'withdraw'),
       (2, 'deposit'),
       (3, 'transfer'),
       (4, 'balance');