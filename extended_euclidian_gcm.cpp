#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://cp-algorithms.com/algebra/extended-euclid-algorithm.html
auto extended_e_gcm_rec(int a, int b, int& x, int& y) -> int
{
  if (b > a) { swap(b, a); }
  if (b == 0) {
    x = 1;
    y = 0;
    return a;
  }
  int x1;
  int y1;
  int greatest_divisor = extended_e_gcm_rec(b, a % b, x1, y1);
  x = y1;
  y = x1 - y1 * (a / b);
  return greatest_divisor;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  int a, b;
  cin >> a >> b;
  int x, y;
  cout << "gcm(" << a << ", " << b << ")= " << extended_e_gcm_rec(a, b, x, y) << ", x = " << x << ", y = " << y << endl;

  return 0;
}
