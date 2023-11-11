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
    auto powerups = priority_queue<int>{};
    auto max_power = i6{ 0 };
    for (auto i = 0; i < n; ++i) {
      auto val = 0;
      cin >> val;
      if (val == 0 && !powerups.empty()) {
        max_power += powerups.top();
        powerups.pop();
      } else {
        powerups.push(val);
      }
    }
    cout << max_power << endl;
  }

  return 0;
}
