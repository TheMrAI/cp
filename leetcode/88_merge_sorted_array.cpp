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

class Solution {
public:
    void merge(vector<int>& nums1, int m, vector<int>& nums2, int n) {
        if(nums2.empty()) {
            return;
        }
        auto i = m-1;
        auto j = n-1;
        auto write = m+n - 1;
        while(j >= 0 && i >= 0) {
            if(nums1[i] <= nums2[j]) {
                nums1[write] = nums2[j];
                --j;
            } else {
                nums1[write] = nums1[i];
                --i;
            }
            --write;
        }
        while(j >= 0) {
            nums1[write] = nums2[j];
            --write;
            --j;
        }
    }
};

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    return 0;
}
