#include <algorithm>
#include <iostream>
#include <vector>

auto main() -> int
{
  auto input = std::string();
  auto calorie_totals = std::vector<uint64_t>{};
  auto calorie_counter = uint64_t{ 0 };
  while (std::getline(std::cin, input)) {
    if (input.empty()) {
      if (calorie_counter != 0) { calorie_totals.push_back(calorie_counter); }
      calorie_counter = 0;
      continue;
    }
    calorie_counter += std::stoull(input);
  }
  std::sort(calorie_totals.begin(), calorie_totals.end(), std::greater<uint64_t>());

  std::cout << "Calorie totals: " << std::endl;
  for (auto const total : calorie_totals) { std::cout << total << std::endl; }

  return 0;
}