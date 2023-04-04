#include <iostream>
#include <vector>
#include <algorithm>
#include <unordered_set>

auto to_priority(char input) -> uint {
    if (std::islower(input)) {
        return 1 + input - 'a';
    }
    return 27 + input - 'A';
}

auto find_duplicate(std::string input) -> char {
    auto seen = std::vector<bool>(52, false);
    auto half_point = input.size() / 2;
    auto index = size_t{0};

    while (index < half_point) {
        seen[to_priority(input[index])-1] = true;
        ++index;
    }
    while (index < input.size()) {
        auto seen_index = to_priority(input[index])-1;
        if (seen[seen_index]) {
            return input[index];
        }
        ++index;
    }
    return '!';
} 

auto main() -> int {
    auto input = std::string();

    auto total_priority = uint64_t{0};
    auto total_badge_priority = uint64_t{0};
    auto triplets = std::vector<std::string>{}; 
    while (std::getline(std::cin, input)) {
        auto duplicate = find_duplicate(input);
        if (duplicate == '!') {
            std::cerr << "No duplicate item!" << std::endl;
            return 1;
        }
        total_priority += to_priority(duplicate);

        triplets.push_back(input);
        if (triplets.size() == 3) {
            auto first = std::vector<char>(triplets[0].cbegin(), triplets[0].cend());
            std::sort(first.begin(), first.end());
            auto second = std::vector<char>(triplets[1].cbegin(), triplets[1].cend());
            std::sort(second.begin(), second.end());
            auto third = std::vector<char>(triplets[2].cbegin(), triplets[2].cend());
            std::sort(third.begin(), third.end());

            auto intersection = std::vector<char>{};
            std::set_intersection(first.cbegin(), first.cend(), second.cbegin(), second.cend(), std::back_inserter(intersection));
            auto final_group = std::vector<char>{};
            std::set_intersection(intersection.cbegin(), intersection.cend(), third.cbegin(), third.cend(), std::back_inserter(final_group));
            
            if (final_group.size() != 1) {
                std::cerr << "Intersection does not contain exactly 1 item!" << std::endl;
                std::cerr << triplets[0] << std::endl;
                std::cerr << triplets[1] << std::endl;
                std::cerr << triplets[2] << std::endl;
                for (auto const c : final_group) {
                    std::cerr << c << std::endl;
                }
            }
            total_badge_priority += to_priority(*final_group.begin());
            
            triplets.clear();
        }
    }

    std::cout << "Total priority: " << total_priority << std::endl;
    std::cout << "Total badge priority: " << total_badge_priority << std::endl;

    return 0;
}