// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

interface ERC20Interface {
    // 根据用户地址查询余额
    function balanceOf(address user) external view returns (uint256);
}
