#define SHYNUR_USED_BY_SEER_ROBOTICS_RBK
#include "include/chaindb.hpp"
namespace chaindb = rbk::chaindb;

int main() {
    const auto uc = chaindb::UserClient{};

    for (const auto& [owner, confirmation] : uc.ListOwners()) {
        std::cout << "OwnerID: " << owner << '\t'
                  << "ConfirmationScore: " << confirmation << '\n';
    }

    for (const auto& [transaction, confirmation] : uc.ListTransactionsOwnedBy(421)) {
        std::cout << "OwnerID: " << transaction.OwnerID << '\t'
                  << "Nonce: " << transaction.Nonce << '\t'
                  << "Data: " << transaction.Data << '\t'
                  << "ConfirmationScore: " << confirmation << '\n';
    }
}
