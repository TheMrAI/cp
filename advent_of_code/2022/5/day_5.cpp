#include <iostream>
#include <regex>
#include <stack>
#include <vector>

auto main() -> int
{
  auto input = std::string();
  auto stacks_init = std::vector<std::vector<char>>{};

  while (std::getline(std::cin, input)) {
    if (isdigit(input[1]) != 0) { break; }
    if (stacks_init.empty()) {
      auto column_count = 1 + ((input.size() - 3) / 4);
      for (auto i = size_t{ 0 }; i < column_count; ++i) { stacks_init.emplace_back(std::vector<char>{}); }
    }
    for (auto index = size_t{ 0 }; index < input.size(); index += 4) {
      stacks_init[(index / 4)].emplace_back(input[index + 1]);
    }
  }
  // empty line
  std::getline(std::cin, input);
  auto stacks = std::vector<std::stack<char>>{};
  for (auto i = size_t{ 0 }; i < stacks_init.size(); ++i) {
    stacks.emplace_back(std::stack<char>{});
    for (auto element = stacks_init[i].rbegin(); element != stacks_init[i].rend(); ++element) {
      if (*element == ' ') { break; }
      stacks[i].push(*element);
    }
  }

  auto regexp = std::regex(R"(move (\d+) from (\d+) to (\d+))");
  while (std::getline(std::cin, input)) {
    std::cout << input << std::endl;
    auto matches = std::smatch{};
    std::regex_match(input, matches, regexp);
    auto count = std::stoul(matches[1].str());
    auto from = std::stoul(matches[2].str()) - 1;
    auto to = std::stoul(matches[3].str()) - 1;

    // while(count != 0) {
    //     auto element = stacks[from].top();
    //     stacks[from].pop();
    //     stacks[to].push(element);
    //     --count;
    // }
    std::vector<char> temp;
    while (count != 0) {
      auto element = stacks[from].top();
      stacks[from].pop();
      temp.push_back(element);
      --count;
    }
    for (auto element = temp.rbegin(); element != temp.rend(); ++element) { stacks[to].push(*element); }
  }

  auto result = std::string{};
  for (auto pile : stacks) { result += pile.top(); }
  std::cout << "Result: " << result << std::endl;

  return 0;
}