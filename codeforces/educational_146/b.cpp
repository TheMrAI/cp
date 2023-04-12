#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(int a, int b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto a = 0;
        auto b = 0;
        cin>>a>>b;

        auto m = 1;
        auto power_of_two = 2;
        auto jump_count = 0;
        auto x = 0;
        auto y = 0;
        
        auto delta_x = a - x;
        auto delta_y = b - y;
        if(delta_x % 2 != 0) {
            x += m;
            ++jump_count;
        }
        if(delta_y % 2 != 0) {
            y += m;
            ++jump_count;
        }
        while(x != a || y != b) {
            delta_x = a - x;
            delta_y = b - y;
            auto to_next_power = power_of_two - m;

            auto steps_x_m = delta_x / m;
            auto steps_x_p = ceil2(delta_x, power_of_two) + to_next_power;
            auto steps_y_m = delta_y / m;
            auto steps_y_p = ceil2(delta_y, power_of_two) + to_next_power;
            
            if(steps_x_p + steps_y_p - to_next_power < steps_x_m + steps_y_m) {
                m += to_next_power;
                jump_count += to_next_power;
                power_of_two = power_of_two << 1;
                continue;
            } else {
                if (delta_x != 0 && steps_x_m <= steps_x_p) {
                    jump_count += delta_x / m;
                    x = a;
                    continue;
                }
                if (delta_y != 0 && steps_y_m <= steps_y_p) {
                    jump_count += delta_y / m;
                    y = b;
                    continue;
                }
            }
        }
        cout << jump_count << endl;
    }

    return 0;
}
