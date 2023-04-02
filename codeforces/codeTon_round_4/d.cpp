#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto min_climb(i6 climb, i6 slide, i6 n) -> i6 {
    return (n-2)*(climb-slide)+climb+1;
}

auto max_climb(i6 climb, i6 slide, i6 n) -> i6 {
    return (n-1)*(climb-slide)+climb;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto events=i6{0};
        cin>>events;
        auto h_min=i6{1};
        auto h_max=static_cast<i6>(1e18);

        while(events--) {
            auto type = i6{0};
            cin>>type;
            if(type == 1) {
                auto climb=i6{0};
                cin>>climb;
                auto slide=i6{0};
                cin>>slide;
                auto n=i6{0};
                cin>>n;

                auto h_min_prop = min_climb(climb, slide, n);
                auto h_max_prop = max_climb(climb, slide, n);
                if (n == 1) {
                    h_min_prop = 1;
                    h_max_prop = climb;
                }

                if(h_min_prop > h_max || h_max_prop < h_min) {
                    cout<<0<<" ";
                } else {
                    h_min = max(h_min, h_min_prop);
                    h_max = min(h_max, h_max_prop);
                    cout<<1<<" ";
                }
            } else {
                auto climb=i6{0};
                cin>>climb;
                auto slide=i6{0};
                cin>>slide;
                auto valid_n_1 = max(i6{1}, (h_min - slide - 1) / (climb - slide) + 1);
                auto valid_n_2 = max(i6{1}, (h_max - slide - 1) / (climb - slide) + 1);

                if (valid_n_1 == valid_n_2) {
                    cout<<valid_n_1<<" ";
                } else {
                    cout<<-1<<" ";
                }
            }
        }
        cout<<endl;
    }

    return 0;
}