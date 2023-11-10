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

struct MonkeyNode {
    bool ready;
    int64_t value;
    std::string operation;
    std::string leftName;
    MonkeyNode* left;
    std::string rightName;
    MonkeyNode* right;
};

auto evaluate_node(MonkeyNode* node) -> int64_t {
    if (node->ready) {
        return node->value;
    }
    auto left_value = evaluate_node(node->left);
    auto right_value = evaluate_node(node->right);

    if (node->operation == "+") {
        node->value = left_value + right_value;
    } else if (node->operation == "-") {
        node->value = left_value - right_value;
    } else if (node->operation == "*") {
        node->value = left_value * right_value;
    } else if (node->operation == "/") {
        node->value = left_value / right_value;
    }
    node->ready = true;
    return node->value;
}

auto main() -> int {
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    auto line = std::string{};
    auto dependency_regex = std::regex(R"((\w+): (\w+) (.{1}) (\w+))");
    auto leaf_regex = std::regex(R"((\w+): (\d+))");

    auto node_map = std::unordered_map<std::string, std::unique_ptr<MonkeyNode>>{};

    while(std::getline(std::cin, line)) { 
        auto match = std::smatch{};

        auto node = std::make_unique<MonkeyNode>();
        if (std::regex_match(line, match, dependency_regex)) {
            node->ready = false;
            node->value = 0;
            node->leftName = match[2];
            node->operation = match[3];
            node->left = nullptr;
            node->rightName = match[4];
            node->right = nullptr;
            node_map[match[1]] = std::move(node);
        } else if (std::regex_match(line, match, leaf_regex)) {
            node->ready = true;
            node->value = std::stoi(match[2]);
            node->leftName = "";
            node->operation = "";
            node->left = nullptr;
            node->rightName = "";
            node->right = nullptr;
            node_map[match[1]] = std::move(node);
        } else {
            std::cerr << "Oh nooo! Explosion imminent!" << std::endl;
            return -1;
        }
    }
    // connect the graph pointers, now that we have all entries
    for (auto& entry : node_map) {
        if (!entry.second->ready) {
            entry.second->left = node_map.at(entry.second->leftName).get();
            entry.second->right = node_map.at(entry.second->rightName).get();
        }
    }

    std::cout << "Root value: " << evaluate_node(node_map.at("root").get()) << std::endl;

    return 0;
}
