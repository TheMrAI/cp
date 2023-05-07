#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://cp-algorithms.com/algebra/euclid-algorithm.html
auto e_gcd(int a, int b) -> int {
    if (b > a) {
        swap(a, b);
    }
    if (b == 0) {
        return a;
    }
    while(b != 0) {
        a = a%b;
        swap(a, b);
    }
    return a;
}

auto e_lcm(int a, int b) -> int {
    return a / e_gcd(a, b) * b;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    int a, b;
    cin>>a>>b;
    cout << e_gcd(a,b) << endl;

    return 0;
}
