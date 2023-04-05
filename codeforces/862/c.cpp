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
        auto line_count = 0;
        cin>>line_count;
        auto parabola_count = 0;
        cin>>parabola_count;

        auto kay_v = vector<i6>{};
        for (auto i=0; i < line_count; ++i) {
            auto val= 0;
            cin>>val;
            kay_v.push_back(val);
        }
        sort(kay_v.begin(), kay_v.end());
        
        for (auto i=0; i < parabola_count; ++i) {
            auto a=i6{};
            auto b=i6{};
            auto c=i6{};
            cin>>a>>b>>c;
            
            auto index = distance(kay_v.begin(), lower_bound(kay_v.begin(), kay_v.end(), b));
            if(index < line_count && (b - kay_v[index]) * (b - kay_v[index]) < 4 * a * c) {
                cout << "YES" << endl;
                cout << kay_v[index] << endl;
                continue;
            }
            if(index > 0 && (b - kay_v[index-1]) * (b - kay_v[index-1]) < 4 * a * c) {
                cout << "YES" << endl;
                cout << kay_v[index-1] << endl;
                continue;
            }
            cout << "NO" << endl;
        }
        cout << endl;
    }

    return 0;
}
