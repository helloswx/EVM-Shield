// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UninitializedStruct3 {
    struct Data {
        uint a;
        uint b;
    }

    Data public data;

    // Vulnerability: Partially initialized struct
    function setData(uint _a) public {
        Data memory newData; // Uninitialized memory variable
        newData.a = _a;
        // newData.b is not initialized
        data = newData;
    }
}
