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
        auto n=0;
        cin>>n;
        auto digit=0;
        cin>>digit;
        auto number = string{""};
        cin>>number;

        auto i=0;
        for(; i<n; ++i) {
            if(digit > number[i] - '0') {
                break;
            }
        }
        auto start = i > 0 ? number.substr(0, i) : "";
        auto end = i < n ? number.substr(i, n - i) : "";
        cout << start << digit << end << endl;
    }

    return 0;
}
