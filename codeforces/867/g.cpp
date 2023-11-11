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

  constexpr auto max_value = static_cast<i6>(1e6);
  auto counter = array<i6, static_cast<i6>(1e6) + 1>{};

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = 0;
    cin >> n;
    auto numbers = vector<i6>{};
    while (n--) {
      auto val = i6{ 0 };
      cin >> val;
      numbers.push_back(val);
      counter[val] += 1;
    }
    auto magic_triplets = i6{ 0 };
    for (auto i = 0; i < numbers.size(); ++i) {
      magic_triplets += (counter[numbers[i]] - 1) * (counter[numbers[i]] - 2);
      for (auto b = i6{ 2 }; numbers[i] * b * b <= max_value; ++b) {
        magic_triplets += counter[numbers[i] * b] * counter[numbers[i] * b * b];
      }
    }
    cout << magic_triplets << endl;
    for (int i = 0; i < numbers.size(); ++i) { --counter[numbers[i]]; }
  }

  return 0;
}