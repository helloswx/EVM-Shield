// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FreezingEther5 {
    uint public unlockTime;

    constructor() {
        unlockTime = block.timestamp + 1 years; // 设定一年后解锁
    }

    function attack() public {
        // 只有在特定时间之后才允许提款
        require(block.timestamp > unlockTime, "Funds are locked");
        payable(msg.sender).transfer(address(this).balance);
    }

    //禁用掉此函数防止存钱
    receive() external payable {}
}
