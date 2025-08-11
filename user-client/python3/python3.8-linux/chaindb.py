import datetime
import typing
import random
import tempfile
import os
import time
import platform
import subprocess
import requests


class chaindb_config:
    DiscoveryInterval = datetime.timedelta(seconds=10)
    TCPPortUserService = 56784


class blockchain:
    class Transaction(typing.NamedTuple):
        OwnerID: int
        Nonce: int
        Data: str


class UserClient:
    def __init__(self, origin=f"http://localhost:{chaindb_config.TCPPortUserService}"):
        self.origin = origin  # const

        if self.origin.startswith("http://localhost:"):
            if platform.uname().system == "Linux" and platform.uname().machine in {
                "AMD64",
                "x86_64",
            }:
                if os.system("bash -c 'PATH+=: type -P chaindb.x64-linux.exe'") == 0:
                    os.system(
                        """ bash -c "PATH+=: bash -c 'chaindb.x64-linux.exe 2>./chaindb.log.txt &'" """
                    )
                else:
                    with requests.get(
                        "https://github.com/shynur/chaindb/releases/latest/download/chaindb.x64-linux.exe",
                        stream=True,
                    ) as resp:
                        if resp.ok:
                            rand_temp_exe_path = os.path.join(
                                tempfile.gettempdir(),
                                f"chaindb.x64-linux.__{random.randint(0, 2**32)}__.exe",
                            )
                            with open(
                                rand_temp_exe_path,
                                "xb",
                            ) as fchaindb:
                                for chunk in resp.iter_content(chunk_size=8192):
                                    fchaindb.write(chunk)
                            os.system(f"chmod a+x {rand_temp_exe_path}")
                            os.system(
                                f"bash -c '{rand_temp_exe_path} 2>./chaindb.log.txt &'"
                            )
            else:
                raise RuntimeError("当前平台不受支持")
            time.sleep(chaindb_config.DiscoveryInterval.total_seconds())

    def ListOwners(
        self,
        *,
        confirmations=0,
    ):
        resp = requests.get(f"{self.origin}/api/v1/owners")
        if resp.status_code not in range(200, 300):
            raise RuntimeError(
                f"请求失败: {resp.status_code} {resp.reason} {resp.text}"
            )

        class Owner(typing.NamedTuple):
            OwnerID: int
            ConfirmationScore: int

        owners = []
        for owner in resp.json() or []:
            owners.append(
                Owner(
                    OwnerID=owner["OwnerID"],
                    ConfirmationScore=owner["ConfirmationScore"],
                )
            )

        return [owner for owner in owners if owner.ConfirmationScore >= confirmations]

    def ListTransactionsOwnedBy(
        self,
        owner,
        *,
        confirmations=0,
    ):
        resp = requests.get(f"{self.origin}/api/v1/transactions?owner={owner}")
        if resp.status_code not in range(200, 300):
            raise RuntimeError(
                f"请求失败: {resp.status_code} {resp.reason} {resp.text}"
            )

        class Tx(typing.NamedTuple):
            Transaction: blockchain.Transaction
            ConfirmationScore: int

        txs = []
        for tx in resp.json() or []:
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
        transaction,
        *,
        confirmations,
    ):
        resp = requests.post(
            f"{self.origin}/api/v1/transactions?confirmation={confirmations}",
            json={
                "OwnerID": transaction.OwnerID,
                "Nonce": transaction.Nonce,
                "Data": transaction.Data,
            },
        )
        return resp.status_code in range(200, 300)
