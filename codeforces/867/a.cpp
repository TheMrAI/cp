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
        auto lunch_time = i6{0};
        cin>>lunch_time;
        auto durations = vector<i6>(n, 0);
        for(auto i = 0; i < n; ++i) {
            cin>>durations[i];
        }
        auto value = vector<i6>(n, 0);
        for(auto i = 0; i < n; ++i) {
            cin>>value[i];
        }

        auto max_index = -1;
        auto cost = 0;
        for(auto i=0; i < n; ++i) {
            if(cost + durations[i] <= lunch_time) {
                if(max_index == -1 || value[max_index-1] < value[i]) {
                    max_index = i+1;
                }
            }
            ++cost;
        }
        cout << max_index << endl;
    }

    return 0;
}
