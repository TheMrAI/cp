#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = 0;
    cin >> n;
    auto longest_space = 0;
    auto zero_count = 0;
    while (n--) {
      auto val = 0;
      cin >> val;
      if (val == 1) {
        zero_count = 0;
        continue;
      }
      ++zero_count;
      longest_space = max(longest_space, zero_count);
    }
    cout << longest_space << endl;
  }

  return 0;
}
