// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InconsistentState3 {
    mapping(address => bool) public voted;
    uint public yesVotes;
    uint public noVotes;

    function vote(bool _vote) public {
        require(!voted[msg.sender], "Already voted");

        if (_vote) {
            yesVotes += 1;
        } else {
            noVotes += 1;
        }

        // 【*】
        voted[msg.sender] = true;
    }
}
