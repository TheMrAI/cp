#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  auto count = 0;
  while (t--) {
    auto p = 0;
    auto v = 0;
    auto t = 0;
    cin >> p >> v >> t;


    if (p + v + t >= 2) { ++count; }
  }
  cout << count << endl;

  return 0;
}
