// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EtherLeak3 {
    mapping(address => uint) public balances;

    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    function attack() public {
        uint bal = balances[msg.sender];
        require(bal > 0, "Insufficient funds");

        // Ether leak due to reentrancy attack
        (bool sent, ) = msg.sender.call{value: bal}("");
        require(sent, "Failed to send Ether");

        balances[msg.sender] = 0;
    }

    receive() external payable {}
}
