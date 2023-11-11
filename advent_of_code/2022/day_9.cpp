#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto direction = ' ';
  auto step = 0;
  auto line = string{};

  auto position_watcher = set<pair<int, int>>{};
  auto rope = vector<pair<int, int>>(10, make_pair(0, 0));

  while (getline(cin, line)) {
    sscanf(line.c_str(), "%c %d", &direction, &step);

    while (step--) {
      switch (direction) {
      case 'U': {
        ++rope[0].second;
        break;
      }
      case 'D': {
        --rope[0].second;
        break;
      }
      case 'R': {
        ++rope[0].first;
        break;
      }
      case 'L': {
        --rope[0].first;
        break;
      }
      }

      for (auto i = 1; i < static_cast<int>(rope.size()); ++i) {
        auto delta_x = rope[i - 1].first - rope[i].first;
        auto delta_y = rope[i - 1].second - rope[i].second;

        if (rope[i - 1].first == rope[i].first || rope[i - 1].second == rope[i].second) {
          if (abs(delta_x) == 2) { rope[i].first += delta_x / abs(delta_x); }
          if (abs(delta_y) == 2) { rope[i].second += delta_y / abs(delta_y); }
        } else if (abs(delta_x) != 1 || abs(delta_y) != 1) {
          rope[i].first += delta_x / abs(delta_x);
          rope[i].second += delta_y / abs(delta_y);
        }
      }
      auto [it, ok] = position_watcher.insert(make_pair(rope[9].first, rope[9].second));
      // if(ok) {
      //     cout << "x: " << rope[9].first << " y: " << rope[9].second << endl;
      // }
      cout << direction << endl;
      for (auto entry : rope) { cout << entry.first << " " << entry.second << ", "; }
      cout << endl;
      for (int y = 20; y >= -20; --y) {
        for (int x = -20; x < 20; ++x) {
          auto match =
            find_if(rope.begin(), rope.end(), [x, y](auto entry) { return entry.first == x && entry.second == y; });
          if (match != rope.end()) {
            auto dis = distance(rope.begin(), match);
            if (dis == 0) {
              cout << "H";
            } else {
              cout << dis;
            }
          } else {
            cout << ".";
          }
        }
        cout << endl;
      }
      cout << endl;
      this_thread::sleep_for(50ms);
    }
  }
  cout << position_watcher.size() << endl;
  for (int i = 14; i >= -15; --i) {
    for (int j = -15; j < 15; ++j) {
      if (position_watcher.count(make_pair(j, i))) {
        cout << "X";
      } else {
        cout << ".";
      }
    }
    cout << endl;
  }

  return 0;
}
