#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  string data;
  getline(cin, data);

  while (t--) {
    auto colors = unordered_map<int, int>{};
    getline(cin, data);

    for (auto i = 0; i < 4; ++i) {
      auto color = data[i] - '0';
      ++colors[color];
    }

    auto types = colors.size();
    if (types == 4) {
      cout << 4 << endl;
    } else if (types == 3) {
      cout << 4 << endl;
    } else if (types == 1) {
      cout << -1 << endl;
    } else {
      auto t1_count = colors.begin()->second;
      auto t2_count = next(colors.begin(), 1)->second;
      if (t1_count % 2 != 0 && t2_count % 2 != 0) {
        cout << 6 << endl;
      } else {
        cout << 4 << endl;
      }
    }
  }

  return 0;
}
