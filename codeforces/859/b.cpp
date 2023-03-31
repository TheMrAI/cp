#include <bits/stdc++.h>

using i64t = int64_t;
using namespace std;

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t=i64t{0};
    cin>>t;
    while(t--) {
        auto n=i64t{0};
        cin>>n;
        auto even_sum =i64t{0};
        auto odd_sum=i64t{0};
        while(n--){
            auto val=i64t{0};
            cin>>val;
            if(val % 2==0) {
                even_sum+=val;
            } else{
                odd_sum+=val;
            }
        }
        if(even_sum > odd_sum) {
            cout<<"YES"<<endl;
        } else {
            cout<<"NO"<<endl;
        }
    }

    return 0;
}
