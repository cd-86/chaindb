#define SHYNUR_USED_BY_SEER_ROBOTICS_RBK
#include "include/chaindb.hpp"
namespace chaindb = rbk::chaindb;

int main() {
    const auto uc = chaindb::UserClient{};

    std::cout << std::endl;

    for (char op; std::cin; std::cout << std::endl) {
        std::cout << "Operation (o=列出索引, t=列出数据, i=插入数据): ";
        std::cin >> op;
        try {
            switch (op) {
                case 'o': {
                    std::cout << "(默认是 0) Confirmation Score >= ";
                    const auto confirmation = [] -> unsigned {
                        std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
                        auto confirmation = std::string{};
                        std::getline(std::cin, confirmation);
                        if (confirmation.empty())
                            return 0ul;
                        return std::stoul(confirmation);
                    }();

                    for (const auto& [owner, confirmation] : uc.ListOwners(confirmation))
                        std::cout << "OwnerID: " << owner << '\t'
                                  << "ConfirmationScore: " << confirmation << '\n';
                }
                    break;
                case 't': {
                    std::cout << "OwnerID=";
                    unsigned owner;
                    std::cin >> owner;

                    std::cout << "(默认是 0) Confirmation Score >= ";
                    const auto confirmation = [] -> unsigned {
                        std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
                        auto confirmation = std::string{};
                        std::getline(std::cin, confirmation);
                        if (confirmation.empty())
                            return 0ul;
                        return std::stoul(confirmation);
                    }();

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

                    std::cout << "(默认是 0) Confirmation Score >= ";
                    const auto confirmation = [] -> unsigned {
                        std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
                        auto confirmation = std::string{};
                        std::getline(std::cin, confirmation);
                        if (confirmation.empty())
                            return 0;
                        return std::stoul(confirmation);
                    }();

                    std::cout << "Data: ";
                    const auto data = [] {
                        auto input_line = std::string{};
                        std::getline(std::cin, input_line);
                        return input_line;
                    }();

                    const auto ok = uc.Insert({owner_id, nonce, data}, confirmation);
                    std::cout << "\n\t";
                    if (ok)
                        std::cout << "OK";
                    else
                        std::cout << "Failed";
                    std::cout << '\n';
                }
            }
        } catch (const std::exception& e) {
            std::cerr << "Error: " << e.what() << std::endl;
        }
    }
}
