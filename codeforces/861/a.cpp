#include <bits/stdc++.h>

using namespace std;

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t=int64_t{0};
    cin>>t;
    while(t--) {
        int64_t l = 0;
        cin>>l;
        int64_t r = 0;
        cin>>r;

        int64_t luckiest = l;
        int64_t max_l = 0;
        for (auto i = l; i<=r && max_l<9; ++i) {
            auto values = to_string(i);
            sort(values.begin(), values.end());
            int64_t l = values[values.size()-1] - values[0];
            if (max_l < l) {
                luckiest = i;
                max_l = l;
            }
        }
        cout<<luckiest<<endl;
    }

    return 0;
}
