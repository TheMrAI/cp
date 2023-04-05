#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto get_ring(int x, int y, int n) -> int {
    auto ring_count = n/2;
    x -= 1;
    y -= 1;
    auto transformed_x = x >= ring_count ? (ring_count-1) - (x%ring_count) : x;
    auto transformed_y = y >= ring_count ? (ring_count-1) - (y%ring_count) : y;
    auto ring = min(transformed_x, transformed_y);
    return ring;
}   

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n=0;
        auto x1=0;
        auto y1=0;
        auto x2=0;
        auto y2=0;
        cin>>n>>x1>>y1>>x2>>y2;

        auto ring_1 = get_ring(x1, y1, n);
        auto ring_2 = get_ring(x2, y2, n);
        if (ring_2 > ring_1) {
            swap(ring_1, ring_2);
        }

        cout << ring_1 - ring_2 << endl;
    }

    return 0;
}
