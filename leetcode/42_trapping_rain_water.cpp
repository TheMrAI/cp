#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

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

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

class Solution
{
public:
  int trap(vector<int>& heights)
  {
    auto left_maximums = std::vector<int>(heights.size());
    auto left_max = 0;
    for (auto i = size_t{ 0 }; i < heights.size(); ++i) {
      left_maximums[i] = left_max;
      left_max = std::max(left_max, heights[i]);
    }

    auto right_maximums = std::vector<int>(heights.size());
    auto right_max = 0;
    for (auto i = static_cast<int>(heights.size() - 1); i >= 0; --i) {
      right_maximums[i] = right_max;
      right_max = std::max(right_max, heights[i]);
    }

    auto total_volume = 0;
    for (auto i = size_t{ 0 }; i < heights.size(); ++i) {
      auto value = std::min(left_maximums[i], right_maximums[i]) - heights[i];
      total_volume += std::max(value, 0);
    }

    return total_volume;
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
