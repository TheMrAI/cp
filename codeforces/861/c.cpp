#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto l = int64_t{ 0 };
    auto r = int64_t{ 0 };
    cin >> l;
    cin >> r;

    auto unluckiest = l;
    auto luckiness = 9;
    for (; l <= r && luckiness != 0; ++l) {
      auto d = to_string(l);
      sort(d.begin(), d.end());
      auto luckiness_candidate = d[d.size() - 1] - d[0];
      if (luckiness_candidate < luckiness) {
        unluckiest = l;
        luckiness = luckiness_candidate;
      }
    }
    cout << unluckiest << endl;
  }

  return 0;
}
