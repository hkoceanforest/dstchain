// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/structs/EnumerableSetUpgradeable.sol";
import "./token.sol";
import "./Stake.sol";

interface IToken {
  function setMintAddress(address _mintAddress) external;
}

contract StakeFactory is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable{
  using SafeERC20Upgradeable for IERC20Upgradeable;
      using EnumerableSetUpgradeable
  for EnumerableSetUpgradeable.AddressSet;

  event SetYearBlockNum(uint256 oldYearBlockNum, uint256 newYearBlockNum);
  event SetDenominator(uint256 oldDenominator, uint256 newDenominator);
  event CreateStakeLog(address indexed gatewayAddress, address stakeAddress, address tokenAddress, uint256 time);

   struct StakeInfo {
     address tokenAddress;
     address stakeAddress;
     string tokenSymbol;
     uint256 decimals;
     uint256 stakeRate;
     uint256 unStakeRate;
     address gatewayAddress;
     uint256 minInflationRate;
     uint256 maxInflationRate;
     uint256 targetStakeRate;
     uint256 scaleFactor;
     uint256 chatThreshold;
   }
   
   struct CreateParams {
     address gatewayAddress;
     string name;
     string symbol;
     uint256 preMintAmount;
     uint8 tokenDecimals;
     uint256 stakeRate;
     uint256 unStakeRate;
     uint256 minInflationRate;
     uint256 maxInflationRate;
     uint256 targetStakeRate;
     uint256 scaleFactor;
     uint256 chatThreshold;
   }

   uint256 public yearBlockNum;
   uint256 public denominator;

   mapping(address => StakeInfo) public stakeInfos;
   mapping(address => EnumerableSetUpgradeable.AddressSet) private userCreates;

   function initialize() public initializer {
     __Ownable_init();
     __ReentrancyGuard_init();

     yearBlockNum = 5256000;
     denominator = 100000;
   }
   
   function setYearBlockNum(uint256 _yearBlockNum) external onlyOwner {
     
     emit SetYearBlockNum(yearBlockNum, _yearBlockNum);
     yearBlockNum = _yearBlockNum;
   }
   function setDenominator(uint256 _denominator) external onlyOwner {
     
     emit SetDenominator(denominator, _denominator);
     denominator = _denominator;
   }

   function createStake(CreateParams memory params) external nonReentrant {
      require(params.tokenDecimals <= 18, 'tokenDecimals is more than 18');
      Token tokenContract = new Token(params.name, params.symbol, params.preMintAmount,params.tokenDecimals,params.gatewayAddress);
      Stake stakeContract = new Stake(yearBlockNum,params.stakeRate,params.unStakeRate,address(tokenContract), address(tokenContract), params.gatewayAddress, params.minInflationRate, params.maxInflationRate, params.targetStakeRate, params.scaleFactor, params.chatThreshold ,denominator);
      IToken(address(tokenContract)).setMintAddress(address(stakeContract));
      stakeInfos[address(stakeContract)] = StakeInfo(
        address(tokenContract),
        address(stakeContract),
        params.symbol,
        params.tokenDecimals,
        params.stakeRate,
        params.unStakeRate,
        params.gatewayAddress,
        params.minInflationRate, 
        params.maxInflationRate, 
        params.targetStakeRate, 
        params.scaleFactor,
        params.chatThreshold
      );
      userCreates[params.gatewayAddress].add(address(stakeContract));
      emit CreateStakeLog(params.gatewayAddress,address(stakeContract), address(tokenContract), block.number);
   }

   function getUserAllCreate(address user, uint256 _page, uint256 _limit) external view returns(StakeInfo[] memory StakeArr, uint256 curPage, uint256 total){
     uint256 length = userCreates[user].length();
      StakeArr = new StakeInfo[](_limit);
      curPage = _page;
      total = length;
      uint256 start = (_page - 1) * _limit;
      uint256 end  = _page * _limit >= length ? length : _page * _limit;
      for(uint256 i = start; i < end; i++){
        StakeArr[i - start] = stakeInfos[userCreates[user].at(i)];
      }  
   }
   
}