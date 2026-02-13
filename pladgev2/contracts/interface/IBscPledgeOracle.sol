// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

interface IBscPledgeOracle {
    // 根据资产地址查询价格
    function getPrice(address asset) external view returns (uint256);
    // 获取代币标的价格
    function getUnderlyingPrice(uint256 cToken) external view returns (uint256);
    // 批量查询资产价格
    function getPrices(uint256[] calldata assets) external view returns (uint256[] memory);
}
