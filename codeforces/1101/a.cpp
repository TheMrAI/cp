#include <bits/stdc++.h>
#include <cstdint>
#include <functional>

using namespace std;
#define bl '\n';

// Source: https://codeforces.com/blog/entry/78852
auto ceil2(int64_t a, int64_t b) -> int64_t { return (a + b - 1) / b; }

auto main() -> int
{
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);

  auto t = int64_t{ 0 };
  std::cin >> t;
  while (t--) {
    int64_t n;
    cin >> n;
    // auto counts = map<int64_t, int64_t, greater<int64_t>>{};
    auto nums = vector<int64_t>{};
    for (auto i = 0; i < n; ++i) {
      int64_t num;
      cin >> num;
      nums.push_back(num);
      // ++counts[num];
    }
    int64_t count = 0;
    sort(nums.begin(), nums.end());
    int64_t l = 0;
    int64_t r = nums.size() - 1;
    while (l <= r && (nums[l] != nums[r])) {
      ++l;
      --r;
      ++count;
    }
    cout << count << bl;
  }

  return 0;
}
