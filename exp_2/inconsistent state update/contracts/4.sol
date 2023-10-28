// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InconsistentState4 {
    bool public locked;

    function lock() public {
        require(!locked, "Already locked");

        // Executing an action and updating the locked state later can result in duplicate operations
        doSomething();
        
        locked = true;
    }

    function doSomething() private {
        // xxx
    }
}
