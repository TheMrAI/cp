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

  auto expected_pi = string{ "314159265358979323846264338327" };

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto number = string{};
    cin >> number;

    auto correct = 0;
    for (auto i = size_t{ 0 }; i < number.size() && i < expected_pi.size(); ++i) {
      if (number[i] != expected_pi[i]) { break; }
      ++correct;
    }
    cout << correct << endl;
  }
  return 0;
}
