#include <bits/stdc++.h>

using namespace std;

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t=int64_t{0};
    cin>>t;
    auto groups=1;
    auto prev_magnet = string{};
    cin>>prev_magnet;
    --t;
    while(t--) {
        auto magnet = string{};
        cin>>magnet;
        if(prev_magnet[1] == magnet[0]) {
            ++groups;
        }
        swap(prev_magnet, magnet);
    }
    cout<< groups<<endl;

    return 0;
}
