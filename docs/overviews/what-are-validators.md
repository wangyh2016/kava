# What are validators

Kava is based on the Cosmos SDK and uses tendermint consensus. This relies on a fixed set of validators; full nodes responsible for running the network and keeping it secure. Validators commit new blocks to the blockchain by listening for incoming transactions submitted by users, forming them into new blocks, and voting on blocks proposed by other validators.

A validator's weight in voting on new blocks is proportional to it's stake, the amount of tokens delegated to it from it's own funds and others. This stake is subject to slashing (an in protocol burning) if the validator is detected behaving unfavorably to the network, such as not being online, or double signing blocks.

For this work validators are rewarded with inflationary rewards and transaction fees. They are free to share a portion of these rewards with their delegators. Validators also take an active role in network governance by proposing and voting on changes to the software and network parameters.

Since there is a relatively small set of validators the security and performance requirements are much higher than in other blockchains. It is recommended that they run in high-availability setups with dedicated hardware, redundant power and networking. Security is paramount, with HSMs recommended for storing private keys, and a sentry node architecture to protect against DDoS attacks.

If you're interested in becoming a professional validator, refer to cosmos hub documentation for more information.
