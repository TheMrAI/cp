#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto l = int64_t{ 0 };
  auto r = int64_t{ 0 };
  cin >> l;
  cin >> r;

  auto numbers = map<int, vector<int64_t>>{};
  for (; l <= r; ++l) {
    auto d = to_string(l);
    sort(d.begin(), d.end());
    auto luckiness = d[d.size() - 1] - d[0];
    numbers[luckiness].push_back(l);
  }

  auto max = 100;
  for (auto l : numbers) {
    cout << "Unluckiness: " << l.first << endl;
    for (auto v : l.second) {
      cout << v << endl;
      --max;
      if (max == 0) {
        cout << "..." << endl;
        break;
      }
    }
    break;
  }

  return 0;
}
