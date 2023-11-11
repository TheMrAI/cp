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
  auto codeforces = string{ "codeforces" };
  while (t--) {
    auto s = string{};
    cin >> s;
    auto mismatch = 0;
    for (auto i = 0; i < s.size(); ++i) {
      if (codeforces[i] != s[i]) { ++mismatch; }
    }
    cout << mismatch << endl;
  }

  return 0;
}
