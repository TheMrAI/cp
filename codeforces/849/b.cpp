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
    auto length = 0;
    cin >> length;
    auto directions = string{};
    cin >> directions;
    auto location = make_pair(0, 0);
    auto found = false;
    for (auto dir : directions) {
      switch (dir) {
      case 'U':
        location.second += 1;
        break;
      case 'D':
        location.second -= 1;
        break;
      case 'L':
        location.first -= 1;
        break;
      case 'R':
        location.first += 1;
        break;
      }
      if (location.first == 1 && location.second == 1) {
        found = true;
        break;
      }
    }
    if (found) {
      cout << "YES" << endl;
      continue;
    }
    cout << "NO" << endl;
  }

  return 0;
}
