// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

interface IRewardToken {
  function stakeMint(address to, uint256 amount) external;
}

contract Stake is ReentrancyGuard{
   using SafeERC20 for IERC20;

    event StakedLog(address indexed user, uint256 amount,uint256 time);
    event WithdrawnLog(address indexed user, uint256 amount,uint256 time);
    event RewardPaidLog(address indexed user, uint256 reward,uint256 time);

    address public burnAddress;
    address public rewardsToken;
    address public stakingToken;
    address public gatewayAddress;
    uint256 public rewardPerTokenStored;
    uint256 public lastUpdateTime;
    uint256 public totalSupply;
    uint256 public denominator;
    uint256 public minInflationRate;
    uint256 public maxInflationRate;
    uint256 public targetStakeRate;
    uint256 public scaleFactor;
    
    uint256 public yearBlockNum;
    uint256 public stakeRate;
    uint256 public unStakeRate;

    struct UserInfo {
        uint256 userRewardPerTokenPaid;
        uint256 rewards;
        uint256 balance;
        uint256 stakeRate;
        uint256 unStakeRate;
    }
    
   mapping(address => UserInfo) public userInfo;
   uint256 public chatThreshold;
   constructor(uint256 _yearBlockNum, uint256 _stakeRate, uint256 _unStakeRate, address _rewardsToken, address _stakingToken, address _gatewayAddress,uint256 _minInflationRate,uint256 _maxInflationRate, uint256 _targetStakeRate, uint256 _scaleFactor,uint256 _chatThreshold, uint256 _denominator) {
      denominator = _denominator;
      burnAddress = address(0x000000000000000000000000000000000000dEaD);
      rewardsToken = _rewardsToken;
      stakingToken = _stakingToken;
      gatewayAddress = _gatewayAddress;
      minInflationRate = _minInflationRate;
      maxInflationRate = _maxInflationRate;
      targetStakeRate = _targetStakeRate;
      scaleFactor = _scaleFactor;
      stakeRate = _stakeRate;
      unStakeRate = _unStakeRate;
      yearBlockNum = _yearBlockNum;
      chatThreshold = _chatThreshold;
   }
   
   function isCanChat(address account) public view returns (bool) {
     return userInfo[account].balance >= chatThreshold;
   }

    function getRewardPerTime() public view returns(uint256 amount){
      uint256 total = IERC20(stakingToken).totalSupply();
      uint256 burnAmount = IERC20(stakingToken).balanceOf(burnAddress);
      uint256 netStakeRate = totalSupply * denominator/ (total - burnAmount);
      // int256 totalInflationRate = netStakeRate * 1e18/targetStakeRate >= 1e18 ? -int256((netStakeRate * 1e18/targetStakeRate - 1e18) * scaleFactor)/1e18: int256((1e18 - netStakeRate * 1e18/targetStakeRate) * scaleFactor / 1e18);
      // uint256 rate = totalInflationRate <= int256(minInflationRate) ? minInflationRate : totalInflationRate >= int256(maxInflationRate) ? maxInflationRate : totalInflationRate; 
      uint256 rate;
      if(netStakeRate * 1e18/targetStakeRate >= 1e18) {
        rate = minInflationRate;
      } else {
        uint256 totalInflationRate = (1e18 - netStakeRate * 1e18/targetStakeRate) * scaleFactor / 1e18;
        rate = totalInflationRate <= minInflationRate ? minInflationRate : totalInflationRate >= maxInflationRate ? maxInflationRate : totalInflationRate;  
      }
      amount = (total - burnAmount) * rate / denominator / yearBlockNum;
    }
     function rewardPerToken() public view returns (uint256 perTokenStored,uint256 mintAmount) {
        if (totalSupply == 0) {
            return (rewardPerTokenStored, 0);
        }
        uint256 addPerTokenStored =  getRewardPerTime() * (block.number - lastUpdateTime) *
            1e18 / totalSupply;
        return (rewardPerTokenStored + addPerTokenStored, getRewardPerTime() * (block.number - lastUpdateTime));
    }
    function earned(address account) public view returns (uint256) {
        (uint256 perTokenReward,) = rewardPerToken();
        return
            (userInfo[account].balance *
                (perTokenReward - userInfo[account].userRewardPerTokenPaid)) /
            1e18 +
            userInfo[account].rewards;
    }
      function stake(uint256 amount) external nonReentrant updateReward(msg.sender)  {
        require(amount > 0, "Cannot stake 0");
        uint256 fee = amount * stakeRate / denominator;
        uint256 burnFee = fee / 3;
        totalSupply = totalSupply + amount - fee;
        userInfo[msg.sender].balance = userInfo[msg.sender].balance + amount - fee;
        IERC20(stakingToken).safeTransferFrom(msg.sender, address(this), amount);
        IERC20(stakingToken).safeTransfer(burnAddress, burnFee);
        IERC20(stakingToken).safeTransfer(gatewayAddress, fee - burnFee);
        emit StakedLog(msg.sender, amount - fee, block.number);
     }
    function withdraw(uint256 amount) public nonReentrant updateReward(msg.sender) {
        require(amount > 0, "Cannot withdraw 0");
        require(amount <= userInfo[msg.sender].balance, "balance < amount");
        totalSupply = totalSupply - amount;
        uint256 fee = amount * unStakeRate / denominator;
        uint256 burnFee = fee / 3;
        userInfo[msg.sender].balance = userInfo[msg.sender].balance - amount;
        IERC20(stakingToken).safeTransfer(msg.sender, amount - fee);
        IERC20(stakingToken).safeTransfer(burnAddress, burnFee);
        IERC20(stakingToken).safeTransfer(gatewayAddress, fee - burnFee);
        emit WithdrawnLog(msg.sender, amount,block.number);
    }
     modifier updateReward(address account) {
        (uint256 stored, uint256 mintAmount) = rewardPerToken();
        IRewardToken(rewardsToken).stakeMint(address(this), mintAmount);

        rewardPerTokenStored = stored;
        lastUpdateTime = block.number;
        if (account != address(0)) {
            userInfo[account].rewards = earned(account);
            userInfo[account].userRewardPerTokenPaid = rewardPerTokenStored;
        }
        _;
    }
   
  function getReward() public nonReentrant updateReward(msg.sender) {
        uint256 rewardAmount = userInfo[msg.sender].rewards;
        if (rewardAmount > 0) {
            userInfo[msg.sender].rewards = 0;
            IERC20(rewardsToken).safeTransfer(msg.sender, rewardAmount);
            emit RewardPaidLog(msg.sender, rewardAmount, block.number);
        }
    }
}