# EVM-Shield

- **AST_JSON**
  - AST files for all contracts
- **CSV**
  - Static analysis results
- **contracts**
  - All contract datasets
- **exp_1**
  - Replay vulnerability detection results for 8 million blocks
- **exp_2**
  - Vulnerability contracts used for comparison experiments
- **geth_shield**
  - geth, an Ethereum client embedded with EVM-Shield
- **static_analysis**
  - Static analysis is performed on the contract code and the abstract syntax tree
- **policy translator**
  - Policy translator

## How to use EVM-Shield
Setting Up a Private Chain with Geth and Connecting to Remix. If you want to fully synchronize the mainnet, simply configure mongoDB and change the shield in contract.go to record instead.

### 1. Compile Geth from Source

```bash
git clone https://github.com/helloswx/EVM-Shield.git
cd geth_shield
make geth
make all
```
### 2. Initialize the Private Chain
For initializing a new private chain, you'll need a genesis block file, e.g., genesis.json.

Use the above genesis block to initialize the private chain:
```bash
./build/bin/geth init genesis.json --datadir ./myPrivateNetwork
```

### 3. Start the Private Chain

```bash
./build/bin/geth --datadir ./myPrivateNetwork --networkid 15 --nodiscover --rpc --rpcapi="personal,db,eth,net,web3,txpool,miner" --rpcport "8545" --rpcaddr "127.0.0.1" --rpccorsdomain "*" --allow-insecure-unlock console
```
### 4. Connect to Remix
Finally, you can send the hexadecimal policy to the smart contract through Low level interactions in remix
