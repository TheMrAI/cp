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
        auto c = 0;
        cin>>n>>c;
        auto candidates = vector<pair<int, int>>{};
        for(auto i = 0; i < n; ++i) {
            auto value = 0;
            cin>>value;
            candidates.emplace_back(value + min(i + 1, (n - i)), value + i + 1);
        }
        ranges::sort(candidates);

        auto teleport_cost_pref_sums = vector<i6>(n, 0);
        teleport_cost_pref_sums[0] = candidates[0].first;
        for(auto i = 1; i < n; ++i) {
            teleport_cost_pref_sums[i] += teleport_cost_pref_sums[i-1] + candidates[i].first;
        }
        // for_each(candidates.begin(), candidates.end(), [](auto const& l) { cout << l.first << " - " << l.second << ", ";});
        // cout << endl;
        auto teleports_used = 0;
        for(auto i = 0; i < n; ++i) {
            auto remaining_c = c - candidates[i].second;
            if(remaining_c >= 0) {
                teleports_used = max(teleports_used, 1);
            }
            // cout << "Remaining " << remaining_c << endl;
            auto l = 0;
            auto h = n-1;
            while (l <= h) {
                // cout << "L " << l << " H " << h << endl;
                auto m = l + ((h-l)/2);
                // cout << "M " << m << endl;
                auto cost = teleport_cost_pref_sums[m];
                // cout << "Cost " << cost << endl;
                auto count = m+2;
                if(m >= i) {
                    cost -= candidates[i].first;
                    // cout << "Mod Cost " << cost << endl;
                    --count;
                }
                if(cost <= remaining_c) {
                    // cout << "C " << count << endl;
                    teleports_used = max(teleports_used, count);
                    l = m + 1;
                } else {
                    h = m - 1;
                }
            }
        }
        cout << teleports_used << endl;
    }

    return 0;
}
