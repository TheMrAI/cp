#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto desired_candy=i6{0};
        cin>>desired_candy;
        if(desired_candy %2==0) {
            cout<<-1<<endl;
        } else {
            auto max_step_count = i6{40};
            auto spells = vector<int>{};
            while(desired_candy != 1 && max_step_count--) {
                auto val=(desired_candy-1)/2;
                if(val % 2 != 0) {
                    spells.push_back(2);
                    desired_candy = val;
                    continue;
                }
                val = (desired_candy+1)/2;
                if(val % 2 != 0) {
                    spells.push_back(1);
                    desired_candy = val;
                    continue;
                }
            }
            if(desired_candy == 1) {
                cout << spells.size() << endl;
                for(auto v = spells.rbegin(); v!=spells.rend(); ++v ) {
                    cout << *v << " ";
                }
                cout << endl;
            } else {
                cout<<-1<<endl;
            }
        }
    }

    return 0;
}
