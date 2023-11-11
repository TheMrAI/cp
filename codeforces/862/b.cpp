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

    auto data = string{ "" };
    cin >> data;
    auto min_index = 0;
    for (auto i = 0; i < data.size(); ++i) {
      if (data[i] <= data[min_index]) { min_index = i; }
    }
    cout << string{ data[min_index] } + data.substr(0, min_index) + data.substr(min_index + 1) << endl;
  }

  return 0;
}
