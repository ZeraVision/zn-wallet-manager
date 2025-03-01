# zn-wallet-manager
ZERA Network (zn) Wallet Manager is an open-source repository designed to facilitate the management of incoming and outgoing CoinTXN / MintTXN transactions on the ZERA Network. It supports both uniquely generated wallets and memo-based (incoming only) configurations, making it a versatile solution for developers building financial applications or blockchain-based platforms that may require managable and scalable deposits and withdrawls.

## Note About Address / Memo Deposists
Generally, address based is reccomended as it limits errors and complexities for the end user.

For **deposits only**, memo support is added. If you choose to use memos, it is integral that **you check** to make sure it is to an authorized address. In theory this means that a deposit webhook could be sent and if you do not check the to address being one of your authorized addresses you may think you have a deposit when you don't.

In most cases, if the memo is present, you will receive a webhook (covered in more detail below).

There are specific circumstances (for security reasons) where the webhook will not be sent to you:
1. If it is a withdrawl
2. If the transaction is a x-to-many transaction and the memo is mistakenly included in the base memo and not the output memo. (Generally unlikely to occur)

In the unlikely event of user error causing 2, manual review will be required by your support team.

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

## Intended Design
This application is intended to satisfy most use cases but requires developers to do "last-mile" customizability for their spefic needs and use cases.

This application takes you to the final functions of the webhook, by default parameterizing deposit/withdraw functions with a simple to use object.

Developers will want to consider:
1. Integrating their KMS / secure storage system into the application.
2. Specific deposit / withdrawl functionality customized to your needs (ie email, notifications, database updates)

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

