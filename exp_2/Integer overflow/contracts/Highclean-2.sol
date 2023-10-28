pragma solidity 0.4.0;

contract Test {
    uint32 public a = 0xFFFFFFFF;
    bool public b;

    function run() public {
        var x = 1;
        a = uint32(a + x);
    }
}
