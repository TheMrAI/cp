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
    int x1, y1, x2, y2;
    cin >> x1 >> y1 >> x2 >> y2;
    int a = y2 - y1;
    if (a < 0) {
      cout << -1 << endl;
      continue;
    }
    int b = x1 - x2 + a;
    if (b < 0) {
      cout << -1 << endl;
      continue;
    }
    cout << a + b << endl;
  }

  return 0;
}
