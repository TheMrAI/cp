#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto n = i6{ 0 };
    cin >> n;
    auto sequence = vector<i6>{};
    for (auto i = 0; i < n; ++i) {
      auto val = i6{ 0 };
      cin >> val;
      sequence.push_back(val);
    }

    auto beautiful = false;
    for (auto i = 0; i < n; ++i) {
      if (sequence[i] <= i + 1) {
        beautiful = true;
        break;
      }
    }
    if (beautiful) {
      cout << "YES" << endl;
    } else {
      cout << "NO" << endl;
    }
  }

  return 0;
}
