#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

struct Entry{
    int64_t sensor_x;
    int64_t sensor_y;
    int64_t beacon_x;
    int64_t beacon_y;
    int64_t manhattan_distance;
};

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto line = string{};
    auto entries = vector<Entry>{};
    auto data_regex = regex(R"((-?\d+))");

    auto min_x = numeric_limits<int64_t>::max();
    auto max_x = numeric_limits<int64_t>::min();
    auto beacons_or_sensor_in_row_at_x = map<int64_t,set<int64_t>>{};
    
    auto target_y = i6{2000000};
    while(getline(cin, line)) { 
        auto numbers_begin = std::sregex_iterator(line.begin(), line.end(), data_regex);
        auto numbers_end = std::sregex_iterator();
        auto numbers = vector<int64_t>{};
        for(auto i = numbers_begin; i != numbers_end; ++i) {
            numbers.push_back(atoll(i->str().c_str()));
        }
        int64_t sensor_x = numbers[0];
        int64_t sensor_y = numbers[1];
        int64_t beacon_x = numbers[2];
        int64_t beacon_y = numbers[3];
        int64_t manhattan_distance = abs(beacon_x - sensor_x) + abs(beacon_y - sensor_y);
        auto entry = Entry{sensor_x, sensor_y, beacon_x, beacon_y, manhattan_distance};
        entries.push_back(move(entry));
        min_x = min({min_x, sensor_x - manhattan_distance, beacon_x});
        max_x = max({max_x, sensor_x + manhattan_distance, beacon_x});
            
        beacons_or_sensor_in_row_at_x[sensor_y].insert(sensor_x);
        beacons_or_sensor_in_row_at_x[beacon_y].insert(beacon_x);
    }
  
    sort(entries.begin(), entries.end(), 
        [target_y](auto const& lhs, auto const& rhs){
            auto lhs_delta_y = abs(target_y - lhs.sensor_y); 
            auto rhs_delta_y = abs(target_y - rhs.sensor_y);
            return lhs_delta_y < rhs_delta_y;
        });

    auto no_beacon_pos = i6{0};
    for(auto x = min_x; x <= max_x; ++x) {
        if(beacons_or_sensor_in_row_at_x[target_y].count(x) == 1) {
            continue;
        }
        for (auto entry = entries.begin(); entry != entries.end(); ++entry) {
            auto delta_y = abs(target_y - entry->sensor_y);

            auto covered_x_l = max(min_x, entry->sensor_x - (entry->manhattan_distance - delta_y));
            auto covered_x_r = min(max_x, entry->sensor_x + (entry->manhattan_distance - delta_y));
            if(covered_x_l <= x && x <= covered_x_r) {
                ++no_beacon_pos;
                break;
            }
        }
    }
    cout << no_beacon_pos << endl;

    auto candidate_points = vector<pair<Entry, Entry>>{};
    for (auto i = 0; i < entries.size() - 1; ++i) {
        for(auto j= i + 1; j < entries.size(); ++j) {
            auto sensors_manhattan_distance = abs(entries[i].sensor_x - entries[j].sensor_x) + abs(entries[i].sensor_y - entries[j].sensor_y);
            if(sensors_manhattan_distance == entries[i].manhattan_distance + entries[j].manhattan_distance + 2) {
                cout << "Candidate" << endl;
                candidate_points.push_back(make_pair(entries[i] , entries[j]));
            }
        }
    }
    for (auto const& can : candidate_points) {
        cout << can.first.sensor_x << " " << can.first.sensor_y << " - " << can.second.sensor_x << " " << can.second.sensor_y << endl;
    }

    return 0;
}
