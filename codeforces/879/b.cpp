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
    auto lower_bound = ""s;
    auto upper_bound = ""s;
    cin >> lower_bound >> upper_bound;
    if (upper_bound.size() > lower_bound.size()) {
      auto starter_diff = upper_bound[0] - '0';
      cout << static_cast<uint64_t>(starter_diff) + uint64_t{ 9 } * (upper_bound.size() - 1) << endl;
    } else {// can only be equal length
      auto index = 0;
      while (index < upper_bound.size() && lower_bound[index] == upper_bound[index]) { ++index; }
      auto starter_diff = 0;
      if (index < upper_bound.size()) {
        starter_diff = std::abs(upper_bound[index] - lower_bound[index]);
        ++index;
      }
      cout << static_cast<uint64_t>(starter_diff) + uint64_t{ 9 } * (upper_bound.size() - index) << endl;
    }
  }

  return 0;
}
