// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract FogmetaRebuilderV1 is Ownable {
    mapping(string => mapping(address => uint)) public cidDeposits;
    address accountAddress = 0xBf4eF4147Aac5FD3C1F8b6b4B8c2F2A70Fb0efF1;

    event Deposit(address account, string payloadCid, uint amount);
    event Withdraw(address account, string payloadCid, uint amount);
    event AddressBalance(address from, address to, uint amount);

    function deposit(string memory payloadCid) public payable {
        cidDeposits[payloadCid][msg.sender] += msg.value;

        emit Deposit(msg.sender, payloadCid, msg.value);
    }

    function withdraw(string memory payloadCid, address from, uint amount) public onlyOwner {
        require(cidDeposits[payloadCid][from] >= amount, 'balance too low');
        cidDeposits[payloadCid][from] -= amount;
        payable(msg.sender).transfer(amount);

        emit Withdraw(from, payloadCid, amount);
    }

    function transfer(uint amount) public onlyOwner{
        if (address(msg.sender).balance < amount )
            revert InsufficientBalance({
                requested: amount,
                available: address(msg.sender).balance
            });
        address(msg.sender).balance -= amount;
        accountAddress += amount;
        emit AddressBalance(address(msg.sender),accountAddress);
    }

    error InsufficientBalance(uint requested, uint available);

}

