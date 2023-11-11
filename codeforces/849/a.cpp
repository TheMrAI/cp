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
    auto character = ' ';
    cin >> character;
    auto found = false;
    auto word = "codeforces"s;
    for (auto c : word) {
      if (c == character) {
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
