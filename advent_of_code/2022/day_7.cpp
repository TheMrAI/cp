#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

struct FileOrDir {
    bool dir;
    int size_or_index;
    string name;
};


// Part 1
auto calculate_total_sum(vector<vector<FileOrDir>>const& filesystem, int in_dir, uint64_t& total_dir_sum) -> uint64_t {
    auto local_dir_sum = uint64_t{0};
    for(auto const& entry : filesystem[in_dir]) {
        if (entry.dir) {
            local_dir_sum += calculate_total_sum(filesystem, entry.size_or_index, total_dir_sum);
        } else {
            local_dir_sum += entry.size_or_index;
        }
    }

    if(local_dir_sum <= 100'000) {
        total_dir_sum += local_dir_sum;
    }

    return local_dir_sum;
}

// Part 2
auto find_smallest_dir_for_delete(vector<vector<FileOrDir>>const& filesystem, int in_dir, uint64_t space_needed, uint64_t& smallest_sufficient_dir_size) -> uint64_t {
    auto local_dir_sum = uint64_t{0};
    for(auto const& entry : filesystem[in_dir]) {
        if (entry.dir) {
            local_dir_sum += find_smallest_dir_for_delete(filesystem, entry.size_or_index, space_needed, smallest_sufficient_dir_size);
        } else {
            local_dir_sum += entry.size_or_index;
        }
    }

    if(space_needed <= local_dir_sum) {
        smallest_sufficient_dir_size = min(smallest_sufficient_dir_size, local_dir_sum);
    }

    return local_dir_sum;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto cd_regex = regex{R"(^\$+ +cd +(.+))"};
    auto path_regex = regex{R"((\w|\.\.|\.)+)"};
    auto ls_regex = regex{R"(^\$+ +ls)"};
    auto dir_regex = regex{R"(^dir +(.+))"};
    auto file_regex = regex{R"(^(\d*) +(.*))"};
    auto match_result = smatch{};

    auto filesystem = vector<vector<FileOrDir>>{};
    filesystem.push_back({}); // put in root at index 0

    auto location = stack<int>{};
    location.push(0);

    for (auto line = string{""}; getline(cin, line);) {
        if(regex_match(line, match_result, cd_regex)) {
            auto path = match_result[1].str();
            if (path == string("/")) {
                location = stack<int>{};
                location.push(0);
                continue;
            }

            for (auto entry = sregex_iterator(path.begin(), path.end(), path_regex); entry != sregex_iterator{}; ++entry) {
                auto move_to = entry->str();
                if (move_to == string{".."}) {
                    location.pop();
                } else if (move_to == string{"."}) {
                    continue;
                } else {
                    auto dir = find_if(filesystem[location.top()].begin(), filesystem[location.top()].end(), 
                     [&move_to](auto val) {
                        return val.dir == true && val.name == move_to;
                    });
                    location.push(dir->size_or_index);
                }
            }
        } else if(regex_match(line, match_result, ls_regex)) {
            continue;
        } else if (regex_match(line, match_result, dir_regex)) {
            filesystem.push_back({});
            filesystem[location.top()].push_back(FileOrDir{true, static_cast<int>(filesystem.size()) - 1, match_result[1].str()});
        } else if (regex_match(line, match_result, file_regex)) {
            auto name = match_result[2].str();
            auto size = stoi(match_result[1].str());
            filesystem[location.top()].push_back(FileOrDir{false, size, name});
        }
    }

    auto total_sum = uint64_t{0};
    auto used_space = calculate_total_sum(filesystem, 0, total_sum);
    cout<<"Used space: " << used_space << endl;
    cout<<"Part 1: "<<total_sum<<endl;
    auto available_space = uint64_t{70'000'000};
    auto free_space = available_space - used_space;
    auto required_space = free_space > 30'000'000 ? 0 : 30'000'000 - free_space;
    cout << "Free: " << free_space << " Required: " << required_space << endl;
    auto smallest_sufficient_dir_size = available_space;
    find_smallest_dir_for_delete(filesystem, 0, required_space, smallest_sufficient_dir_size);
    cout<<"Part 2: "<<smallest_sufficient_dir_size<<endl;

    return 0;
}
