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

auto to_key(int x, int y, int z) -> std::string {
    return std::to_string(x) + "," + std::to_string(y) + "," + std::to_string(z);
}

struct Triplet {
    int x;
    int y;
    int z;
};

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto line = string{};
    auto data_regex = regex(R"((\d+),(\d+),(\d+))");

    auto droplet_voxels = std::unordered_set<std::string>{};
    auto side_count = size_t{0};
    constexpr auto offsets = std::array<Triplet, 6>{Triplet{0, 0, 1}, Triplet{0, 0, -1}, Triplet{0, 1, 0}, Triplet{0, -1, 0}, Triplet{1, 0, 0}, Triplet{-1, 0, 0}};

    auto min_x = std::numeric_limits<int32_t>::max();
    auto max_x = std::numeric_limits<int32_t>::min();
    auto min_y = std::numeric_limits<int32_t>::max();
    auto max_y = std::numeric_limits<int32_t>::min();
    auto min_z = std::numeric_limits<int32_t>::max();
    auto max_z = std::numeric_limits<int32_t>::min();

    while(getline(cin, line)) {
        auto match = std::smatch{};
        if (std::regex_match(line, match, data_regex)) {
            std::cout << match[1].str() << ", " << match[2].str() << ", " << match[3].str() << std::endl;
        } else {
            std::cerr << "Noooo!" << std::endl;
            return -1;
        }

        auto x = std::stoi(match[1]);
        auto y = std::stoi(match[2]);
        auto z = std::stoi(match[3]);
        auto key = to_key(x,y,z);
        if (droplet_voxels.contains(key)) {
            continue;
        }
        // update location limits
        min_x = std::min(min_x, x);
        max_x = std::max(max_x, x);
        min_y = std::min(min_y, y);
        max_y = std::max(max_y, y);
        min_z = std::min(min_z, z);
        max_z = std::max(max_z, z);

        droplet_voxels.insert(key);
        side_count += 6;
        for (auto const& offset : offsets) {
            key = to_key(x + offset.x, y + offset.y, z + offset.z);
            if (droplet_voxels.contains(key)) {
                side_count -= 2;
            }
        }
    }
    // Part 1
    std::cout << "Droplet total side count: " << side_count << std::endl;
    // Part 2
    // extend bounding volume by 1 on each side so that it fully envelopes droplet
    --min_x;
    ++max_x;
    --min_y;
    ++max_y;
    --min_z;
    ++max_z;
    // std::cout << max_x << " " << min_x << " " << max_y << " " << min_y << " " << max_z << " " << min_z << std::endl;

    auto bounding_object_side_count = size_t{0};
    auto bounding_voxels = std::unordered_set<std::string>{};
    auto to_check = std::stack<Triplet>{};
    to_check.push(Triplet{min_x, min_y, min_z});
    
    while(!to_check.empty()) {
        auto checking = to_check.top();
        to_check.pop();
        // std::cout << checking.x << " " << checking.y << " " << checking.z << std::endl;
        auto checking_key = to_key(checking.x, checking.y, checking.z);
        if (bounding_voxels.contains(checking_key)) {
            continue;
        }
        bounding_voxels.insert(checking_key);
        bounding_object_side_count += 6;
        for (auto const& offset : offsets) {
            auto candidate_x = checking.x + offset.x;
            auto candidate_y = checking.y + offset.y;
            auto candidate_z = checking.z + offset.z;
            auto candidate_key = to_key(candidate_x, candidate_y, candidate_z);
            if (bounding_voxels.contains(candidate_key)) {
                bounding_object_side_count -= 2;
            }
            // skip all invalid voxels from side counting
            if (droplet_voxels.contains(candidate_key)
                || candidate_x < min_x || candidate_x > max_x
                || candidate_y < min_y || candidate_y > max_y
                || candidate_z < min_z || candidate_z > max_z) {
                continue;
            }
            to_check.push(Triplet{candidate_x, candidate_y, candidate_z});
        }
    }
    // Now from the total side count simply remove the outer perimeter which should leave us with the droplets
    // outside area
    // the plus one is necessary because we need the length of the sides not the distance between minimum and maximum
    auto delta_x = 1 + max_x - min_x;
    auto delta_y = 1 + max_y - min_y;
    auto delta_z = 1 + max_z - min_z;
    auto droplet_outside_perimeter = bounding_object_side_count - (2 * delta_x * delta_y) - (2 * delta_z * delta_y) - (2 * delta_x * delta_z);
    std::cout << "Droplet outer side count: " << droplet_outside_perimeter << std::endl;

    return 0;
}
