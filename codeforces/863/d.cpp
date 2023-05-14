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

auto valid_fn(array<int64_t, 46> const& fibs, int64_t n, int64_t x, int64_t y) -> bool {
    if(n == 1) {
        return true;
    }
    auto height = fibs[n];
    auto width = fibs[n+1];
    // loop check while adjusting the coordinates and rotating them 90 degrees
    // the rotation direction doesn't matter as the problem is centrally symmetric
    if(y > height) {
        return valid_fn(fibs, n-1, y - height, x);
    }
    if(width - y >= height) {
        return valid_fn(fibs, n-1, y, x);
    }
    return false;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto fib_array = array<int64_t, 46>{};
    fib_array[0] = 1;
    fib_array[1] = 1;
    for(auto i = 2; i < 46; ++i) {
        fib_array[i] = fib_array[i-1] + fib_array[i-2];
    }

    auto t=i6{0};
    cin>>t;
    while(t--) {
        int64_t n = 0;
        cin>>n;
        int64_t x = 0;
        int64_t y = 0;
        cin>>x>>y;
        if(valid_fn(fib_array, n, x, y)) {
            cout << "YES" << endl;
            continue;
        }
        cout << "NO" << endl;
    }

    return 0;
}
