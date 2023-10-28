pragma solidity 0.4.0;

contract NameRegistrar{

    bool public unlocked = false;
    struct NameRecord{
        bytes32 name;
        address mappedAddress; 
    }
    
    mapping(address => NameRecord) public registeredNameRecord;
    mapping(bytes32 => address) public resolve;


    function register(bytes32 _name, address _mappedAddress) public{
        NameRecord storage newRecord;
        newRecord.name = _name;
        newRecord.mappedAddress = _mappedAddress;

        resolve[_name] = _mappedAddress;
        registeredNameRecord[msg.sender] =newRecord;  

    }

}
