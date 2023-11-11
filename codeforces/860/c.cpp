#include <bits/stdc++.h>

using namespace std;

auto main() -> int
{
  ios::sync_with_stdio(0);
  cin.tie(0);

  auto t = 0;
  cin >> t;
  while (t--) {
    int64_t n = 0;
    cin >> n;
    vector<pair<int64_t, int64_t>> candies;
    for (auto i = 0; i < n; ++i) {
      int64_t a = 0;
      int64_t b = 0;
      cin >> a;
      cin >> b;
      candies.emplace_back(make_pair(a, b));
    }

    auto prev_gcd = candies[0].first * candies[0].second;
    auto prev_lcm = candies[0].second;
    auto individual_price_tags = 1;
    for (auto i = int64_t{ 1 }; i < n; ++i) {
      auto multiple = candies[i].first * candies[i].second;
      auto new_gcd = gcd(prev_gcd, multiple);
      auto new_lcm = lcm(prev_lcm, candies[i].second);

      if (new_gcd % new_lcm == 0) {
        prev_gcd = new_gcd;
        prev_lcm = new_lcm;
      } else {
        prev_gcd = multiple;
        prev_lcm = candies[i].second;
        ++individual_price_tags;
      }
    }
    cout << individual_price_tags << endl;
  }


  return 0;
}