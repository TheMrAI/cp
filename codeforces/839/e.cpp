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
        auto data = vector<int>(n);
        for(auto i = 0; i < n; ++i) {
            cin>>data[i];
        }
        auto first_count = 0;
        auto second_count = 0;
        auto both_count = 0;
        for(auto i = 0; i < n; ++i) {
            auto first_needs = i + 1 != data[i] ? true : false;
            auto second_needs = n - i != data[i] ? true : false;
            if(first_needs) {
                ++first_count;
            }
            if(second_needs) {
                ++second_count;
            }
            if(first_needs && second_needs) {
                ++both_count;
            }
        }
        if(first_count <= second_count - both_count) {
            cout << "First" << endl;
        } else if (second_count < first_count - both_count) {
            cout << "Second" << endl;
        } else {
            cout << "Tie" << endl;
        }
    }

    return 0;
}
