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
    auto a = vector<int>(n, 0);
    auto b = vector<int>(n - 1, 0);
    for (auto i = 0; i < n - 1; ++i) { cin >> b[i]; }

    a[0] = b[0];
    a[n - 1] = b[n - 2];
    for (auto i = 0; i < n - 2; ++i) { a[i + 1] = min(b[i], b[i + 1]); }

    for (auto val : a) { cout << val << " "; }
    cout << endl;
  }

  return 0;
}
