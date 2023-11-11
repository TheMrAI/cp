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
    auto s = ""s;
    cin >> s;
    auto one_seen = false;
    for (auto i = 0; i < s.size(); ++i) {
      if (s[i] == '?') {
        if (one_seen) {
          s[i] = '1';
        } else {
          s[i] = '0';
        }
      } else if (s[i] == '1') {
        one_seen = true;
      } else {
        one_seen = false;
      }
    }
    cout << s << endl;
  }

  return 0;
}
