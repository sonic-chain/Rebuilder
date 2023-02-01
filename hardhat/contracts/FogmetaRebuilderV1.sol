// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract FogmetaRebuilderV1 is Ownable {
    mapping(string => mapping(address => uint)) public cidDeposits;

    event Deposit(address account, string payloadCid, uint amount);
    event Withdraw(address account, string payloadCid, uint amount);

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

}

