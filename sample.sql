CREATE TABLE wallets (
    address TEXT PRIMARY KEY,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL -- update this to use a secret to encrypt / decrypt the private key as needed
);

CREATE TABLE wallet_balances (
    address TEXT NOT NULL,
    coin_symbol VARCHAR(10) NOT NULL,
    amount BIGINT NOT NULL, -- amount parts
    CONSTRAINT fk_wallet FOREIGN KEY (address) REFERENCES wallets(address) ON DELETE CASCADE,
    PRIMARY_KEY (address, coin_symbol)
);

