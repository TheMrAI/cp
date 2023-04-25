#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n = i6{0};
        cin>>n;
        auto max_beauty = numeric_limits<i6>::min();
        auto min_v = i6{0};
        auto max_v = i6{0};
        for(auto i = 0; i < n; ++i) {
            auto val = i6{0};
            cin>>val;
            if(i == 0) {
                min_v = val;
                max_v = val;
                continue;
            }
            max_beauty = max({max_beauty, min_v * val, max_v * val});
            min_v = min(min_v, val);
            max_v = max(max_v, val);
        }
        cout << max_beauty << endl;
    }

    return 0;
}
