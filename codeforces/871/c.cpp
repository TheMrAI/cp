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
        auto n = 0;
        cin>>n;
        auto min_first = numeric_limits<int>::max();
        auto min_second = numeric_limits<int>::max();
        auto min_cost = numeric_limits<int>::max();
        while(n--) {
            auto cost = 0;
            cin>>cost;
            auto skill = string{};
            cin>>skill;
            if(skill == "11") {
                min_cost = min(min_cost, cost);
            } else {
                if(skill[0] == '1') {
                    min_first = min(min_first, cost);
                }
                if(skill[1] == '1') {
                    min_second = min(min_second, cost);
                }
            }
        }

        auto duo_cost = numeric_limits<int>::max();
        if(min_first != numeric_limits<int>::max() && min_second != numeric_limits<int>::max()) {
            duo_cost = min_first + min_second;
        }
        if(duo_cost != numeric_limits<int>::max() && duo_cost <= min_cost ) {
            cout << duo_cost << endl;
            continue;
        }
        if(min_cost == numeric_limits<int>::max()) {
            cout << -1 << endl;
            continue;
        }
        cout << min_cost << endl;
    }

    return 0;
}
