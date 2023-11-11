#include <algorithm>
#include <iostream>
#include <vector>

// A - Rock
// B - Paper
// C - Scissors

// X - Rock
// Y - Paper
// Z - Scissors

// X - Lose
// Y - Draw
// Z - Win

auto get_score(char opponent, char you) -> uint64_t
{
  auto opponent_value = static_cast<uint64_t>(opponent - 'A');
  auto you_value = static_cast<uint64_t>(you - 'X');

  auto score = you_value + 1;
  auto diff = static_cast<int64_t>(you_value) - static_cast<int64_t>(opponent_value);
  if (diff == 0) {
    score += 3;
  } else if (diff == 1 || diff == -2) {
    score += 6;
  }

  return score;
}

auto get_updated_score(char opponent, char result) -> uint64_t
{
  auto opponent_value = static_cast<uint64_t>(opponent - 'A') + 1;

  auto score = 0;
  if (result == 'Y') {
    score += 3 + opponent_value;
  } else if (result == 'Z') {
    score += 6;
    if (opponent_value == 3) {
      score += 1;
    } else {
      score += opponent_value + 1;
    }
  } else {
    if (opponent_value == 1) {
      score += 3;
    } else {
      score += opponent_value - 1;
    }
  }

  return score;
}

auto main() -> int
{
  auto total_score = uint64_t{ 0 };
  auto updated_score = uint64_t{ 0 };

  auto input = std::string{};
  while (std::getline(std::cin, input)) {
    total_score += get_score(input[0], input[2]);
    updated_score += get_updated_score(input[0], input[2]);
  }
  std::cout << "Total score: " << total_score << std::endl;
  std::cout << "Updated score: " << updated_score << std::endl;

  return 0;
}