#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto line = string{};
    auto x = 1;
    auto cycle = 0;

    auto sum_at_cycle = set<int>{40, 80, 120, 160, 200, 240};
    auto total_signal_strength = i6{0};
    while(getline(cin, line)) {
        char operation[5];
        auto value = 0;
        sscanf(line.c_str(), "%s %d", operation, &value);

        auto to_execute = 0;
        if (strcmp(operation, "noop") == 0) {
            to_execute = 1;
        } else if (strcmp(operation, "addx") == 0) {
            to_execute = 2;
        }
        while(to_execute--) {
            ++cycle;
            if(abs(((cycle - 1) % 40) - x) <= 1) {
                cout << "#";

            } else {
                cout << ".";
            }
            if(sum_at_cycle.contains(cycle)) {
                total_signal_strength += x * cycle;
                cout<<endl;
            }
            if(strcmp(operation, "addx") == 0 && to_execute == 0) {
                x += value;
            }
        }
    }
    cout << total_signal_strength << endl;

    return 0;
}
