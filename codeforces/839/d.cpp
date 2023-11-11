#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
requires std::derived_from<Output, std::ostream> && std::weakly_incrementable<Iter> && std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end)
{
  while (begin != end) {
    output << *begin;
    if (std::next(begin) != end) { output << ' '; }
    ++begin;
  }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = 0;
    cin >> n;
    auto values = vector<i6>(n, 0);
    for (auto i = 0; i < n; ++i) { cin >> values[i]; }
    auto top_x = numeric_limits<i6>::max();
    auto bottom_x = i6{ 0 };// as per constraints
    for (auto i = 0; i + 1 < n; ++i) {
      auto delta = values[i + 1] - values[i];
      auto sum = values[i + 1] + values[i];
      if (delta == 0) {
        continue;// no changes in constraints
      }
      if (delta < 0) {// order has to change
        bottom_x = max(bottom_x, ceil2(sum, i6{ 2 }));
      } else {
        top_x = min(top_x, sum / i6{ 2 });
      }
    }
    // cout << top_x << " " << bottom_x << endl;
    if (top_x < bottom_x) {
      cout << -1 << endl;
      continue;
    }
    cout << bottom_x << endl;
  }

  return 0;
}
