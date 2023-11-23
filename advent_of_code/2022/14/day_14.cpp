#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto parse_line(regex const& position_regex, string_view line, int& max_x, int& max_y) -> vector<pair<int, int>>
{
  auto positions = vector<pair<int, int>>{};

  auto positions_begin = cregex_iterator(line.cbegin(), line.cend(), position_regex);
  auto positions_end = cregex_iterator();

  for (auto position = positions_begin; position != positions_end; ++position) {
    auto x = stoi((*position)[1].str());
    auto y = stoi((*position)[2].str());
    max_x = max(max_x, x);
    max_y = max(max_y, y);
    positions.emplace_back(x, y);
  }

  return positions;
}

auto print_cave(vector<vector<char>> const& cave) -> void
{
  for (auto y = 0; y < cave.size(); ++y) {
    for (auto x = 0; x < cave[0].size(); ++x) { cout << cave[y][x]; }
    cout << endl;
  }
}

auto simulate(vector<vector<char>>& cave, bool print) -> int
{
  auto grain_settlement_count = 0;

  while (true) {
    auto grain_y = 0;
    auto grain_x = 500;
    if (cave[grain_y][grain_x] == 'o' || cave[grain_y][grain_x] == 'x') { return grain_settlement_count; }

    while (grain_y != cave.size() || grain_x != cave[0].size()) {
      if (print) {
        cave[grain_y][grain_x] = '~';
        print_cave(cave);
        cout << endl;
        cave[grain_y][grain_x] = '.';
        // this_thread::sleep_for(30ms);
      }

      if (grain_y < cave.size() - 1) {
        if (cave[grain_y + 1][grain_x] == '.') {
          ++grain_y;
          continue;
        }
        if (grain_x > 0) {
          if (cave[grain_y + 1][grain_x - 1] == '.') {
            ++grain_y;
            --grain_x;
            continue;
          } else if (grain_x < cave[0].size() - 1) {
            if (cave[grain_y + 1][grain_x + 1] == '.') {
              ++grain_y;
              ++grain_x;
              continue;
            } else {
              cave[grain_y][grain_x] = 'o';
              ++grain_settlement_count;
              break;
            }
          } else {
            return grain_settlement_count;
          }
        } else {
          return grain_settlement_count;
        }
      } else {
        return grain_settlement_count;
      }
    }
  }

  return grain_settlement_count;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto position_regex = regex(R"((\d+),(\d+))");
  auto line = string{ "" };
  auto max_x = 0;
  auto max_y = 0;

  auto lines = vector<vector<pair<int, int>>>{};
  while (getline(cin, line)) {
    auto line_segments = parse_line(position_regex, line, max_x, max_y);
    lines.push_back(move(line_segments));
  }

  for (auto& line : lines) {
    for (auto segment : line) { cout << segment.first << "," << segment.second << " "; }
    cout << endl;
  }

  auto cave = vector<vector<char>>(max_y + 3, vector<char>(2 * max_x + 1, '.'));
  for (auto& entry : cave[max_y + 2]) { entry = 'x'; }
  for (auto const& line : lines) {
    for (auto i = 1; i < line.size(); ++i) {
      auto x = line[i - 1].first;
      auto y = line[i - 1].second;
      auto last_x = line[i].first;
      auto last_y = line[i].second;

      auto offset_x = (last_x - x) / abs(last_x - x);
      auto offset_y = (last_y - y) / abs(last_y - y);

      while (x != last_x || y != last_y) {
        cave[y][x] = 'x';

        if (x != last_x) { x += offset_x; }
        if (y != last_y) { y += offset_y; }
      }
      cave[y][x] = 'x';
    }
  }
  lines.clear();

  cout << simulate(cave, false) << endl;

  return 0;
}
