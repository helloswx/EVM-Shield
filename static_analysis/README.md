

## System Requirements

- Linux Distribution
- Python 3.x

## Install

#### Install Dependencies

You can install all the required packages with the following command:

```
python3 setup.py install
```

#### Install Solidity Compilers

You can install all the required Solidity compilers with the following command:

```
python3 install_compilers.py
```

## Running Script

You can run SmartMuv with the following command on the provided example smart contracts:

```
python3 -m try_smartmuv
```

Select the smart contract from the provided list, and SmartMuv will analyze and extract its complete state. 

```
1   0xc9ae11a393a08e86d46ce683fde7699db01a5f15   AUX1769
2   0x51bb7917efcad03ec8b1d37251a06cd56b0c4a72   DSRCoin
3   0x24dd6e1fe742bd8fd3a1d144fece1680f16296aa   OBK
4   0x143e685dd51d467d77663a3be119217185d81b99   CommunityBankCoin
5   0x145f9bbd9f1ca0923e81e05c2ac04fda2310d774   VACCToken

```


#### Slot Layout

```
slot 0 - mapping balances[address] = uint256;
slot 1 - mapping allowed[address][address] = uint256;
slot 2 - uint256 totalSupply;
slot 3 - string name;
slot 4 - uint8 decimals;
slot 5 - string symbol;
```
