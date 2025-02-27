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

## Configuration
To configure the application, provide the following environment variables:

```plaintext
GRPC_ADDR=
API_KEY=
```

This application is set up to use your API Key directly. If you'd prefer to modify it to use a bearer, you can view the [ZV Bearer Issuer](https://github.com/ZeraVision/zv-bearer-issuer) sample program.

## Obtaining an API Key
To request an API key, please contact us at [Zera Vision](https://www.zera.vision/contact).

A developer platform with automated API key issuance and native GUI analytics is planned for future release. Stay updated at [ZV Explorer](https://explorer.zera.vision/apis) *(coming soon)*.

## Docker Deployment
Coming soon

