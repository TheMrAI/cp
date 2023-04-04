#include <iostream>
#include <vector>
#include <algorithm>
#include <regex>

auto ranges_contained(uint64_t first_start, uint64_t first_end, uint64_t second_start, uint64_t second_end) -> bool {
    if (first_start <= second_start && first_end >= second_end) {
        return true;
    }
    if (second_start <= first_start && second_end >= first_end) {
        return true;
    }
    return false;
}

auto ranges_overlap(uint64_t first_start, uint64_t first_end, uint64_t second_start, uint64_t second_end) -> bool {
    if (first_end < second_start || first_start > second_end) {
        return false;
    }
    return true;
}

auto main() -> int {
    auto input = std::string();

    auto contained_pairs = uint64_t{0};
    auto overlapping_pairs = uint64_t{0};
    auto regexp = std::regex("(\\d+)-(\\d+),(\\d+)-(\\d+)");
    while (std::getline(std::cin, input)) {
        auto matches = std::smatch{};
        std::regex_match(input, matches, regexp);
        if (ranges_contained(std::stoull(matches[1].str()), std::stoull(matches[2].str()), std::stoull(matches[3].str()), std::stoull(matches[4].str()))) {
            ++contained_pairs;
        }

        if (ranges_overlap(std::stoull(matches[1].str()), std::stoull(matches[2].str()), std::stoull(matches[3].str()), std::stoull(matches[4].str()))) {
            ++overlapping_pairs;
        }
    }
    std::cout << "Contained pairs: " << contained_pairs << std::endl;
    std::cout << "Overlapping pairs: " << overlapping_pairs << std::endl;

    return 0;
}