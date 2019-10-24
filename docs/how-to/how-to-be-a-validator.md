# How To Run A Validator

Start by reading the [validator overview](../overviews/what-are-validators.md).

>Note: This guide covers the bare minimum needed to get a validator up and running. To operate a validator professionally it is recommended that you consult the cosmos docs on secure private keys storage, sentry node architectures, slashing risks, etc.

>Note: This guide is for creating a validator on a live chain, not for participating in genesis.

## 1) Setup A Compute Environment

Setup a cloud Ubuntu server with an externally accessible IP.

## 2) Install the Software

1) install go - there are several good guides online
2) clone the repo
3) run `make install`

This will install `kvd` and `kvcli` into `$GOPATH/bin` (`normally ~/go/bin`). Add this to your `$PATH` for convenience.

## 3) Setup

Run `kvd init <your validator name>` to create the initial data directories and validator keys.

<!-- TODO say where they can get the genesis file from -->
Download the genesis file for the chain. Place this in the config directory; default location `$HOME/.kvd/config/genesis.json`.

<!-- TODO are we setting up a seed node? -->
Add some peer addresses for your node to connect to, in the node config file (default `$HOME/.kvd/config/config.toml`). Ask online for these addresses.

Start the node with `kvd start`. The node should connect to peers, and start to catch up with the chain. Once caught up you have a full node.

## 4) Convert the Full Node To A Validator

Setup `kvcli` as described in [the cli guide](./how-to-use-the-command-line-wallet.md). By default `kvcli` will connect to your node on `localhost`.

You will need a minimum of 1 kava token (plus gas fees) in a local account to run a validator.

Create a validator with `kvcli tx staking create-validator ...`.

e.g.

```bash
    kvcli tx staking create-validator \
        --amount 1000000ukava \
        --pubkey $(kvd tendermint show-validator) \
        --from <your key name> \
        --gas-adjustment 1.3 \
        --moniker <your validator name> \
```

The `--help` flag will print out more information on the available options. Certain details can be edited with `kvcli tx staking edit-validator ...`.

## Troubleshooting

### Unjailing

If your validator is offline too long some stake is slashed and your validator becomes jailed. To unjail, submit the `kvcli tx slashing unjail`.

### Slow Catching Up To The Latest Block Height

A slow network connection can make syncing with the network very slow. Try allotting a higher bandwidth connection.

### Error Starting a Full Node

Check that you are connecting to peers and that `persistent_peers` or `seeds` in `config.toml` is set correctly.

`kvd unsafe-reset-all` is a useful command to delete all blockchain history stored locally and start syncing again from the genesis file.
