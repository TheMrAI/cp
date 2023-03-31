#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t=int64_t{0};
    cin>>t;
    while(t--) {
        auto n=i6{};
        cin>>n;
        auto s=string{};
        cin>>s;
        auto assigned = vector<int>(26, -1);
        auto val = true;
        auto can = true;
        for(auto i=size_t{0}; i<s.size(); ++i) {
            auto converted = static_cast<int>(val);
            if(assigned[s[i]-'a'] == -1) {
                assigned[s[i]-'a'] = converted;
            }
            if(assigned[s[i]-'a'] != converted) {
                can = false;
                break;
            }
            val = !val;
        }
        if(can) {
            cout << "YES"<<endl;
        } else {
            cout <<"NO"<<endl;
        }
    }

    return 0;
}
