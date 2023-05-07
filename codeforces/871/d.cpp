#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto can_pile(int n, int m) -> bool {
    if(n == m) {
        return true;
    }
    if(n < m || n % 3 != 0) {
        return false;
    }
    auto piece = n / 3;
    return can_pile(piece, m) || can_pile(2*piece, m);
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n = 0;
        cin>>n;
        auto m = 0;
        cin>>m;
        if(can_pile(n, m)) {
            cout << "YES" << endl;
            continue;
        }
        cout << "NO" << endl;
    }

    return 0;
}
