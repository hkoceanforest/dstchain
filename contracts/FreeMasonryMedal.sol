// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

 struct Attribute {
      uint256 tokenId;
      address createAddress;
      uint256 level;
      uint256 createTime;
    }

contract FreeMasonryMedal is ERC721, Ownable, ReentrancyGuard{
    using Strings for uint256;
    using EnumerableSet for EnumerableSet.AddressSet;
    using EnumerableSet for EnumerableSet.UintSet;

    event CreateLog(address indexed createAddress, uint256 tokenId, uint256 level, uint256 createTime);

   string public baseUrl;
   uint256 private _totalSupply;
   mapping(address => EnumerableSet.UintSet) private _userTokens;
   mapping(uint256 => Attribute) public tokenInfo;
  
  /**
   * @dev Initialization parameters
   */
   constructor () ERC721('Free Masonry Medal', 'Free Masonry Medal') {
     baseUrl = '';
   }
  
  function setUrl (string memory url) external onlyOwner {
    baseUrl = url;
  }

  function mint(uint256 level) external nonReentrant {
     _totalSupply = _totalSupply + 1;
     _mint(msg.sender, _totalSupply);
     tokenInfo[_totalSupply] = Attribute(_totalSupply, msg.sender, level, block.timestamp);
     emit CreateLog(msg.sender, _totalSupply, level, block.timestamp);  
  }
  function totalSupply () external view returns(uint256) {
    return _totalSupply;
  }
  function tokenURI(uint256 tokenId) public view override returns (string memory) {
       Attribute memory attribute = tokenInfo[tokenId];
       return string(
                abi.encodePacked(
                                '{"name": "Free Masonry Medal #',
                                tokenId.toString(),
                                '", "image": "',
                                string(abi.encodePacked(baseUrl,'/', attribute.level.toString(), '.png')),
                                '", "attributes":[',
                                  '{"trait_type":"level", "value":"',
                                    attribute.level.toString(),
                                  '"},',
                                  '{"trait_type":"createTime", "value":"',
                                    attribute.createTime.toString(),
                                  '"},',
                                   '{"trait_type":"createAddress", "value":"',
                                     Strings.toHexString(uint160(attribute.createAddress), 20),
                                  '"}]}'
                            )
            );
    }
  function getUserNfts (address _user) external view returns(Attribute[] memory tokens){
    uint256 length = _userTokens[_user].length();
    tokens = new Attribute[](length);
    for(uint256 i = 0; i < length; i++){
      tokens[i] = tokenInfo[_userTokens[_user].at(i)];
    }  
  }
  function _afterTokenTransfer(address from, address to,uint256 tokenId) override internal virtual {
    if(from != address(0)) {
      _userTokens[from].remove(tokenId);
    }
    if(to != address(0)) {
      _userTokens[to].add(tokenId);  
    }
  }
}
