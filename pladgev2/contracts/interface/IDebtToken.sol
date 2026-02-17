// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity 0.6.12;

interface IDebtToken {
    // 债务
    // 查用户余额
    function balanceOf(address account) external view returns (uint256);
    // 查存在的代币数量
    function totalSupply() external view returns (uint256);
    // 铸造
    function mint(address accounnt, uint256 amount) external;
    // 销毁
    function burn(address account, uint256 amount) external;
}
