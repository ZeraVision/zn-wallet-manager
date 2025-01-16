# zn-wallet-manager

ZERA Network (zn) Wallet Manager is an open-source repository designed to facilitate the management of incoming and outgoing CoinTXN transactions on the ZERA Network. It supports both uniquely generated wallets and memo-based configurations, making it a versatile solution for developers building financial applications or blockchain-based platforms that may require managable and scalable deposits and withdrawls.

## Use Cases
There are various use cases that this repository can help satisfy.

This repo can potentially be useful for:
1. Wallet or Memo based deposit & withdrawl systems interacting with a simple PostgreSQL compatible database.
2. Wallet creation and key management (based on BIP39 or random entropy)
3. Varying levels of transfer functionality, including one-to-one, one-to-many, many-to-one, and many-to-many.

## Dependencies
Most of the functionality of this program is intended to work simply with go.

The following gobinding(s) are required in your environment:
- [Libsodium](github.com/GoKillers/libsodium-go/cryptosign)

To set up the environment, you can either use Docker (as described below) or install dependencies locally.

## Docker Deployment

