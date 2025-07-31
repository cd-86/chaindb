#pragma once
#include <bits/stdc++.h>
using namespace std::literals;

namespace shynur::chaindb {

    struct [[gnu::weak]] Transaction {
        /**
         * @brief 数据条目在 ChainDB 中的索引号.
         * @details 因为一个索引号对应的数据条目可能有很多个修订历史,
         *          即一对多的关系, 所以称索引号为 owner.
         * @warning 只能是自然数, 即 [0, 2**31).  其余范围供内部使用.
         */
        const std::int32_t OwnerID;

        /**
         * @brief 修订历史, 表示这是第几次修订.
         * @note 某个 owner 第一次插入数据, 该字段为 0.
         * @details
         */
        const unsigned Nonce;

        /**
         * @brief 数据条目的主体.
         */
        const std::string Data;
    };

    struct [[gnu::weak]] UserClient {
        const std::string server_host;
        const std::uint16_t server_port;

        UserClient(
            const std::string server_host, const std::uint16_t server_port
        ): server_host{server_host}, server_port{server_port} {}

        auto ListOwners() const {
            struct Owner {
                const std::int32_t OwnerID;
				const unsigned ConfirmationScore;
            };
            auto owners = std::vector<Owner>{};


            return owners;
        }

        auto ListTransactionsOwnedBy(const std::uint32_t owner) const {
            struct TxWithConfirmation {
                const Transaction Tx;
                const unsigned ConfirmationScore;
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
namespace rbk::chaindb {
    using UserClient = ::shynur::chaindb::UserClient;
}
#endif
