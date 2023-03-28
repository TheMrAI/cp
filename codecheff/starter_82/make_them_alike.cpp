#include <bits/stdc++.h>

using namespace std;

auto perm_transform(vector<int64_t> const & p, vector<pair<int64_t, int64_t>>& a) {
    for(auto i=int64_t{0}; i<a.size(); ++i) {
        a[i].first = a[p[i]].first;
        a[i].second = a[p[i]].second;
    }
}

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t=int64_t{0};
    cin>>t;
    while(t--) {
        auto n=int64_t{0};
        cin>>n;
        auto m=int64_t{0};
        cin>>m;
        auto p = vector<int64_t>{};
        for(auto i=int64_t{0}; i<n; ++i) {
            auto pv = int64_t{0};
            cin>>pv;
            p.push_back(pv-1);
        }

        auto a = vector<pair<int64_t, int64_t>>{};
        auto zero_count = int64_t{0};
        for(auto i=int64_t{0}; i<n; ++i) {
            auto av = int64_t{0};
            cin>>av;
            if(av == 0) {
                ++zero_count;
            }
            a.push_back(make_pair(av, i));
        }
        perm_transform(p, a);       

        // for(auto i=int64_t{0}; i<n; ++i) {
        //     cout << a[i].first << " ";
        // }
        // cout << endl;
        // for(auto i=int64_t{0}; i<n; ++i) {
        //     cout << a[i].second << " ";
        // }
        // cout << endl;

        auto distinct_zero_indexes = unordered_set<int64_t>{};
        auto distinct_numbers = unordered_set<int64_t>{};

        for(auto i=int64_t{0}; i<n; ++i) {
            if(a[i].first == 0) {
                distinct_zero_indexes.insert(a[i].second);
            } else {
                distinct_numbers.insert(a[i].first);
            }
        }

        if (distinct_numbers.size() <= 1) {
            auto max_power = zero_count;
            max_power -= distinct_zero_indexes.size();
            if(distinct_zero_indexes.size() != 0 && distinct_numbers.size() == 0) {
                max_power += 1;
            }

            auto result = int64_t(1);
            auto mod = static_cast<int64_t>(1e9 + 7);
            for(auto i=int64_t{0}; i<max_power; ++i) {
                result = ((result%mod) * (m%mod));
            }
            // cout << "Indexes: " << distinct_zero_indexes.size() << " Numbers: " << distinct_numbers.size() << " POW: " << max_power << endl;
            cout<<result%mod<<endl;
        } else {
            cout<<0<<endl;
        }
    }

    return 0;
}