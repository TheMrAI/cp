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

    auto data = string{};
    cin >> data;
    for (int i = 0; i < data.size(); ++i) { data[i] = tolower(data[i]); }
    auto last = unique(data.begin(), data.end());
    data.erase(last, data.end());

    if (data == "meow") {
      cout << "YES" << endl;
    } else {
      cout << "NO" << endl;
    }
  }

  return 0;
}
