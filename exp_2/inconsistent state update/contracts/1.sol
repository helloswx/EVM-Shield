// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InconsistentState1 {
    uint public balanceReceived;

    receive() external payable {
        // Increasing the amount first and receiving the payment later may result in an inconsistent state due to a failed call：
        balanceReceived += msg.value;
        
    }
}
