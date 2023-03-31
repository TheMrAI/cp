#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto main() -> int {
    // ios::sync_with_stdio(false);
    // cin.tie(nullptr);

    auto l=i6{1};
    auto h=i6{1'000'000};
    while(l<h) {
        auto mid = (l + h + 1) / 2;
        cout<<mid<<endl;
        auto verdict = string{};
        cin>>verdict;
        if(verdict == string{"<"}) {
            h = mid - 1;
        } else {
            l = mid;
        }
    }
    cout << "! " << l << endl; 

    return 0;
}
