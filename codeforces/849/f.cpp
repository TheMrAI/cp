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
        auto n = 0;
        auto q = 0;
        cin>>n>>q;
        auto data = vector<int>(n, 0);
        auto can_be_updated = set<int>{};
        for(auto i = 0; i < n; ++i) {
            cin>>data[i];
            if(data[i] >= 10) {
                can_be_updated.insert(i);
            }
        }
        while(q--) {
            auto type = 0;
            cin>>type;
            if(type == 2) {
                auto index = 0;
                cin>>index;
                cout<<data[index-1]<<endl;
            } else if(type == 1) {
                auto l = 0;
                auto h = 0;
                cin>>l>>h;
                auto to_update = can_be_updated.lower_bound(l-1);
                while(to_update != can_be_updated.end() && *to_update <= h-1) {
                    auto updated_val = 0;
                    auto index = *to_update;
                    while(data[index] != 0) {
                        updated_val += data[index] % 10;
                        data[index] = data[index]/10;
                    }
                    data[index] = updated_val;
                    if(data[index] < 10) {
                        to_update = can_be_updated.erase(to_update);
                        continue;
                    }
                    ++to_update;
                }
            }
        }
    }

    return 0;
}
