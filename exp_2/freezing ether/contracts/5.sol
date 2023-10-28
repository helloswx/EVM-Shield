// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FreezingEther3 {
    // 合约接收 Ether 但没有提供提款方法
    receive() external payable {}
}