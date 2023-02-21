// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@zondax/filecoin-solidity/contracts/v0.8/MarketAPI.sol";
import "@zondax/filecoin-solidity/contracts/v0.8/types/MarketTypes.sol";
import "@zondax/filecoin-solidity/contracts/v0.8/cbor/BigIntCbor.sol";

contract FogmetaRebuilder is Ownable {
    MarketTypes.GetDealTermReturn public dealTerm;
    MarketTypes.GetDealEpochPriceReturn public dealPrice;

    mapping(uint64 => DealData) public dealInfo;

    struct DealData {
        uint storagePrice;
        uint balance;
        uint replicaFee;
    }


    function deposit(uint64 dealId) public payable {
        uint storagePrice = getStoragePrice(dealId);
        require(msg.value >= storagePrice, 'not enough deposit');

        DealData storage d = dealInfo[dealId];
        d.storagePrice = storagePrice;
        d.balance += msg.value;
        d.replicaFee = storagePrice * 5 / 100;
    }

    function replicateDeal() public onlyOwner {
        // uint fee = cidBalance[payloadCid] * 5 / 100;
        // cidBalance[payloadCid] -= fee;
    }

    function getStoragePrice(uint64 dealId) public returns (uint) {
        dealTerm = MarketAPI.getDealTerm(dealId);
        dealPrice = MarketAPI.getDealTotalPrice(dealId);

        return uint(uint64(dealTerm.end - dealTerm.start)) ;//* uint(bytes32(dealPrice.price_per_epoch.val));
    }
}

