#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto value = int64_t{ 0 };
    cin >> value;

    auto low = int64_t{ 0 };
    auto high = value;
    while (low < high) {
      auto distance = low + floorl((high - low) / 2);

      if (distance * distance >= value) {
        high = distance;
      } else {
        low = distance + 1;
      }
    }
    cout << high - 1 << endl;
  }

  return 0;
}
