#include <bits/stdc++.h>

template<typename Output, typename Iter>
requires std::derived_from<Output, std::ostream> && std::weakly_incrementable<Iter> && std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end)
{
  while (begin != end) {
    output << *begin;
    if (std::next(begin) != end) { output << ' '; }
    ++begin;
  }
}

auto snafu_symbol_to_decimal(char c) -> int64_t
{
  switch (c) {
  case '2':
    return 2;
  case '1':
    return 1;
  case '0':
    return 0;
  case '-':
    return -1;
  case '=':
    return -2;
  }
  throw std::logic_error("Invalid symbol.");
}

auto decimal_digit_to_snafu_symbol(int64_t digit) -> char
{
  switch (digit) {
  case 2:
    return '2';
  case 1:
    return '1';
  case 0:
    return '0';
  case -1:
    return '-';
  case -2:
    return '=';
  }
  throw std::logic_error("Invalid digit.");
}

auto snafu_to_decimal(std::string const& data) -> int64_t
{
  auto result = int64_t{ 0 };
  for (auto i = 0; i < data.size(); ++i) {
    result += snafu_symbol_to_decimal(data[data.size() - 1 - i]) * static_cast<int64_t>(std::pow(double{ 5 }, i));
  }
  return result;
}

auto decimal_to_snafu(int64_t to_convert) -> std::string
{
  auto current = int64_t{ 1 };
  auto exponent = int64_t{ 1 };
  // find the maximum exponent we would need
  while (std::abs(current) < std::abs(to_convert)) {
    ++exponent;
    current = static_cast<int64_t>(std::pow(double{ 5 }, exponent));
  }

  constexpr auto digits = std::array<int64_t, 5>{ 2, 1, 0, -1, -2 };
  auto result = std::string{};
  current = int64_t{ 0 };
  for (; exponent >= 0; --exponent) {
    auto remaining = to_convert - current;
    auto value = static_cast<int64_t>(std::pow(double{ 5 }, exponent));
    auto max_correction = value / 2;
    auto digit_index = 0;
    for (; digit_index < digits.size(); ++digit_index) {
      auto added = (digits[digit_index] * value);
      auto up_corrected = added + max_correction;
      auto down_corrected = added - max_correction;
      if ((added >= remaining && down_corrected <= remaining) || (up_corrected >= remaining && added <= remaining)) {
        current += added;
        break;
      }
    }
    result += decimal_digit_to_snafu_symbol(digits[digit_index]);
    // std::cout << result << std::endl;
  }
  return result;
}

auto test_decimal_to_snafu(int64_t value) { std::cout << value << " => " << decimal_to_snafu(value) << std::endl; }

auto main() -> int
{
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);

  auto line = std::string{};
  auto total = int64_t{ 0 };
  while (std::getline(std::cin, line)) { total += snafu_to_decimal(line); }
  std::cout << "Total: " << total << ", in SNAFU: " << decimal_to_snafu(total) << std::endl;

  return 0;
}
