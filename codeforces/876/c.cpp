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
    auto a = vector<int>(n, 0);
    for (auto i = 0; i < n; ++i) { cin >> a[i]; }
    if (a.back() == 1) {
      cout << "NO" << endl;
      continue;
    }
    auto last = static_cast<int>(a.size()) - 1;
    auto indexes = vector<int>{};
    while (last >= 0) {
      auto prior = last - 1;
      if (prior < 0) {
        indexes.push_back(0);
      } else {
        while (prior >= 0 && a[prior] == 1) {
          indexes.push_back(0);
          --prior;
        }
        indexes.push_back(last - prior - 1);
      }
      last = prior;
    }
    cout << "YES" << endl;
    dump_to(cout, indexes.cbegin(), indexes.cend());
    cout << endl;
  }

  return 0;
}
