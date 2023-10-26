pragma solidity >= 0.4.24;

interface token {
    function transfer(address receiver, uint amount) external;
	function transferFrom(address from, address to, uint value) external;
    function balanceOf(address tokenOwner) constant external returns (uint balance);
}

contract againstGraveyard {
    mapping(address => bool) public active;
    mapping(address => string) public tkname;
	mapping(address => uint) public rate; //9 decimal places
	mapping(address => uint) public buyoffer; //in AGAINST
	token tokenReward = token(0xF7Be133620a7D944595683cE2B14156591EFe609);
	
	
    string public name = "AGAINST GR";
    string public symbol = "AGAINST";
    string public comment = "AGAINST Graveyard Tokens";
    address internal owner;
	
	constructor() public {
       owner = address(msg.sender); 
    }
	
	function registerToken(address _token, string _name, uint _rate, uint _buyoffer) public {
	   if (msg.sender == owner) {
	     active[_token] = false;
		 tkname[_token] = _name;
		 rate[_token] = _rate; //with 9 decimal places
		 buyoffer[_token] = _buyoffer;	//with 18 decimal places	 
	   }
	}
	
	function enableToken(address _token) public {
	   if (msg.sender == owner) {
	      active[_token] = true;
	   }
	}	
	
	function disableToken(address _token) public {
	   if (msg.sender == owner) {
	      active[_token] = false;
	   }
	}
	
	function exchangeIt(address _token, uint _qtd) public payable {
	   require(active[_token],'Token Disabled');
	   uint amount = _qtd*(10**9)*rate[_token];
	   require(amount <= buyoffer[_token]);
	   buyoffer[_token] = buyoffer[_token]-amount;
	   token deadToken = token(_token);
	   deadToken.transferFrom(msg.sender,owner,_qtd);
	   tokenReward.transfer(msg.sender, amount);
	}
	
}