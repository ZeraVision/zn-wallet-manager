CREATE TABLE wallets (
    address TEXT PRIMARY KEY, -- TODO secure according to your requirements (see InsertWallet func)
    public_key TEXT NOT NULL, -- TODO secure according to your requirements (see InsertWallet func)
    private_key TEXT NOT NULL, -- TODO secure according to your requirements (see InsertWallet func)
    kms_id TEXT,
    empty BOOLEAN NOT NULL DEFAULT TRUE,
);