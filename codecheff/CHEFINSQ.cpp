#include <bits/stdc++.h>
#include <cstdint>
using namespace std;

#define bl '\n'

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto binom = vector<vector<uint64_t>>(60, vector<uint64_t>(60));
	binom[0][0] = 1;
	for (int i = 1; i < 60; ++i)
	{
		binom[i][0] = 1;
		for (int j = 1; j < 60; ++j)
		{
			binom[i][j] = binom[i - 1][j - 1] + binom[i - 1][j];
		}
	}

    int t;
    cin >> t;
    while(t--) {
        int n, k;
        cin >> n >> k;

        auto nums = vector<int>{};
        while(n--) {
            int val;
            cin >> val;
            nums.push_back(val);
        }

        sort(nums.begin(), nums.end());

        int biggest = nums[k-1];
        int total_count = 0;
        int need_count = 0;
        for (auto i = 0; i < nums.size(); ++i) {
            if (nums[i] == biggest) {
                if (i < k) {
                    ++need_count;
                }
                ++total_count;
            }
        }
        cout << binom[total_count][need_count] << bl;
    }
}
