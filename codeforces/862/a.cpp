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
    auto n = int{ 0 };
    cin >> n;
    auto numbers = vector<int>{};

    for (auto i = 0; i < n; ++i) {
      auto val = int{ 0 };
      cin >> val;
      numbers.push_back(val);
    }

    auto x = 0;
    auto total_xor = 0;
    for (; x < pow(2, 8); ++x) {
      total_xor = numbers[0] ^ x;
      for (auto i = 1; i < n; ++i) { total_xor = numbers[i] ^ x ^ total_xor; }
      if (total_xor == 0) { break; }
    }
    if (total_xor == 0) {
      cout << x << endl;
    } else {
      cout << -1 << endl;
    }
  }

  return 0;
}
