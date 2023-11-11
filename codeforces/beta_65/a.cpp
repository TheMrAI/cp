#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = int64_t{ 0 };
  cin >> t;
  while (t--) {
    auto word = string{};
    cin >> word;
    if (word.size() > 10) {
      cout << word[0] << word.size() - 2 << word[word.size() - 1] << endl;
    } else {
      cout << word << endl;
    }
  }

  return 0;
}
