// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UninitializedStorage {
    struct Data {
        uint a;
        uint b;
    }

    Data public data;
    Data public defaultData;

    // Vulnerability: Uninitialized storage pointer
    function setData(uint _a, uint _b) public {
        Data storage d; // Uninitialized storage variable
        d.a = _a;
        d.b = _b;

        // The assignment above actually modifies 'defaultData', not 'data'
    }
}
