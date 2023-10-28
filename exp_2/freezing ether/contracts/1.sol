// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FreezingEther1 {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    receive() external payable {}

    function withdraw() public {
        require(msg.sender == owner, "Not owner");
        // 错误的地址导致资金冻结
        payable(address(0)).transfer(address(this).balance);
    }
}
