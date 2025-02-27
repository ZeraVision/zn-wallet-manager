CREATE TABLE wallets (
    address TEXT PRIMARY KEY, -- TODO secure according to your requirements (see InsertWallet func)
    public_key TEXT NOT NULL, -- TODO secure according to your requirements (see InsertWallet func)
    private_key TEXT NOT NULL, -- TODO secure according to your requirements (see InsertWallet func)
    kms_id TEXT NOT NULL, -- TODO if not using kms / ext encryption/decryption, remove this field & update as nessasary
    empty BOOLEAN NOT NULL DEFAULT TRUE,
);