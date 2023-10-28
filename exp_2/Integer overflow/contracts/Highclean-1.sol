pragma solidity 0.4.0;

contract Test {
    uint32 public a = 0xFFFFFFFF;
    uint32 public b;

    function run() public {
        var x = 1;
        a = uint32(a + x);
    }
}
