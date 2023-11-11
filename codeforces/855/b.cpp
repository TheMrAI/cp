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
    int n = 0;
    cin >> n;
    int k = 0;
    cin >> k;

    auto lower_counts = vector<int>(26, 0);
    auto upper_counts = vector<int>(26, 0);
    auto word = string{};
    cin >> word;
    for (auto c : word) {
      if (isupper(c)) {
        ++upper_counts[c - 'A'];
      } else {
        ++lower_counts[c - 'a'];
      }
    }
    auto burls = 0;
    for (auto i = 0; i < 26; ++i) {
      auto match = min(lower_counts[i], upper_counts[i]);
      burls += match;
      lower_counts[i] -= match;
      upper_counts[i] -= match;
    }
    for (auto i = 0; i < 26; ++i) {
      if (k > 0) {
        auto unpaired = max(lower_counts[i], upper_counts[i]);
        auto to_pair = min(unpaired / 2, k);
        burls += to_pair;
        k -= to_pair;
      }
    }
    cout << burls << endl;
  }

  return 0;
}
