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
    auto x = 0;
    cin >> x;
    auto k = 0;
    cin >> k;
    if (x % k != 0) {
      cout << 1 << endl;
      cout << x << endl;
      continue;
    }
    auto low = 1;
    auto high = x - 1;
    while (low % k == 0) { ++low; }
    while (high % k == 0) { --high; }
    while (low + high != x && low <= high) {
      auto sum = low + high;
      if (sum < x) {
        while (low % k == 0) { ++low; }
      } else if (sum > x) {
        while (high % k == 0) { --high; }
      }
    }
    cout << 2 << endl;
    cout << low << " " << high << endl;
  }

  return 0;
}
