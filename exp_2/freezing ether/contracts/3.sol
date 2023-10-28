// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FreezingEther4 {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    function withdraw() public {
        // 仅当调用者不是所有者时才允许提款，这是一个逻辑错误
        require(msg.sender != owner, "Owner cannot withdraw");
        payable(msg.sender).transfer(address(this).balance);
    }

    receive() external payable {}
}