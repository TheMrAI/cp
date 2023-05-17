#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
requires 
    std::derived_from<Output, std::ostream> &&
    std::weakly_incrementable<Iter> && 
    std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end) {
    while(begin != end) {
        output << *begin;
        if(std::next(begin) != end) {
            output << ' ';
        }
        ++begin;
    }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T {
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
        auto k = 0;
        cin>>k>>n;
        auto result = vector<int>(k, 0);
        iota(result.begin(), result.end(), 1);
        if(k == 2 || k == n) {
            dump_to(cout, result.begin(), result.end());
            cout << endl;
            continue;
        }
        auto diff = 1;
        for(auto i = 2; i < result.size(); ++i) {
            auto next_best = result[i-1] + diff + 1;
            if(k - i - 1 <= n - (next_best)) {
                result[i] = next_best;
                ++diff;
                continue;
            }
            diff = 1;
            result[i] = result[i-1]+1;
        }
        dump_to(cout, result.begin(), result.end());
        cout << endl;
    }

    return 0;
}
