pragma solidity ^0.8.3;

contract TestDelegateCall{
    uint public num;
    uint private number;
    address public sender;
    uint public value;

    function setVars(uint _num) external payable {
        num =_num;
        number = _num;
        sender = msg.sender;
        value = msg.value;
    }

}

contract DelegateCall{
    uint private number;
    address public sender;
    uint public value;

    function setVars(address _test, uint _num) external payable{
        (bool success,bytes memory data) =_test.delegatecall(
            abi.encodeWithSelector(TestDelegateCall.setVars.selector,_num)
        );
        require(success,"delegatecall failed");
    }
}
