#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = 0;
  cin >> t;
  while (t--) {
    auto a = 0;
    auto b = 0;
    auto c = 0;
    auto d = 0;
    cin >> a;
    cin >> b;
    cin >> c;
    cin >> d;

    auto x = 1;
    while ((a + x) % b != (c + x) % d) { ++x; }
    cout << x << endl;
  }
  return 0;
}