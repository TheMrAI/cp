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
    auto s = ""s;
    cin >> s;
    auto longest_sequence = 1;
    auto current_sequence = 1;
    auto sequence_char = s[0];
    for (auto i = 1; i < n; ++i) {
      if (s[i] == sequence_char) {
        ++current_sequence;
        longest_sequence = max(longest_sequence, current_sequence);
        continue;
      }
      current_sequence = 1;
      sequence_char = s[i];
    }
    cout << longest_sequence + 1 << endl;
  }

  return 0;
}
