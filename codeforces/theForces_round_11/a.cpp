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
        auto numbers = vector<i6>(n, 0);
        auto maxi = i6{0};
        for(auto i = 0; i < n; ++i) {
            cin>>numbers[i];
            if (numbers[i] >= numbers[maxi]){
                maxi = i;
            }
        }

        auto k_count = numbers.size() - maxi;
        for(auto i = 0; i < maxi; ++i) {
            if(numbers[i] == numbers[maxi]) {
                --k_count;
                break;
            }
        }
        cout << k_count << endl;
    }

    return 0;
}
