#include <bits/stdc++.h>

using namespace std;


// too much memory!!!
// auto _find_sequence(vector<vector<int>> const & days, vector<int>& sequence, int day, unordered_set<int> const & seen) -> bool {
//     if(day < 0) {
//         return true;
//     }
//     auto day_to_check = days[day];
    
//     auto new_seen = unordered_set<int>(seen);
//     for(auto i=0; i<day_to_check.size(); ++i) {
//         new_seen.insert(day_to_check[i]);
//     }

//     for(auto i=0; i<day_to_check.size(); ++i) {
//         if (seen.count(day_to_check[i]) == 1) {
//             continue;
//         }
//         sequence[day] = day_to_check[i];
//         auto ok = _find_sequence(days, sequence, day-1, new_seen);
//         if (ok) {
//             return true;
//         }
//     }
//     return false;
// }

// auto find_sequence(vector<vector<int>> const & days, vector<int>& sequence) -> bool {
//     auto seen=unordered_set<int>{};
//     auto day=days.size()-1;
//     return _find_sequence(days, sequence, day, seen);
// }

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);

    auto t = 0;
    cin>>t;
    while(t--) {
        auto m=0;
        cin>>m;
        auto days = vector<unordered_set<int>>{};
        while(m--) {
            auto n=0;
            cin>>n;
            auto day = unordered_set<int>{};
            while(n--) {
                auto player = 0;
                cin>>player;
                day.insert(player);
            }
            days.push_back(std::move(day));
        }

        // prune
        auto cant_be_winner = unordered_set<int>(days[days.size()-1]);
        for (int i = days.size()-2; i>=0; --i) {
            for(auto element = days[i].begin(); element != days[i].end();) {
                if (cant_be_winner.count(*element) == 1) {
                    element = days[i].erase(element);
                } else {
                    cant_be_winner.insert(*element);
                    ++element;
                }
            }
        }

        auto winner_sequence = vector<int>(days.size(), 0);
        auto ok = true;
        for (auto i=0; i<days.size(); ++i) {
            if (days[i].empty()) {
                ok = false;
                break;
            }
            winner_sequence[i] = *(days[i].begin());
        }

        if (ok) {
            for (auto v = winner_sequence.begin(); v!= winner_sequence.end(); ++v) {
                cout<< *v << " ";
            }
            cout<<endl;
        } else {
            cout<<-1<<endl;
        }
    }

    return 0;
}