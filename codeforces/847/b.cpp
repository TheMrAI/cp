#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
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
    auto s = 0;
    cin >> s;
    auto r = 0;
    cin >> r;

    auto sequence = vector<int>{};
    auto max_number = s - r;
    sequence.push_back(max_number);

    while (r != 0) {
      auto max_count = 0;
      if (max_number == 1) {
        max_count = n - static_cast<int>(sequence.size());
      } else {
        max_count = (r + static_cast<int>(sequence.size()) - n) / (max_number - 1);
      }

      for (auto i = 0; i < max_count; ++i) { sequence.push_back(max_number); }
      r -= max_count * max_number;
      --max_number;
    }

    for (auto number : sequence) { cout << number << " "; }
    cout << endl;
  }

  return 0;
}
