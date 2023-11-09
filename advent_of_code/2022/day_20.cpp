#include <bits/stdc++.h>

template<typename Output, typename Iter>
requires 
    std::derived_from<Output, std::ostream> &&
    std::weakly_incrementable<Iter> && 
    std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end) {
    while(begin != end) {
        output << *begin;
        if(std::next(begin) != end) {
            output << ' ';
        }
        ++begin;
    }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto move_element(std::list<int64_t>& ordering, std::unordered_map<size_t, std::list<int64_t>::iterator>& index_map, size_t index, bool move_right) -> void {
    if (!index_map.contains(index)) {
        std::cerr << "Unknown index given: " << index << std::endl;
        return;
    }
    auto element = index_map.at(index);
    auto value = *element;

    auto moves_to_make = std::abs(value) % (ordering.size() - 1);
    if (moves_to_make == 0) {
        return;
    }

    auto walker = ordering.erase(element);
    if (move_right && walker == ordering.end()) {
        walker = ordering.begin();
    }
    if (!move_right && walker == ordering.begin()) {
        walker = ordering.end();
    }

    auto offset = 1;
    if (!move_right) {
        offset = -1;
    }

    while (moves_to_make > 0) {
        std::advance(walker, offset);
        if (move_right && walker == ordering.end()) {
            walker = ordering.begin();
        }
        if (!move_right && walker == ordering.begin()) {
            walker = ordering.end();
        }
        --moves_to_make;
    }
    // reinsert element and update map
    element = ordering.insert(walker, value);
    index_map[index] = element;
}

auto solve(std::vector<int64_t> numbers, int64_t decryption_key, int64_t mixing_times) -> void {
    auto ordering = std::list<int64_t>{};
    auto index_map = std::unordered_map<size_t, std::list<int64_t>::iterator>{};
    
    for (auto i = size_t{0}; i < numbers.size(); ++i) {
        numbers[i] *= decryption_key;
        ordering.push_back(numbers[i]);
        index_map.emplace(i, std::prev(ordering.end()));
    }
    // dump_to(std::cout, numbers.begin(), numbers.end());
    // std::cout << std::endl;

    for (auto mixing_count = 0; mixing_count < mixing_times; ++mixing_count) {
        for (auto i = size_t{0}; i < numbers.size(); ++i) {
            auto move_right = true;
            if (numbers[i] < 0) {
                move_right = false;
            }
            move_element(ordering, index_map, i, move_right);
        }
    }

    numbers.clear();
    for (auto entry : ordering) {
        numbers.push_back(entry);
    }
    // dump_to(std::cout, numbers.begin(), numbers.end());
    // std::cout << std::endl;

    auto zero_at = std::find_if(numbers.begin(), numbers.end(), [](auto val){ return val == 0;});
    auto zero_index = std::distance(numbers.begin(), zero_at);
    auto thousandth = numbers[(zero_index + 1000) % numbers.size()];
    auto two_thousandth = numbers[(zero_index + 2000) % numbers.size()];
    auto three_thousandth = numbers[(zero_index + 3000) % numbers.size()];
    std::cout << thousandth << " " << two_thousandth << " " << three_thousandth << std::endl;
    std::cout << "Grove coordinates: " << thousandth + two_thousandth + three_thousandth << std::endl;
}

auto main() -> int {
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    auto numbers = std::vector<int64_t>{};

    auto number=int64_t{0};
    while(std::cin >> number) {
        numbers.push_back(number);
    }

    std::cout << "Part one: " << std::endl;
    solve(numbers, 1, 1);

    std::cout << "Part two: " << std::endl;
    solve(numbers, 811589153, 10);

    return 0;
}
