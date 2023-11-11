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
    auto data = string{};
    cin >> data;
    auto l = 0;
    auto h = data.size() - 1;
    while (l < h) {
      if (data[l] == data[h]) { break; }
      ++l;
      --h;
    }
    cout << data.size() - (2 * l) << endl;
  }

  return 0;
}
