// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract FogmetaRebuilderV1 is Ownable {
    address accountAddress;
    event AddressBalance(address from, address to, uint amount);

    constructor() {
        accountAddress = 0xBf4eF4147Aac5FD3C1F8b6b4B8c2F2A70Fb0efF1;
    }

    error InsufficientBalance(uint requested, uint available);

    function transfer(uint amount) public onlyOwner{
        if (address(msg.sender).balance < amount )
            revert InsufficientBalance({
            requested: amount,
            available: address(msg.sender).balance
            });
        payable(accountAddress).transfer(amount);
        emit AddressBalance(address(msg.sender),accountAddress,amount);
    }

    function withdraw(uint amount) public returns(bool){
        payable(msg.sender).transfer(amount);
        return true;
    }

    function getBalance() view public returns(uint){
        return address(this).balance;
    }
}
