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
        auto a = vector<int>(n);
        for(auto i = 0; i < n; ++i) {
            cin>>a[i];
        }
        auto b = vector<int>(n);
        for(auto i = 0; i < n; ++i) {
            cin>>b[i];
        }
        ranges::sort(a);
        ranges::sort(b, 
            [](auto const& lhs, auto const& rhs) { 
                return lhs > rhs;
            });
        auto used = 0;
        auto total_order_count = i6{1};
        constexpr auto modulo = static_cast<i6>(1e9 + 7);
        for(auto i = 0; i < b.size(); ++i) {
            auto upper = upper_bound(a.begin(), a.end(), b[i]);
            auto count = std::distance(upper, a.end());
            auto multiple = count - used;
            ++used;
            if(multiple < 1) {
                total_order_count = 0;
                break;
            }
            total_order_count = (total_order_count * multiple) % modulo;
        }
        cout << total_order_count << endl;
    }

    return 0;
}
