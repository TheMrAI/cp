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
        cin >> n;
        auto data = string{};
        cin >> data;
        auto right_chars = array<int, 26>{};
        auto left_chars = array<int, 26>{};
        for(auto c : data) {
            right_chars[c - 'a'] += 1;
        }
        auto max_distinct_split = count_if(right_chars.begin(), right_chars.end(), [](auto const& val) {
            return val != 0;
        });
        for(auto i = 0; i < data.size(); ++i) {
            left_chars[data[i]-'a'] += 1;
            right_chars[data[i]-'a'] -= 1;
            max_distinct_split = max(max_distinct_split, count_if(left_chars.begin(), left_chars.end(), [](auto const& val) {
            return val != 0;
        }) + count_if(right_chars.begin(), right_chars.end(), [](auto const& val) {
            return val != 0;
        }));
        }
        cout << max_distinct_split << endl;
    }

    return 0;
}
