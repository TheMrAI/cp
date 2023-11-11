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
  auto last_calculated = i6{ 1 };
  auto row_width = 1;
  auto next_row_from = 2;
  auto prev_row_start = 1;
  auto prev_row_end = 1;
  auto dp = vector<i6>(2, 0);
  dp[1] = 1;
  while (t--) {
    auto n = 0;
    cin >> n;
    if (n < dp.size()) {
      cout << dp[n] << endl;
      continue;
    }
    for (auto i = last_calculated + 1; i <= n; ++i) {
      if (i == next_row_from) {
        prev_row_end = i - 1;
        prev_row_start = 1 + prev_row_end - row_width;
        ++row_width;
        next_row_from = i + row_width;
      }
      auto up_left_index = i - row_width;
      auto up_right_index = 1 + up_left_index;
      auto prior_pyramid_index = 2 + i - (2 * row_width);
      dp.push_back(i * i);
      auto left_valid = prev_row_start <= up_left_index && up_left_index <= prev_row_end;
      auto right_valid = prev_row_start <= up_right_index && up_right_index <= prev_row_end;
      if (left_valid && right_valid) { dp[i] -= dp[prior_pyramid_index]; }
      if (left_valid) { dp[i] += dp[up_left_index]; }
      if (right_valid) { dp[i] += dp[up_right_index]; }
      last_calculated = n;
    }
    cout << dp[n] << endl;
  }

  return 0;
}
