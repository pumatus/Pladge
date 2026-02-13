// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import "../multiSignature/multiSignatureClient.sol";
import "@openzeppelin/contracts/utils/EnumerableSet.sol";

contract AddressPrivileges is multiSignatureClient {
    constructor(address multiSignature) public multiSignatureClient(multiSignature) {}

    using EnumerableSet for EnumerableSet.AddressSet;
    EnumerableSet.AddressSet private _minters;

    // 添加铸币者
    function addMinter(address _addMinter) public validCall returns (bool) {
        require(_addMinter != address(0), "Tokenn:_addMinter is the zero address");
        return EnumerableSet.add(_minters, _addMinter);
    }

    // 删除铸币者
    function delMinter(address _delMinter) public validCall returns (bool) {
        require(_delMinter != address(0), "Token: _delMinter is the zero address");
        return EnumerableSet.remove(_minters, _delMinter);
    }

    // 获取铸币者列表长度
    function getMinterLength() public view returns (uint256) {
        return EnumerableSet.length(_minters);
    }

    // 是否存在铸币者
    function isMinter(address account) public view returns (bool) {
        return EnumerableSet.contains(_minters, account);
    }

    // 根据索引获取铸币者账户
    function getMinter(uint256 _index) public view returns (address) {
        require(_index <= getMinterLength() - 1, "Token: index out of bounds");
        return EnumerableSet.at(_minters, _index);
    }

    // 铸币者修饰符
    modifier onlyMinter() {
        require(isMinter(msg.sender), "Token: caller is not the minter");
        _;
    }
}
