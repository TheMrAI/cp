#include <bits/stdc++.h>

using namespace std;

auto check(vector<int>& a, vector<int>& b) -> bool {
    auto size = static_cast<int>(a.size());
    for(auto i=size-2; i>=0; --i) {
            if(a[i] > a[size-1] || b[i] > b[size-1]) {
                swap(a[i], b[i]);
                if(a[i] > a[size-1] || b[i] > b[size-1]) {
                    return false;
                }
            }
    }
    return true;
}

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t = 0;
    cin>>t;
    while(t--) {
        auto n =0;
        cin>>n;

        auto a = vector<int>{};
        for(auto i=0; i<n; ++i) {
            auto value=0;
            cin>>value;
            a.push_back(value);
        }

        auto b = vector<int>{};
        for(auto i=0; i<n; ++i) {
            auto value=0;
            cin>>value;
            b.push_back(value);
        }

        if (n==1) {
            cout<<"Yes"<<endl;
            continue;
        }

        if (check(a, b)) {
            cout<<"Yes"<<endl;
        } else {
            cout<<"No"<<endl;
        }
    }

    return 0;
}