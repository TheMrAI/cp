#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto n = int64_t{ 0 };
    cin >> n;
    auto m = int64_t{ 0 };
    cin >> m;

    auto number_columns = vector<vector<int64_t>>(m);
    for (auto i = 0; i < n; ++i) {
      for (auto j = 0; j < m; ++j) {
        auto value = int64_t{ 0 };
        cin >> value;
        number_columns[j].push_back(value);
      }
    }

    for (auto i = 0; i < m; ++i) { sort(number_columns[i].begin(), number_columns[i].end()); }

    // for(auto column : number_columns) {
    //     for (auto val: column) {
    //         cout << val << " ";
    //     }
    //     cout <<endl;
    // }

    auto chips = int64_t{ 0 };
    for (auto k = 0; k < m; ++k) {
      for (auto i = 0; i < n; ++i) { chips += number_columns[k][i] * (1 + 2 * i - n); }
    }
    cout << chips << endl;
  }

  return 0;
}
