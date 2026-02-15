// SPDX-License-Identifier: MIT

pragma solidity 0.6.12;

interface IMultiSignature {
    // 多签接口 获取验证签名
    function getValidSignature(bytes32 msghash, uint256 lastIndex) external view returns (uint256);
}

// 多签客户端 使用伪随机存储槽存储多签服务端合约地址
// validCall 给业务合约关键函数功能加锁
contract multiSignatureClient {
    // 伪随机存储槽
    uint256 private constant multiSignaturePositon = uint256(keccak256("org.multiSignature.storage"));
    uint256 private constant defaultIndex = 0;

    constructor(address multiSignature) public {
        // 防止初始化错误
        require(multiSignature != address(0), "multiSignatureClient : Multiple signature contract address is zero!");
        saveValue(multiSignaturePositon, uint256(multiSignature));
    }

    function getMultiSignatureAddress() public view returns (address) {
        return address(getValue(multiSignaturePositon));
    }

    // 验证锁
    modifier validCall() {
        // 多签校验逻辑
        checkMultiSignature();
        _;
    }

    // 生成哈希：keccak256(abi.encodePacked(msg.sender, address(this)))。它将调用者地址和当前合约地址绑定，生成一个唯一的业务标识。
    // 外部调用：调用 IMultiSignature 接口的 getValidSignature 函数。
    // 逻辑判定：要求返回的 newIndex 必须大于 defaultIndex (0)。如果返回 0，说明多签合约里还没通过这笔交易的授权，直接 revert。
    function checkMultiSignature() internal view {
        uint256 value;
        assembly {
            value := callvalue()
        }
        // 哈希拼接
        bytes32 msgHash = keccak256(abi.encodePacked(msg.sender, address(this)));
        address multiSign = getMultiSignatureAddress();
        //        uint256 index = getValue(uint256(msgHash));
        // 去多签服务端查找 调用者+业务合约 的管理员认证数是否达到阈值
        uint256 newIndex = IMultiSignature(multiSign).getValidSignature(msgHash, defaultIndex);
        require(newIndex > defaultIndex, "multiSignatureClient : This tx is not aprroved");
        //        saveValue(uint256(msgHash),newIndex);
    }

    function saveValue(uint256 position, uint256 value) internal {
        assembly {
            sstore(position, value)
        }
    }

    function getValue(uint256 position) internal view returns (uint256 value) {
        assembly {
            value := sload(position)
        }
    }
}
