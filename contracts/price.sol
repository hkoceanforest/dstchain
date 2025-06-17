
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IRouter{
  function getAmountsOut(uint amountIn, address[] calldata path) external view returns (uint256[] memory amounts);
}

contract Swap  {

    address public controller;
    address public USDT;
    address public WDST;
    IRouter public router;
 
   constructor() {
      router = IRouter(0xFd91bA542d92a16E74Ee85f1811297548CE72f82);
      USDT = address(0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd);
      WDST = address(0x80Cc816f98ecEc5EafE363c17C5d2a9ac890C84b);
   }
  function getPrice(uint256 _amountUsdtIn) external view returns (uint256 amountDstOut) {
      address[] memory path = new address[](2);
      path[0] = USDT;
      path[1] = WDST;
      uint256[] memory amounts = IRouter(router).getAmountsOut(_amountUsdtIn, path);
      amountDstOut = amounts[1];
    }
}
