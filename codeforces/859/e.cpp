#include <bits/stdc++.h>

using namespace std;

auto ask(vector<int64_t> const& stone_sums) -> int64_t
{
  auto l = int64_t{ 1 };
  auto h = static_cast<int64_t>(stone_sums.size() - 1);

  while (l < h) {
    auto middle = l + floorl((h - l) / 2);
    cout << "? " << 1 + middle - l << " ";
    for (auto i = l; i <= middle; ++i) { cout << i << " "; }
    cout << endl;
    auto sum = int64_t{ 0 };
    cin >> sum;
    if (sum > (stone_sums[middle] - stone_sums[l - 1])) {
      h = middle;
    } else {
      l = middle + 1;
    }
  }
  return l;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto n = int64_t{ 0 };
    cin >> n;
    auto stone_sums = vector<int64_t>(n + 1, 0);
    for (auto i = int64_t{ 1 }; i <= n; ++i) {
      auto count = int64_t{};
      cin >> count;
      stone_sums[i] = stone_sums[i - 1] + count;
    }

    auto result = ask(stone_sums);
    cout << "! " << result << endl;
  }

  return 0;
}
