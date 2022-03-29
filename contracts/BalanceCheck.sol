// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.0;

contract BalanceCheck {
    uint256 balance;
    address public admin;

    constructor() {
        admin = msg.sender;
        balance = 0;
        updateBalance();
    }

    function updateBalance() internal {
        balance += msg.value;
    }

    function Withdraw(uint256 _amt) public{
        require(msg.sender == admin);
        balance = balance - _amt;
    }

    function Deposit(uint256 amt) public returns (uint256) {
        return balance = balance + amt;
    }

    function Balance() public view returns (uint256) {
        return balance;
    }

}
