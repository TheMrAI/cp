#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = 0;
    cin >> n;
    auto start = 1;
    auto delta = 1;
    if (n % 2 == 0) {
      start = 2;
      delta = 2;
    }
    for (auto i = 0; i < n; ++i) {
      cout << start << " ";
      start += delta;
    }
    cout << endl;
  }

  return 0;
}
