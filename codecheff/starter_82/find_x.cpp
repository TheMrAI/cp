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

    auto x = 0;
    a = a % b;
    c = c % d;

    auto diff_ba = b - a;
    auto diff_dc = d - c;

    if (diff_ba > 1 && diff_dc > 1) {
      x = 1;
    } else if (diff_ba == diff_dc) {
      x = 1;
    } else {
      x = ((b * d) / gcd(b, d)) - (a % b);
    }
    cout << x << endl;
  }
  return 0;
}