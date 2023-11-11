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
    auto lhs = ""s;
    auto rhs = ""s;
    cin >> lhs >> rhs;
    auto block_time = 0;
    cin >> block_time;
    auto q = 0;
    cin >> q;

    auto bad_count = 0;
    auto bad_chars = vector<bool>(lhs.size(), false);
    auto block_expiry = queue<pair<int, int>>{};
    for (auto i = 0; i < lhs.size(); ++i) {
      if (lhs[i] != rhs[i]) {
        bad_chars[i] = true;
        ++bad_count;
      }
    }

    auto time_point = 0;
    while (q--) {
      while (!block_expiry.empty() && block_expiry.front().first == time_point) {
        auto [_, index] = block_expiry.front();
        block_expiry.pop();
        if (bad_chars[index] == true) { ++bad_count; }
      }
      auto command_type = 0;
      cin >> command_type;
      switch (command_type) {
      case 1: {
        auto blocking_index = 0;
        cin >> blocking_index;
        if (bad_chars[--blocking_index] == true) { --bad_count; }
        block_expiry.emplace(time_point + block_time, blocking_index);
        break;
      }
      case 2: {
        auto* lhs_ptr = &lhs;
        auto* rhs_ptr = &rhs;

        auto first = 1;
        cin >> first;
        if (first == 2) { lhs_ptr = &rhs; }
        auto swap_from = 0;
        cin >> swap_from;
        --swap_from;

        auto second = 2;
        cin >> second;
        if (second == 1) { rhs_ptr = &lhs; }
        auto swap_to = 0;
        cin >> swap_to;
        --swap_to;

        auto swap_from_bad = bad_chars[swap_from];
        auto swap_to_bad = bad_chars[swap_to];

        swap((*lhs_ptr)[swap_from], (*rhs_ptr)[swap_to]);
        if (lhs[swap_from] != rhs[swap_from] && !swap_from_bad) {
          bad_chars[swap_from] = true;
          ++bad_count;
        }
        if (lhs[swap_from] == rhs[swap_from] && swap_from_bad) {
          bad_chars[swap_from] = false;
          --bad_count;
        }
        if (lhs[swap_to] != rhs[swap_to] && !swap_to_bad) {
          bad_chars[swap_to] = true;
          ++bad_count;
        }
        if (lhs[swap_to] == rhs[swap_to] && swap_to_bad) {
          bad_chars[swap_to] = false;
          --bad_count;
        }
        break;
      }
      case 3: {
        if (bad_count == 0) {
          cout << "YES" << endl;
        } else {
          cout << "NO" << endl;
        }
        break;
      }
      }
      ++time_point;
    }
  }

  return 0;
}
