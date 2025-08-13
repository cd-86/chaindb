import logging
import functools

# polyfills
functools.cache = functools.lru_cache(None)

_Logger = logging.getLogger(__name__)

import chaindb

RequiredConfirmations: int = 10
AvailableOwners = range(100_0000, 200_0000)


@functools.cache
def _Client():
    _Logger.info("正在创建 ChainDB 客户端...")
    uc = chaindb.UserClient()
    _Logger.info("ChainDB 客户端创建完成")
    return uc


def Owners() -> list:
    owners = _Client().ListOwners(confirmations=RequiredConfirmations)
    return [owner.OwnerID for owner in owners if owner.OwnerID in AvailableOwners]


def Get(owner: int) -> list:
    txs = _Client().ListTransactionsOwnedBy(owner, confirmations=RequiredConfirmations)
    txs.sort(key=lambda tx: tx.Transaction.Nonce)
    return [tx.Transaction.Data for tx in txs]


def Insert(owner: int, data: str) -> int:
    if owner not in AvailableOwners:
        raise ValueError(f"{owner=} isn't in {AvailableOwners=}")

    nonce = len(Get(owner))

    while True:
        _Logger.info(f"尝试插入数据到 {owner=} 的索引 {nonce} 处...")

        ok = _Client().Insert(
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
            _Logger.info("该位置已被其它客户端插入数据, 即将重试...")
            return Insert(owner, data)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
