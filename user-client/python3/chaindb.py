#! python3.13 -m

import datetime
import typing
import os
import time
import json
import sys
import requests


class chaindb_config:
    DiscoveryInterval: typing.ReadOnly[datetime.timedelta] = datetime.timedelta(
        seconds=10
    )
    TCPPortUserService: typing.ReadOnly[int] = 56784


class blockchain:
    class Transaction(typing.NamedTuple):
        OwnerID: int
        Nonce: int
        Data: str

    class UserClient:
        def __init__(
            self, origin: str = f"http://localhost:{chaindb_config.TCPPortUserService}"
        ):
            self.origin: str = origin  # const

            if self.origin.startswith("http://localhost:"):
                match os.name:
                    case "nt":
                        # os.system("start chaindb.x64-mswindows.exe")
                        ...
                    case "posix":
                        os.system("chaindb.x64-linux.exe &")
                    case _:
                        raise RuntimeError("当前平台不受支持")
                time.sleep(chaindb_config.DiscoveryInterval.total_seconds())

        def ListOwners(
            self,
            *,
            confirmations: int = 0,
        ) -> list:
            resp = requests.get(f"{self.origin}/api/v1/owners")
            if resp.status_code not in range(200, 300):
                raise RuntimeError(
                    f"请求失败: {resp.status_code} {resp.reason} {resp.text}"
                )

            class Owner(typing.NamedTuple):
                OwnerID: int
                ConfirmationScore: int

            owners: list[Owner] = []
            for owner in resp.json() or []:
                owners.append(
                    Owner(
                        OwnerID=owner["OwnerID"],
                        ConfirmationScore=owner["ConfirmationScore"],
                    )
                )

            return [
                owner for owner in owners if owner.ConfirmationScore >= confirmations
            ]

        def ListTransactionsOwnedBy(
            self,
            owner: int,
            *,
            confirmations: int = 0,
        ) -> list:
            resp = requests.get(f"{self.origin}/api/v1/transactions?owner={owner}")
            if resp.status_code not in range(200, 300):
                raise RuntimeError(
                    f"请求失败: {resp.status_code} {resp.reason} {resp.text}"
                )

            class Tx(typing.NamedTuple):
                Transaction: blockchain.Transaction
                ConfirmationScore: int

            txs: list[Tx] = []
            for tx in json.loads(resp.json()) or []:
                txs.append(
                    Tx(
                        Transaction=blockchain.Transaction(
                            OwnerID=tx["Tx"]["OwnerID"],
                            Nonce=tx["Tx"]["Nonce"],
                            Data=tx["Tx"].get("Data", ""),
                        ),
                        ConfirmationScore=tx["ConfirmationScore"],
                    )
                )
            return [tx for tx in txs if tx.ConfirmationScore >= confirmations]

        def Insert(
            self,
            transaction: "blockchain.Transaction",
            *,
            confirmations: int,
        ) -> bool:
            resp = requests.post(
                f"{self.origin}/api/v1/transactions?confirmation={confirmations}",
                json={
                    "OwnerID": transaction.OwnerID,
                    "Nonce": transaction.Nonce,
                    "Data": transaction.Data,
                },
            )
            return resp.status_code in range(200, 300)


if __name__ == "__main__":
    uc = blockchain.UserClient()

    print()

    while True:
        try:
            op = input("Operation (o=列出索引, t=列出数据, i=插入数据): ").strip()
            match op:
                case "o":
                    confirmations = int(
                        input("(默认是 0) Confirmation Score >= ") or "0"
                    )
                    for owner, confirmation in uc.ListOwners(
                        confirmations=confirmations
                    ):
                        print(
                            f"OwnerID: {owner.OwnerID}\tConfirmationScore: {confirmation}"
                        )
                case "t":
                    owner = int(input("OwnerID="))
                    confirmations = int(
                        input("(默认是 0) Confirmation Score >= ") or "0"
                    )
                    for transaction, confirmation in uc.ListTransactionsOwnedBy(
                        owner, confirmations=confirmations
                    ):
                        print(
                            f"OwnerID: {transaction.OwnerID}\tNonce: {transaction.Nonce}\tConfirmationScore: {confirmation}{transaction.Data and f'\tData: {transaction.Data}'}"
                        )
                case "i":
                    owner = int(input("OwnerID="))
                    nonce = int(input("Nonce="))
                    confirmations = int(
                        input("(默认是 0) Confirmation Score >= ") or "0"
                    )
                    data = input("Data: ")
                    ok = uc.Insert(
                        blockchain.Transaction(OwnerID=owner, Nonce=nonce, Data=data),
                        confirmations=confirmations,
                    )
                    print()
                    print(f"\t{'OK' if ok else 'Failed'}")
        except Exception as e:
            print("Error: ", e, file=sys.stderr)
