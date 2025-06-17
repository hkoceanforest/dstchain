// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract Token is ERC20 {
  
  event StakeMintLog(address indexed to, uint256 amount);
 
  uint8 decimalsBase;
  address mintAddress;
  constructor(string memory name, string memory symbol, uint256 preMintAmount, uint8 tokenDecimals, address receiveAddress) ERC20(name, symbol) {
    decimalsBase = tokenDecimals;
    _mint(receiveAddress, preMintAmount);
  }
  function setMintAddress(address _mintAddress) external {
    if(mintAddress == address(0)) {
       mintAddress =  _mintAddress;
    }
  }
  function stakeMint(address to, uint256 amount) external {
    require(mintAddress == msg.sender, 'not mint address');
    _mint(to, amount);
    emit StakeMintLog(to, amount);
  }
  function decimals() public view virtual override returns (uint8) {
    return decimalsBase;
  }
}