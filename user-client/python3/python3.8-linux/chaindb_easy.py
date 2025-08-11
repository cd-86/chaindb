import chaindb

RequiredConfirmations: int = 10
AvailableOwners = range(100_0000, 200_0000)

uc = chaindb.UserClient()


def Owners() -> list:
    owners = uc.ListOwners(confirmations=RequiredConfirmations)
    return [owner.OwnerID for owner in owners if owner.OwnerID in AvailableOwners]


def Get(owner: int) -> list:
    txs = uc.ListTransactionsOwnedBy(owner, confirmations=RequiredConfirmations)
    txs.sort(key=lambda tx: tx.Transaction.Nonce)
    return [tx.Transaction.Data for tx in txs]


def Insert(owner: int, data: str) -> int:
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
