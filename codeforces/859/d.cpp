#include <bits/stdc++.h>

using namespace std;

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t=int64_t{0};
    cin>>t;
    while(t--) {
        auto n=int64_t{0};
        auto q=int64_t{0};
        cin>>n>>q;
        auto rolling_sum = vector<int64_t>(n+1);
        rolling_sum[0] = 0;
        for(auto i=int64_t{1}; i<n+1; ++i) {
            auto val = int64_t{0};
            cin>>val;
            rolling_sum[i] = rolling_sum[i-1] + val;
        }

        for(auto i=int64_t{0}; i<q; ++i) {
            auto l=int64_t{0};
            auto r=int64_t{0};
            auto k=int64_t{0};
            cin>>l>>r>>k;

            auto sum = rolling_sum[rolling_sum.size()-1];
            sum -= rolling_sum[r];
            sum += rolling_sum[l-1];
            sum += (1 + r-l)*k;
            if(sum % 2!=0) {
                cout << "YES" << endl;
            } else {
                cout << "NO" << endl;
            }
        }
    }

    return 0;
}
