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

auto count_options(i6 counter, i6 k)
{
  auto limit = counter - (k - 1);
  return (limit * (limit + 1)) / 2;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = 0;
    auto k = i6{ 0 };
    auto q = i6{ 0 };
    cin >> n >> k >> q;
    auto counter = i6{ 0 };
    auto options = i6{ 0 };
    auto val = i6{ 0 };
    while (n--) {
      cin >> val;
      if (val > q) {
        if (counter >= k) { options += count_options(counter, k); }
        counter = 0;
      } else {
        ++counter;
      }
    }
    if (counter >= k) { options += count_options(counter, k); }
    cout << options << endl;
  }

  return 0;
}
