#define SHYNUR_USED_BY_SEER_ROBOTICS_RBK
#include "include/chaindb.hpp"
namespace chaindb = rbk::chaindb;

int main() {
    const auto uc = chaindb::UserClient{};

    for (char op; std::cin; std::cout << std::endl) {
        std::cout << "Operation (o=列出索引, t=列出数据, i=插入数据): ";
        std::cin >> op;
        try {
            switch (op) {
                case 'o': {
                    std::cout << "Confirmation Score >= ";
                    unsigned confirmation;
                    std::cin >> confirmation;

                    for (const auto& [owner, confirmation] : uc.ListOwners(confirmation))
                        std::cout << "OwnerID: " << owner << '\t'
                                  << "ConfirmationScore: " << confirmation << '\n';
                }
                    break;
                case 't': {
                    std::cout << "OwnerID=";
                    unsigned owner;
                    std::cin >> owner;

                    std::cout << "Confirmation Score >= ";
                    unsigned confirmation;
                    std::cin >> confirmation;

                    for (const auto& [transaction, confirmation] : uc.ListTransactionsOwnedBy(owner, confirmation)) {
                        std::cout << "OwnerID: " << transaction.OwnerID << '\t'
                                  << "Nonce: " << transaction.Nonce << '\t'
                                  << "ConfirmationScore: " << confirmation;
                        if (!transaction.Data.empty())
                            std::cout << '\t' << "Data: " << transaction.Data;
                        std::cout << '\n';
                    }
                }
                    break;
                case 'i': {
                    std::cout << "OwnerID=";
                    std::uint32_t owner_id;
                    std::cin >> owner_id;

                    std::cout << "Nonce=";
                    unsigned nonce;
                    std::cin >> nonce;

                    std::cout << "Data=";
                    std::string data;
                    std::cin >> data;

                    std::cout << "Confirmation Score >= ";
                    unsigned confirmation;
                    std::cin >> confirmation;

                    const auto ok = uc.Insert({owner_id, nonce, data}, confirmation);
                    if (ok)
                        std::cout << "Insert OK";
                    else
                        std::cout << "Insert Failed";
                    std::cout << '\n';
                }
            }
        } catch (const std::exception& e) {
            std::cerr << "Error: " << e.what() << std::endl;
        }
    }
}
