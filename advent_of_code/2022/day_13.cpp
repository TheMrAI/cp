#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

// Nasty but oh well it works
struct Item {
    bool is_number;
    int number;
    vector<Item> items;
};

auto parse(string const& input, regex const& tokens_regex) -> Item {
    auto tokens_begin = 
        std::sregex_iterator(input.begin(), input.end(), tokens_regex);
    auto tokens_end = std::sregex_iterator();

    stack<Item> sections;
    auto result = Item{false, 0, {}};
    for(sregex_iterator i = tokens_begin; i != tokens_end; ++i) {
        string match_str = i->str();
        if(match_str == string{"["}) {
            sections.push({false, 0, {}});
        } else if (match_str == string{"]"}) {
            result = sections.top();
            sections.pop();
            if(!sections.empty()) {
                auto& parent = sections.top();
                parent.items.push_back(result);
            }
        } else {
            auto& section = sections.top();
            section.items.push_back(Item{true, stoi(match_str), {}});
        }
    }

    return result;
}

// just checking this abomination
auto print(Item const& item) -> void {
    if (item.is_number) {
        cout << item.number;
    } else {
        cout << "[";
        for(auto i = 0; i < item.items.size(); ++i) {
            print(item.items[i]);
            if (i < item.items.size() - 1) {
                cout<< ",";
            }
        }
        cout << "]";
    }
}

// -1 -smaller, 0 - equal, 1 - bigger
auto smaller(Item const& left, Item const& right) -> int {
    if (left.is_number && right.is_number) {
        if(left.number < right.number) {
            return -1;
        }
        if(left.number == right.number) {
            return 0;
        }
        return 1;
    }
    if (!left.is_number && !right.is_number) {
        for(auto i=0; i < left.items.size() && i < right.items.size(); ++i) {
            auto result = smaller(left.items[i], right.items[i]);
            if(result != 0) {
                return result;
            }
        }
        if (left.items.size() < right.items.size()) {
            return -1;
        }
        if (left.items.size() == right.items.size()) {
            return 0;
        }
        return 1;
    }
    if (left.is_number && !right.is_number) {
        auto left_listified = Item{false, 0, {Item{true, left.number, {}}}};
        return smaller(left_listified, right);
    }
    if (!left.is_number && right.is_number) {
        auto right_listified = Item{false, 0, {Item{true, right.number, {}}}};
        return smaller(left, right_listified);
    }
    return -1;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto left = string{};
    auto right = string{};
    auto data = string{};
    auto index = 1;
    auto correct_index_sum = 0;
    auto tokens = regex{R"((\[|\]|\d+))"};
    
    auto packets = vector<Item>{};
    while(getline(cin, data)) {
        if(data.empty()) {
            continue;
        }
        if (left.empty()) {
            left = data;
            continue;
        }
        if (right.empty()) {
            right = data;
        }

        auto left_data = parse(left, tokens);
        left.erase();
        auto right_data = parse(right, tokens);
        right.erase();

        auto comparison = smaller(left_data, right_data);
        packets.push_back(move(left_data));
        packets.push_back(move(right_data));
        if(comparison == -1) {
            cout << index << endl;
            correct_index_sum += index;
        }
        ++index;
    }
    cout << "Correct index sum: " << correct_index_sum << endl;

    packets.push_back(Item{false, 0, {Item{false, 0, {Item{true, 2, {}}}}}});
    packets.push_back(Item{false, 0, {Item{false, 0, {Item{true, 6, {}}}}}});
    sort(packets.begin(), packets.end(), 
    [](auto const& lhs, auto const& rhs) -> bool {
        return smaller(lhs, rhs) == -1;
    });
    for(auto i = 0; i < packets.size(); ++i) {
        cout << i + 1 << " : ";
        print(packets[i]);
        cout << endl;
    }

    return 0;
}
