#pragma once
#include <bits/stdc++.h>
using namespace std::literals;

namespace ::shynur::chaindb {
    struct [[gnu::weak]] Transaction {
        const std::uint32_t OwnerID;
        const std::uint32_t Nonce;
        const std::string Data;
    };

    struct [[gnu::weak]] UserClient {
        const std::string server_host;
        const std::uint16_t server_port;

        UserClient(
            const std::string server_host, const std::uint16_t server_port,
        ): server_host{server_host}, server_port{server_port} {}

        auto ListOwners() const {
            struct Owner {
                const std::uint32_t OwnerID;
				const std::uint32_t ConfirmationScore;
            };
            auto owners = std::vector<owner>{};


            return owners;
        }

        auto ListTransactionsOwnedBy(const std::uint32_t owner) const {
            struct TxWithConfirmation {
                const Transaction Tx;
                const std::uint32_t ConfirmationScore;
            };
            auto txs = std::vector<TxWithConfirmation>{};

            return txs;
        }

        auto Insert(const Transaction transaction, const std::uint32_t confirmation_score) const
            -> std::future<bool> {

        }
    };
}

#ifdef SEER_ROBOTICS_RBK
namespace ::rbk::chaindb {
    using UserClient = ::shynur::chaindb::UserClient;
}
#endif
