#! python3.13 -i

import chaindb

RequiredConfirmations: int = 10
# 表示区块链节点间相互同步的数据的可靠性.
# 可靠性会导致更高的延迟.  基本上
#     延迟 = RequiredConfirmations * 1.25 seconds
# 该变量可配置, 建议所有节点使用相同的 RequiredConfirmations.

uc = chaindb.UserClient()

AvailableOwners = range(100_0000, 200_0000)
# ChainDB 所存储的 key-value pairs 中 key 的可用范围.


def Owners() -> list[int]:
    """
    获取 ChainDB 里的 keys (也称 owners).
    """
    owners = uc.ListOwners(confirmations=RequiredConfirmations)
    return [owner.OwnerID for owner in owners if owner.OwnerID in AvailableOwners]


def Get(owner: int) -> list[str]:
    """
    按 owner 获取 value, value 是一个数组, 按插入时间的顺序排列.
    """
    txs = uc.ListTransactionsOwnedBy(owner, confirmations=RequiredConfirmations)
    txs.sort(key=lambda tx: tx.Transaction.Nonce)
    return [tx.Transaction.Data for tx in txs]


def Insert(owner: int, data: str) -> int:
    """
    向 ChainDB 插入数据.
    也即, 在 owner 映射的数组后追加一个 data.

    新建键值对 (owner -> [data, ...]) 时, 请用随机数作为 owner, 避免冲突.

    返回被插入的数据在 owner 的数组中的索引.
    """
    if owner not in AvailableOwners:
        raise ValueError(f"{owner=} isn't in {AvailableOwners=}")

    nonce = len(Get(owner))

    while True:
        ok = uc.Insert(
            chaindb.blockchain.Transaction(
                OwnerID=owner,
                Nonce=nonce,
                Data=data,
            ),
            confirmations=RequiredConfirmations,
        )
        if ok:
            return nonce

        if len(Get(owner)) > nonce:
            if Get(owner)[nonce] == data:
                return nonce
            return Insert(owner, data)
