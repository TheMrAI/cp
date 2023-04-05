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

    auto row = string("");
    auto forest = vector<vector<int>>{};
    while(getline(cin, row)) {
        auto tree_heights = vector<int>(row.size(), 0);
        for(auto i=0; i < row.size(); ++i) {
            tree_heights[i] = row[i]-'0';
        }
        forest.push_back(move(tree_heights));
    }

    auto left_heights = vector<vector<int>>(forest.size(), vector<int>(forest[0].size(), -1));
    for (auto i=0; i < forest.size(); ++i) {
        for(auto j=1; j < forest[0].size(); ++j) {
            left_heights[i][j] = max(left_heights[i][j-1], forest[i][j-1]);
        }
    }
    auto right_heights = vector<vector<int>>(forest.size(), vector<int>(forest[0].size(), -1));
    for (auto i=0; i < forest.size(); ++i) {
        for(auto j=static_cast<int>(forest[0].size())-2; j >= 0 ; --j) {
            right_heights[i][j] = max(right_heights[i][j+1], forest[i][j+1]);
        }
    }
    auto up_heights = vector<vector<int>>(forest.size(), vector<int>(forest[0].size(), -1));
    for(auto j=0; j < forest[0].size(); ++j) {
        for (auto i=1; i < forest.size(); ++i) {
            up_heights[i][j] = max(up_heights[i-1][j], forest[i-1][j]);
        }
    }
    auto down_heights = vector<vector<int>>(forest.size(), vector<int>(forest[0].size(), -1));
    for(auto j=0; j < forest[0].size(); ++j) {
        for (auto i=static_cast<int>(forest.size()) - 2; i >= 0; --i) {
            down_heights[i][j] = max(down_heights[i+1][j], forest[i+1][j]);
        }
    }

    auto visible = 0;
    for (auto i = 0; i < forest.size(); ++i) {
        for(auto j = 0; j < forest[0].size(); ++j) {
            // cout << forest[i][j];
            if (left_heights[i][j] < forest[i][j] 
            || right_heights[i][j] < forest[i][j] 
            || up_heights[i][j] < forest[i][j] 
            || down_heights[i][j] < forest[i][j]) {
                ++visible;
            }
        }
        // cout << endl;
    }
    cout << "Visible trees: " << visible << endl;

    auto max_scenic = 0;
    for (auto i = 0; i < forest.size(); ++i) {
        for(auto j = 0; j < forest[0].size(); ++j) {
            auto to_left = 0;
            if (left_heights[i][j] < forest[i][j]) {
                to_left = j;
            } else {
                for (auto k = j - 1; k >= 0; --k) {
                    ++to_left;
                    if (forest[i][k] >= forest[i][j]) {
                        break;
                    }
                }
            }
            auto to_right = 0;
            if (right_heights[i][j] < forest[i][j]) {
                to_right = forest[0].size() - j - 1;
            } else {
                for (auto k = j + 1; k < forest[0].size(); ++k) {
                    ++to_right;
                    if (forest[i][k] >= forest[i][j]) {
                        break;
                    }
                }
            }
            auto to_up = 0;
            if (up_heights[i][j] < forest[i][j]) {
                to_up = i;
            } else {
                for (auto k = i - 1; k >= 0; --k) {
                    ++to_up;
                    if (forest[k][j] >= forest[i][j]) {
                        break;
                    }
                }
            }
            auto to_down = 0;
            if (down_heights[i][j] < forest[i][j]) {
                to_down = forest.size() - i - 1;
            } else {
                for (auto k = i + 1; k < forest.size(); ++k) {
                    ++to_down;
                    if (forest[k][j] >= forest[i][j]) {
                        break;
                    }
                }
            }
            max_scenic = max(max_scenic, to_left*to_right*to_up*to_down);
        }
    }
    cout << "Scenic score: " << max_scenic << endl;

    return 0;
}
