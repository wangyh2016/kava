# Using The Command Line Wallet

The command line wallet `kvcli` is an interface to submit transactions (such as sending coins and delegating) to the blockchain. It can also be used to query the state of the blockchain.

This guide will show basic setup and usage.

>Note: this is not exhaustive, the reference help documentation for each command can be accessed with the `--help` flag. Also the cosmos hub documentation is helpful due to the shared code base.

## 1) Install The CLI

1) install go - there are several good guides online
2) clone the repo
3) run `make install`

This will install `kvcli` into `$GOPATH/bin` (`normally ~/go/bin`). Add this to your `$PATH` for convenience.

## 2) Setup

`kvcli` requires a full node to connect to. This is configured with the flag `--node <full node url>` or setting it persistently with `kvcli config node <full node url>`.

To get a full node address ask online, or run your own.

Every command that communicates with a full node also needs the `--chain-id <chain id>` flag. The first mainnnet chain id is `kava-1`. This can also be set persistently with `kvcli config chain-id <chain-id>`.

## 3) Run Commands

All commands for submitting a transaction to the blockchain are under `kvcli tx`.

Commands for querying the state of the chain are under `kvcli query`.

And commands for managing private keys are under `kvcli keys`.

<!-- TODO notes on using a ledger -->

<!-- TODO notes on using a multisig - refer to cosmos? -->

### Create a key and receive tokens

Create a new local private key with `kvcli keys add <key name>`. It will ask you to enter the password to encrypt the key.

Then `kvcli keys show <key name>` will display the address associated with this key.

Send this address to your counter-party for them to transfer tokens to it.

Check your token balance with `kvcli query account <your address>`. Note accounts won't exist on chain until they have received tokens.

### Send Some Tokens

`kvcli tx send <key name you want to send from> <receiver's address> <amount>`

Note the amount must be specified as micro kava, where `1000000ukava` is 1 kava.

e.g. `kvcli tx send my_key kava15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw 12000000ukava` will send 12 kava to the address `kava15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw`.

### Other Commands

Other common actions include delegating to a validator with `kvcli tx staking delegate ...`, voting on proposals `kvcli tx gov vote ...` and querying the chain for various state. The help text documents the arguments needed.

<!-- TODO Notes on Gas usage:
--gas, --fees, --gas-adjustment, --gas-prices -->
