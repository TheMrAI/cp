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
    constexpr auto big_prime = static_cast<size_t>(1e9+7);
    while(t--) {
        auto n = 0;
        cin>>n;
        auto k = 0;
        cin>>k;
        auto bitmask_counts = array<size_t, 64>{};
        while(n--) {
            auto val = 0;
            cin>>val;
            auto tmp = bitmask_counts;
            for(auto i = size_t{0}; i < size_t{64}; ++i) {
                auto target = i & val;
                bitmask_counts[target] = (bitmask_counts[target] + tmp[i]) % big_prime;
            }
            bitmask_counts[val] = (bitmask_counts[val] + 1) % big_prime;
        }
        auto total = 0;
        for(auto i = size_t{0}; i < size_t{64}; ++i) {
            if(popcount(i) == k) {
                total = (total + bitmask_counts[i]) % big_prime; 
            }
        }
        cout << total << endl;
    }

    return 0;
}
