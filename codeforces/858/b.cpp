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
    int n = 0;
    cin >> n;

    auto numbers = map<int, int>{};
    for (int i = 0; i < n; ++i) {
      int a = 0;
      cin >> a;
      ++numbers[a];
    }

    auto zeros = numbers.find(0);
    auto ones = numbers.find(1);
    // cout << "Half ceil " << ceil2(n, 2) << endl;
    if (zeros == numbers.end() || zeros->second <= ceil2(n, 2)) {
      cout << 0 << endl;
    } else if (ones == numbers.end() || n - zeros->second - ones->second >= 1) {
      cout << 1 << endl;
    } else {
      cout << 2 << endl;
    }
  }

  return 0;
}
