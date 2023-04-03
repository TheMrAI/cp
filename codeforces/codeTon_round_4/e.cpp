#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto grow(vector<int>const& nodes, vector<vector<int>>const& edges, vector<bool>& already_reached, int from) -> bool {
    auto visited = vector<bool>(200000, false);

    auto to_check = set<pair<int, int>>{};
    to_check.insert(make_pair(nodes[from], from));

    auto slain = 0;

    while(!to_check.empty()) {
        auto checking = *to_check.begin();
        if(visited[checking.second]) {
            to_check.erase(to_check.begin());
            continue;
        }

        if (slain < checking.first) {
            return slain == nodes.size();
        }
        ++slain;
        to_check.erase(to_check.begin());

        for(auto neighbour : edges[checking.second]) {
            if(!visited[neighbour]) {
                to_check.insert(make_pair(nodes[neighbour], neighbour));
            }
        }
        visited[checking.second] = true;
        already_reached[checking.second] = true;
    }

    return slain == nodes.size();
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto v_count = i6{0};
        cin>>v_count;
        auto e_count = i6{0};
        cin>>e_count;

        auto nodes = vector<int>(v_count, 0);
        for(auto i=i6{0}; i<v_count; ++i) {
            auto monster = 0;
            cin>>monster;
            nodes[i] = monster;
        }

        auto edges = vector<vector<int>>(v_count);
        for(auto i=i6{0}; i<e_count; ++i) {
            auto a = i6{0};
            auto b = i6{0};
            cin>>a>>b;
            a-=1;
            b-=1;
            edges[a].push_back(b);
            edges[b].push_back(a);
        }

        auto solvable = false;
        auto already_reached = vector<bool>(200000, false);
        for (auto i=0; i<v_count; ++i) {
            if(nodes[i] == 0 && !already_reached[i]) {
                solvable = grow(nodes, edges, already_reached, i);
                if(solvable) {
                    break;
                }
            }
        }

        if(solvable) {
            cout<<"YES"<<endl;
        } else {
            cout<<"NO"<<endl;
        }
    }

    return 0;
}
