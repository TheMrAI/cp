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
  vector<vector<int>> threeSum(vector<int>& nums)
  {
    std::sort(nums.begin(), nums.end());
    auto triplets = std::vector<std::vector<int>>{};

    for (size_t i = 0; i < nums.size() - 2; ++i) {
      if (i > 0 && nums[i] == nums[i - 1]) { continue; }
      auto front = size_t{ i + 1 };
      auto back = nums.size() - 1;

      while (front < back) {
        auto sum = nums[i] + nums[front] + nums[back];
        if (sum == 0) {
          while (nums[front] == nums[front + 1] && front + 1 < back) { ++front; }
          while (nums[back] == nums[back - 1] && back - 1 > front) { --back; }
          triplets.emplace_back(std::vector<int>{ nums[i], nums[front], nums[back] });
          ++front;
          --back;
        } else if (sum < 0) {
          ++front;
        } else {
          --back;
        }
      }
    }

    return triplets;
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
