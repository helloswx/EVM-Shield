// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TodBidding {
    address public highestBidder;
    uint public highestBid;

    // Anyone can place a bid
    function bid() public payable {
        require(msg.value > highestBid, "Bid not high enough");

        // Refund the previous highest bidder
        if (highestBidder != address(0)) {
            payable(highestBidder).transfer(highestBid);
        }

        // TOD Vulnerability: the state is changed after sending Ether
        highestBidder = msg.sender;
        highestBid = msg.value;
    }
}
