#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto a = int64_t{ 0 };
    auto b = int64_t{ 0 };
    auto c = int64_t{ 0 };

    cin >> a >> b >> c;
    if (a + b == c) {
      cout << "+" << endl;
    } else {
      cout << "-" << endl;
    }
  }

  return 0;
}
