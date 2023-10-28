// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UninitializedStruct4 {
    struct InnerData {
        uint x;
    }

    struct OuterData {
        InnerData inner;
        uint y;
    }

    OuterData public data;

    // Vulnerability: Uninitialized nested struct
    function setOuterData(uint _y) public {
        OuterData memory newData; // Uninitialized memory variable
        newData.y = _y;
        // newData.inner is uninitialized
        data = newData;
    }
}
