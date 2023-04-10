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
        int n = 0;
        cin>>n;
        auto numbers = vector<int>(n, 0);
        for(int i = 0; i < n; ++i) {
            cin>>numbers[i];
            if (numbers[i] == 1) {
                ++numbers[i];
            }
        }

        for(int i = 1; i < n; ++i) {
            if(numbers[i] % numbers[i-1] == 0) {
                ++numbers[i];
            }
        }
        for(auto num : numbers) {
            cout << num << " ";
        }
        cout << endl;
    }

    return 0;
}
