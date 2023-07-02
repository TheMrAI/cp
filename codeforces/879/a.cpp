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
        cin>>n;
        auto pos_count = 0;
        auto neg_count = 0;
        while(n--) {
            auto val = 0;
            cin>>val;
            if(val == -1) {
                ++neg_count;
            } else {
                ++pos_count;
            }
        }
        
        auto mod_count = 0;
        if(neg_count % 2 != 0) {
            --neg_count;
            ++pos_count;
            mod_count = 1;
        }
        if(pos_count >= neg_count) {
            cout << mod_count << endl;
        } else {
            auto diff = ceil2(neg_count - pos_count, 2);
            if (diff % 2 != 0) {
                ++diff;
            }
            cout << mod_count + diff << endl;
        }
    }

    return 0;
}
