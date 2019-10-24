# What Is Kava

Kava is a blockchain for cross-chain DeFi, built with the power of the Cosmos SDK.

It will offer decentralized loans and stable coins for other crypto assets such as BTC, XRP, BNB, and ATOM.
The network is secured through tendermint PBFT and delegated proof of stake.

## How does it work

Decentralized loans are created by locking up collateral and then minting new stable tokens up to some fraction fo the stable asset. If the price of the collateral goes up, more stable tokens can be withdrawn. If the price goes down, a loan will be automatically seized and liquidated by others on the network, with the original creator suffering a fee.

## The Network

Kava is run by approximately 100 validators who keep the network live and secure. However every account can contribute to the security of the network, and earn block rewards for doing so. This is done by users staking funds to validators. Staked assets are subject to slashing if validators act to the detriment of network security. But they also receive inflationary increases in the token supply, and a portion of transaction fees.

## What is the Cosmos SDK and IBC

The Cosmos SDK is a software development framework for building interoperable, fast, PoS blockchains. IBC is a protocol for connecting all blockchains focused on those built with the SDK; the cosmos ecosystem.

IBC allows blockchains to interchange arbitrary messages in a secure way. This enables many applications, including transferring tokens between blockchains.

Kava is built using the Cosmos SDK, and inherits its many benefits. As such much of the SDK documentation is applicable to the Kava software.

## Executables

This repo contains the kava blockchain full node and validator software (`kvd`), and a cli wallet (`kvcli`).

- `kvd` runs a full node, and can additionally run as a validator.
- `kvcli` is a cli to interact with a full node, by sending transactions or querying the chain. Additionally it can be run as a rest-server.

Configuration for each executable is stored by default in `$HOME/.kvd` and `$HOME/.kvcli` respectively.

## Interacting With The Live Chain

The live chain can be visualized using the community block explorers.

There are also community wallets to send transactions easily.

The cli wallet can be used with any full node, remote or local.
