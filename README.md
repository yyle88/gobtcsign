[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/gobtcsign/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/gobtcsign/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/gobtcsign)](https://pkg.go.dev/github.com/yyle88/gobtcsign)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/gobtcsign/master.svg)](https://coveralls.io/github/yyle88/gobtcsign?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.22%2C%201.23-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yyle88/gobtcsign.svg)](https://github.com/yyle88/gobtcsign/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/gobtcsign)](https://goreportcard.com/report/github.com/yyle88/gobtcsign)

# gobtcsign

`gobtcsign` is a concise and efficient Bitcoin transaction signing library designed to help developers quickly build, sign, and verify Bitcoin transactions.

`gobtcsign` is a Golang package that simplifies BTC/DOGECOIN transaction signing and serves as a gateway for developers to explore BTC blockchain knowledge.

---

## CHINESE README

[中文说明](README.zh.md)

---

## Installation

```bash
go get github.com/yyle88/gobtcsign
```

---

## Features

Here are the core features provided by `gobtcsign`:

1. **Transaction Construction**: Efficiently construct transactions with support for multiple inputs and outputs, including automatic change calculation. A dynamic fee adjustment feature allows users to control transaction costs.
2. **Transaction Size Estimation**: Estimate the virtual size (vSize) of transactions based on the number and type of inputs/outputs. This helps developers set appropriate fee rates based on real-time conditions.
3. **Transaction Signing**: Compatible with multiple address types, including P2PKH, P2SH, and SegWit. Developers can use private keys to sign transaction inputs.
4. **Signature Verification**: Ensure transaction signatures are valid, reducing the risk of rejection by the network due to signature issues.
5. **Transaction Serialization**: Serialize signed transactions into hexadecimal strings for direct broadcasting to the Bitcoin network.

---

## Dependencies

`gobtcsign` relies on the following key modules:

- **github.com/btcsuite/btcd**: Implements Bitcoin's core protocol and serves as the foundation for building and parsing transactions.
- **github.com/btcsuite/btcd/btcec/v2**: Handles elliptic curve cryptography for key management and signature generation/verification.
- **github.com/btcsuite/btcd/btcutil**: Provides utilities for encoding/decoding Bitcoin addresses and other common Bitcoin operations.
- **github.com/btcsuite/btcd/chaincfg/chainhash**: Offers hash calculations and chain-related utilities.
- **github.com/btcsuite/btcwallet/wallet/txauthor**: Constructs transaction inputs/outputs and automatically handles change.
- **github.com/btcsuite/btcwallet/wallet/txrules**: Defines transaction rules, including minimum fee calculations and other constraints.
- **github.com/btcsuite/btcwallet/wallet/txsizes**: Calculates the virtual size (vSize) of transactions, enabling dynamic fee adjustments.

`gobtcsign` avoids using packages outside the `github.com/btcsuite` suite. Even so, you should never use this library directly for signing transactions without careful review to avoid potential malicious code that could collect your private keys. The best practice is to **fork the project** or copy the relevant code into your project while thoroughly reviewing the code and configuring strict network whitelists for your servers.

---

## Usage Steps

1. **Initialize Transaction Parameters**: Define transaction inputs (UTXOs), output target addresses, and amounts. Configure Replace-By-Fee (RBF) options as needed.
2. **Estimate Transaction Size and Fees**: Use the library's methods to estimate transaction size and set appropriate fees based on real-time fee rates.
3. **Generate Unsigned Transactions**: Build transactions using the defined parameters.
4. **Sign Transactions**: Sign transaction inputs using the corresponding private keys.
5. **Validate and Serialize**: Verify the signature's validity and serialize the transaction into a hexadecimal string for broadcasting.

---

## Demos

[BTC-SIGN](internal/demos/signbtc/main/main.go) [DOGECOIN-SIGN](internal/demos/signdoge/main/main.go)

---

## Notes

1. **Private Key Security**: Never expose private keys in production environments. Only use demo data for development or testing purposes.
2. **Fee Settings**: Set transaction fees reasonably based on size and network congestion to avoid rejection by miners.
3. **Change Address**: Ensure leftover funds are returned to your address as change to avoid loss of funds.
4. **Network Configuration**: Properly configure network parameters (e.g., `chaincfg.TestNet3Params`) for TestNet or MainNet usage.

---

## Getting Started with Bitcoin (BTC)

Using `gobtcsign`, here is a simple introduction to Bitcoin (`BTC`):

### Step 1 - Create a Wallet

Create a test wallet using **offline tools**. For example, see [create_wallet_test.go](create_wallet_test.go).

Avoid using online tools to create wallets, as they could expose your private key to unauthorized parties.

Blockchain wallets are created offline, and you can use any offline tool you prefer for this purpose. **Generating a private key online is insecure and should be avoided.**

### Step 2 - Obtain Test Coins from a Faucet

Look for a Bitcoin faucet online to receive some test coins. This will provide the UTXOs you need for transaction testing.

### Step 3 - Sign and Send a Transaction

Once you have UTXOs from the faucet, you can construct and send transactions.

In practice, you need extra features like block crawling to automatically get your UTXOs. Without these features, you can't fully automate sending transactions.

You can use blockchain explorers and program code to send transactions manually, while for automated transactions, block crawling is required.

### Additional - Use DOGE to Learn BTC

Since Dogecoin (DOGE) is derived from Litecoin (LTC), which itself is derived from Bitcoin (BTC), this library also supports DOGE signing.

While Litecoin signing hasn't been tested, you can try it if you want.

DOGE provides an excellent environment for learning due to its faster block times, allowing 6-block confirmation in just a few minutes. This makes testing and iteration more efficient compared to BTC, which requires around an hour for 6-block confirmation.

BTC has richer resources and greater adoption, making it more beneficial for learning blockchain concepts. Since DOGE mimics BTC, testing DOGE logic can often reveal BTC-related issues. Supporting both BTC and DOGE is a practical choice for developers.

### Important - Don’t Forget Change Outputs

Forgetting to include change outputs can lead to significant losses. Here is an example:

- Transaction at block height **818087**
- Hash: `b5a2af5845a8d3796308ff9840e567b14cf6bb158ff26c999e6f9a1f5448f9aa`
- The sender transferred **139.42495946 BTC** (worth $5,217,651), but the recipient only received **55.76998378 BTC** (worth $2,087,060).
- The remaining **83.65497568 BTC** (worth $3,130,590) was lost as miner fees.

This is a mistake that would be deeply regrettable and must be avoided.

---

## DISCLAIMER

Crypto coin, at its core, is nothing but a scam. It thrives on the concept of "air coins"—valueless digital assets—to exploit the hard-earned wealth of ordinary people, all under the guise of innovation and progress. This ecosystem is inherently devoid of fairness or justice.

For the elderly, cryptocurrencies present significant challenges and risks. The so-called "high-tech" façade often excludes them from understanding or engaging with these tools. Instead, they become easy targets for financial exploitation, stripped of the resources they worked a lifetime to accumulate.

The younger generation faces a different but equally insidious issue. By the time they have the opportunity to engage, the early adopters have already hoarded the lion’s share of resources. The system is inherently tilted, offering little chance for new entrants to gain a fair footing.

The idea that cryptocurrencies like BTC, ETH, or TRX could replace global fiat currencies is nothing more than a pipe dream. This notion serves only as the shameless fantasy of early adopters, particularly those from the 1980s generation, who hoarded significant amounts of crypto coin before the general public even had an opportunity to participate.

Ask yourself this: would someone holding thousands, or even tens of thousands, of Bitcoin ever genuinely believe the system is fair? The answer is unequivocally no. These systems were never designed with fairness in mind but rather to entrench the advantages of a select few.

The rise of cryptocurrencies is not the endgame. It is inevitable that new innovations will emerge, replacing these deeply flawed systems. At this moment, my interest lies purely in understanding the underlying technology—nothing more, nothing less.

This project exists solely for the purpose of technical learning and exploration. The author of this project maintains a firm and unequivocal stance of *staunch resistance to cryptocurrencies*.

--- 

## License

`gobtcsign` is open-source and released under the MIT License. See the [LICENSE](LICENSE) file for more information.

---

## Support

Welcome to contribute to this project by submitting pull requests or reporting issues.

If you find this package helpful, give it a star on GitHub!

**Thank you for your support!**

**Happy Coding with `gobtcsign`!** 🎉

Give me stars. Thank you!!!

---

## Starring

[![starring](https://starchart.cc/yyle88/gobtcsign.svg?variant=adaptive)](https://starchart.cc/yyle88/gobtcsign)
