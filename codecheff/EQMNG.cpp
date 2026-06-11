#include <bits/stdc++.h>
#include <vector>

using i64 = int64_t;
using u64 = uint64_t;
using i32 = int32_t;
using u32 = uint32_t;
using u8 = uint8_t;

constexpr char brl = '\n';

using namespace std;

struct Tri
{
    i64 a;
    i64 b;
    i64 c;
};

auto solve() -> void
{
    i64 a, b, c;
    cin >> a >> b >> c;

    i64 cur_cost = 0;

    auto states = vector<Tri>();
    auto next_states = vector<Tri>();
    states.push_back({ a, b, c });
    auto found = false;

    while (!found)
    {
        for (auto const& state : states)
        {
            if (state.a == state.b)
            {
                found = true;
                break;
            }
            next_states.push_back({ gcd(state.a, state.c), gcd(state.b, state.c), state.c });
            next_states.push_back({ state.a, state.b, state.c + 1 });
        }
        swap(states, next_states);
        next_states.clear();
        ++cur_cost;
    }

    cout << cur_cost - 1 << brl;
}

auto main() -> int
{
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    auto t = i64{ 0 };
    std::cin >> t;
    while (t--) { solve(); }

    return 0;
}
