#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n=i6{0};
        auto c=i6{0};
        auto d=i6{0};
        cin>>n;
        cin>>c;
        cin>>d;

        auto numbers = vector<i6>(n+1, 0);
        for(auto i=1; i<n+1; ++i) {
            auto val=i6{0};
            cin>>val;
            numbers[i] = val;
        }
        sort(numbers.begin(), numbers.end());
        auto last = unique(numbers.begin(), numbers.end());
        numbers.erase(last, numbers.end());

        auto sequence_until = vector<i6>(numbers.size());
        sequence_until[0] = c*n; // removed all duplicates cost
        for(auto i=1; i < numbers.size(); ++i) {
            auto diff = numbers[i] - numbers[i-1];
            sequence_until[i] = sequence_until[i-1] + (d*(diff-1)) - c;
        }
        auto min = min_element(std::next(sequence_until.begin(), 1), sequence_until.end());
        auto full_rebuild = n*c + d;
        if (*min < full_rebuild) {
            cout<<*min<<endl;
        } else {
            cout << full_rebuild <<endl;
        }
    }

    return 0;
}
