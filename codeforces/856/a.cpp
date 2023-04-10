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
        int n=0;
        cin>>n;
        auto longest = vector<string>{};
        for(int i = 0; i < (2*n) - 2; ++i) {
            auto sequence = string{""};
            cin>>sequence;
            if(sequence.size() == n-1) {
                longest.push_back(sequence);
            }
        }

        auto section_1 = longest[0];
        auto section_2 = longest[1];
        reverse(section_2.begin(), section_2.end());

        if (section_1 == section_2) {
            cout << "YES" << endl;
        } else {
            cout << "NO" << endl;
        }
    }

    return 0;
}
