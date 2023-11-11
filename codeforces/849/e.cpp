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
    auto values = vector<i6>(n, 0);
    auto negative_even_count = true;
    auto has_zero = false;
    auto sum = i6{ 0 };
    for (auto i = 0; i < n; ++i) {
      cin >> values[i];
      if (values[i] < 0) { negative_even_count = not negative_even_count; }
      if (values[i] == 0) { has_zero = true; }
      values[i] = abs(values[i]);
      sum += values[i];
    }
    if (has_zero || negative_even_count) {
      cout << sum << endl;
    } else {
      ranges::sort(values);
      cout << sum - 2 * values[0] << endl;
    }
  }

  return 0;
}
