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
        auto n = i6{0};
        cin>>n;

        if (n % 2 != 0 && n > 1) {
            cout << -1 << endl;
            continue;
        }

        auto bottom = 1;
        auto top = n - 1;
        auto use_bottom = false;
        auto last_remainder = 0;
        auto result = vector<int>(1, n);
        while(bottom <= top) {
            auto number = last_remainder;
            if(use_bottom) {
                if(number < bottom) {
                    number = bottom - number;
                } else if (number > bottom) {
                    number = (n - number) + bottom;
                }
                last_remainder = bottom;
                bottom += 1;
            } else {
                if(number < top) {
                    number = top - number;
                } else if (number > top) {
                    number = (n - number) + top;
                }
                last_remainder = top;
                top -= 1;
            }
            result.push_back(number);
            use_bottom = not use_bottom;
        }
        for(auto i = 0; i < result.size(); ++i) {
            cout << result[i] << " ";
        }
        cout << endl;
    }

    return 0;
}
