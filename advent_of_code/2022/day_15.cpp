#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

struct LineSegments{
    int64_t x_s; // horizontal
    int64_t x_u;
    int64_t y_s; // vertical
    int64_t y_u;
    int64_t center_x;
    int64_t center_y;
    int64_t beacon_x;
    int64_t beacon_y;
    int64_t manhattan_distance;
};

auto within_beacon_range(LineSegments const& segment, int64_t x, int64_t y) -> bool 
{
    return std::abs(segment.center_x - x) + std::abs(segment.center_y - y) <= segment.manhattan_distance;
}

auto locate_beacon(vector<LineSegments> const& segments, int64_t l_x, int64_t h_x, int64_t l_y, int64_t h_y, pair<int64_t, int64_t>& beacon_location) -> bool {
    if(l_x == h_x || l_y == h_y) {
        return false;
    }

    auto all_corners_covered_by_one_segment = false;
    auto corners = vector<pair<int64_t, int64_t>>{{l_x, l_y}, {l_x, h_y - 1}, {h_x - 1, l_y}, {h_x - 1, h_y - 1}};
    for(auto const& segment : segments) {
        all_corners_covered_by_one_segment = std::all_of(corners.cbegin(), corners.cend(),
            [segment](auto const& corner) {
                return within_beacon_range(segment, corner.first, corner.second);
            });
        if(all_corners_covered_by_one_segment) {
            break;
        }
    }

    if(all_corners_covered_by_one_segment) {
        return false;
    }

    if(h_x - l_x == 1 && h_y - l_y == 1) {
        beacon_location.first = l_x;
        beacon_location.second = l_y;
        return true;
    }

    auto mid_x = l_x + (h_x - l_x) / 2;
    auto mid_y = l_y + (h_y - l_y) / 2;
    return locate_beacon(segments, l_x, mid_x, l_y, mid_y, beacon_location) 
        || locate_beacon(segments, mid_x, h_x, l_y, mid_y, beacon_location)
        || locate_beacon(segments, l_x, mid_x, mid_y, h_y, beacon_location) 
        || locate_beacon(segments, mid_x, h_x, mid_y, h_y, beacon_location);
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    
    // constexpr auto target_y = i6{10};
    constexpr auto target_y = i6{2'000'000};

    auto line = string{};
    auto line_segments = vector<LineSegments>{};

    auto maximum_x = i6{0};
    auto maximum_y = i6{0};

    auto data_regex = regex(R"((-?\d+))");
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
        maximum_x = std::max(maximum_x, sensor_x);
        maximum_y = std::max(maximum_y, sensor_y);

        line_segments.emplace_back(sensor_x - manhattan_distance, sensor_x + manhattan_distance, sensor_y - manhattan_distance, sensor_y + manhattan_distance, sensor_x, sensor_y, beacon_x, beacon_y, manhattan_distance);    
    }
  
    ranges::sort(line_segments.begin(), line_segments.end(), 
        [](auto const& lhs, auto const& rhs) {
            return lhs.y_s < rhs.y_s;
        });
    auto candidates_at_line_y = vector<LineSegments>{};
    std::copy_if(line_segments.cbegin(), line_segments.cend(), std::back_inserter(candidates_at_line_y),
        [target_y](auto const& segment){
            return segment.y_s <= target_y && target_y <= segment.y_u;
        });
    ranges::sort(candidates_at_line_y.begin(), candidates_at_line_y.end(), 
        [](auto const& lhs, auto const& rhs){
            return lhs.x_s + std::abs(target_y - lhs.center_y) < rhs.x_s + std::abs(target_y - rhs.center_y);
        });

    auto covered = i6{0};
    auto covered_until = std::numeric_limits<i6>::min();
    auto beacons_on_the_scanline = set<int>{};
    for (auto const& segment : candidates_at_line_y) {
        auto segment_start = segment.x_s + std::abs(target_y - segment.center_y);
        auto segment_until = segment.x_u - std::abs(target_y - segment.center_y);

        if(segment_until <= covered_until) {
            continue;
        }
        
        auto from = std::max(covered_until + 1, segment_start);
        covered += 1 + segment_until - from;
        covered_until = segment_until;

        if(segment.beacon_y == target_y && segment_start <= segment.beacon_x && segment.beacon_x <= segment_until) {
            beacons_on_the_scanline.emplace(segment.beacon_x);
        }
    }
    covered -= beacons_on_the_scanline.size();
    cout << "# of positions where the beacon cannot be at: " << covered << endl;
    
    auto beacon_position = make_pair(i6{-1}, i6{-1});
    locate_beacon(line_segments, 0, 4'000'001, 0, 4'000'001, beacon_position);
    cout << "Beacon position: " << (beacon_position.first * 4'000'000) + beacon_position.second << endl;

    return 0;
}
