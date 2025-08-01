#pragma once
#include <bits/stdc++.h>
#include "nlohmann/json/single_include/nlohmann/json.hpp"
#include "yhirose/cpp-httplib/httplib.h"
using namespace std::literals;

namespace shynur::chaindb {
    namespace chaindb_config {
        constexpr auto DiscoveryInterval = 10s;
        constexpr auto TCPPortUserService = 56784;
    }

    namespace blockchain {
        /**
         * @brief ChainDB 存储的基本元素.  代表一次修订历史.
         */
        struct [[gnu::weak]] Transaction {
            /**
             * @brief 数据条目在 ChainDB 中的索引号.
             * @details 因为一个索引号对应的数据条目可能有很多个修订历史,
             *          即一对多的关系, 所以称索引为 owner.
             */
            const std::uint32_t OwnerID;

            /**
             * @brief 修订历史, 表示这是第几次修订.
             * @note 某个 owner 第一次插入数据时, 请设置该字段为 0.
             * @details ChainDB 是一个巨大的从 OwnerID 映射到 Data 的哈希表,
             *          多个 client 可能同时发起对 Data 的修订请求.  为了保证
             *          操作的原子性, Nonce 被 ChainDB 视为原子变量, Data 的插入
             *          成功与否与 Nonce 绑定.  借鉴自 Ethereum (以太坊).
             *          Nonce 采用 Compare-And-Swap 方式进行更新 (下列操作作为一个整体是原子的):
             *              - ChainDB 找出相同 OwnerID 的数据条目的修订历史记录
             *              - 发现最后一次修订号是 LastNonce
             *              - 如果 Nonce==LastNonce+1, 则插入成功, 且最后一次修订号被更新为了 Nonce
             * @warning Nonce 通常是连续的.  如果你插入一个跳跃很大的 Nonce, 它可能
             *          会驻留在 ChainDB 的待插入队列中足够久的时间, 直到某一时刻被
             *          成功插入, 而你并不知情.
             */
            const unsigned Nonce;

            /**
             * @brief 数据条目的主体.
             */
            const std::string Data;
        };
    }

    /**
     * @brief 向固定的 ChainDB 服务器发送请求的客户端.
     * @note 所有方法都是阻塞的.
     *       服务器正在同步数据时, 会阻塞较久.
     */
    struct [[gnu::weak]] UserClient {
        const std::string origin;

        /**
         * @note 如果 server_host 是 "localhost", 则会尝试启动一个本地的
         *       ChainDB 节点, 并等待它启动.  即使等待了, 也不保证数据同步已经
         *       完成.  如果数据未同步, 后续的请求会继续阻塞, 无需关心.
         */
        UserClient(
            const std::string origin
                = "http://localhost:" + std::to_string(chaindb_config::TCPPortUserService)
        ): origin{origin} {
            if (this->origin.find("http://localhost:") == 0) {
                std::system(
                    R"(bash -c "cd ../..; make chaindb.exe; ./chaindb.exe &>/dev/null &")"
                );
                std::this_thread::sleep_for(chaindb_config::DiscoveryInterval);
            }
        }

        /**
         * @brief 列出所有已知的 owners.
         *        OwnerID 按照它们第一次出现在 ChainDB 中的顺序排列.
         * @param required_confirmation_score
         *        返回的结果中最小的 ConfirmationScore 值.  通常设为 10 即可.
         *        表示 ChainDB 中索引号 OwnerID 存在的概率.
         * @details 根据区块链理论,
         *              p^ConfirmationScore 与 数据可能被篡改的概率
         *          成正比, 其中 p 取值 [0, 1).
         *          新插入的 transaction 大约每 1.x seconds 增加一个 ConfirmationScore.
         * @see blockchain::Transaction
         */
        auto ListOwners(const unsigned required_confirmation_score = 0) const {
            const auto owners_json = [this] {
                const auto resp = httplib::Client{this->origin}.Get("/api/v1/owners");

                if (!resp)
                    throw std::runtime_error{"[shynur/chaindb] HTTP failed"};

                if (resp->status / 100 != 2)
                    throw std::runtime_error{
                        "[shynur/chaindb] HTTP status=" + std::to_string(resp->status) + ' '
                        + "error: " + resp->reason
                    };

                return resp->body;
            }();

            struct Owner {
                const std::uint32_t OwnerID;
                const unsigned ConfirmationScore;
            };
            auto owners = std::vector<Owner>{};
            for (const auto& owner : ::nlohmann::json::parse(owners_json)) {
                owners.push_back({
                    .OwnerID = owner["OwnerID"],
                    .ConfirmationScore = owner["ConfirmationScore"],
                });
            }

            return [&] {
                auto confirmed_owners = std::vector<Owner>{};
                std::copy_if(
                    std::cbegin(owners), std::cend(owners),
                    std::back_inserter(confirmed_owners),
                    [&](const auto& owner) {
                        return owner.ConfirmationScore > required_confirmation_score;
                    }
                );
                return confirmed_owners;
            }();
        }

        /**
         * @param required_confirmation_score
         *        返回的结果中最小的 ConfirmationScore 值.  通常设为 10 即可.
         *        表示 ChainDB 中某个 transaction 存在的概率.
         * @see blockchain::Transaction
         * @details 根据区块链理论,
         *              p^ConfirmationScore 与 数据可能被篡改的概率
         *          成正比, 其中 p 取值 [0, 1).
         *          新插入的 transaction 大约每 1.x seconds 增加一个 ConfirmationScore.
         */
        auto ListTransactionsOwnedBy(
            const std::uint32_t owner,
            const unsigned required_confirmation_score = 0
        ) const {
            const auto txs_json = [&, this] {
                const auto resp = httplib::Client{this->origin}.Get(
                    "/api/v1/transactions?owner=" + std::to_string(owner)
                );

                if (!resp)
                    throw std::runtime_error{"[shynur/chaindb] HTTP failed"};

                if (resp->status / 100 != 2)
                    throw std::runtime_error{
                        "[shynur/chaindb] HTTP status=" + std::to_string(resp->status) + ' '
                        + "error: " + resp->reason
                    };

                return resp->body;
            }();

            struct Tx {
                const blockchain::Transaction Transaction;
                const unsigned ConfirmationScore;
            };
            auto txs = std::vector<Tx>{};
            for (const auto& tx : ::nlohmann::json::parse(txs_json)) {
                txs.push_back({
                    .Transaction = {
                        .OwnerID = tx["Tx"]["OwnerID"],
                        .Nonce = tx["Tx"]["Nonce"],
                        .Data = [&]() -> std::string {
                            try {
                                return tx["Tx"].at("Data");
                            } catch (const ::nlohmann::json::out_of_range&) {
                                return "";
                            }
                        }(),
                    },
                    .ConfirmationScore = tx["ConfirmationScore"],
                });
            }

            return [&]() {
                auto confirmed_txs = std::vector<Tx>{};
                std::copy_if(
                    std::cbegin(txs), std::cend(txs),
                    std::back_inserter(confirmed_txs),
                    [&](const auto& tx) {
                        return tx.ConfirmationScore > required_confirmation_score;
                    }
                );
                return confirmed_txs;
            }();
        }

        /**
         * @brief 向 ChainDB 插入一条 transaction.
         * @param required_confirmation_score
         *        等待, 直到被插入的这条 transaction 在 ChainDB 中
         *        的 ConfirmationScore 达到该值.  通常设为 10 即可.
         * @details 根据区块链理论,
         *              p^ConfirmationScore 与 数据可能被篡改的概率
         *          成正比, 其中 p 取值 [0, 1).
         *          新插入的 transaction 大约每 1.x seconds 增加一个 ConfirmationScore.
         * @note 如果因 CAS 操作失败导致插入失败 (see blockchain::Transaction::Nonce),
         *       返回 false.
         *       如果插入成功, 返回 true.
         *       如果不清楚结果如何, 但是超时了, 也会返回 false, 你可以重试.
         */
        auto Insert(
            const blockchain::Transaction transaction,
            const std::uint32_t required_confirmation_score
        ) const {
            const auto tx_json = ::nlohmann::json{
                {"OwnerID", transaction.OwnerID},
                {"Nonce", transaction.Nonce},
                {"Data", transaction.Data},
            }.dump(4);

            const auto resp = httplib::Client{this->origin}.Post(
                "/api/v1/transactions?confirmation=" + std::to_string(required_confirmation_score),
                tx_json,
                "application/json"
            );
            if (!resp)
                throw std::runtime_error{"[shynur/chaindb] HTTP failed"};
            return resp->status / 100 == 2;
        }
    };
}

#ifdef SHYNUR_USED_BY_SEER_ROBOTICS_RBK
namespace rbk::chaindb {
    using UserClient = ::shynur::chaindb::UserClient;
    using Transaction = ::shynur::chaindb::blockchain::Transaction;
}
#endif
