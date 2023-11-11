#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
requires std::derived_from<Output, std::ostream> && std::weakly_incrementable<Iter> && std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end)
{
  while (begin != end) {
    output << *begin;
    if (std::next(begin) != end) { output << ' '; }
    ++begin;
  }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

const auto shape_a = std::vector<std::string>{ "####" };

const auto shape_b = std::vector<std::string>{ ".#.", "###", ".#." };

// this shape is mirrored horizontally
const auto shape_c = std::vector<std::string>{ "###", "..#", "..#" };

const auto shape_d = std::vector<std::string>{ "#", "#", "#", "#" };

const auto shape_e = std::vector<std::string>{
  "##",
  "##",
};

const auto shapes = std::array<std::vector<std::string>, 5>{ shape_a, shape_b, shape_c, shape_d, shape_e };
// const auto jets = std::string{">>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>"};
const auto jets = std::string{
  ">><<<>>>><<>>><<>><>>>><<<><<><>>><<<><<<>>><<>>><<<>>><>>>><<<>>>><<<><<>>><>>>><<>>>><<>>><><>>>><<<>><<<>>>><<>>>"
  "><<<>>><<>>><>>><<<>><<<<>><<<<>>>><<<<>>>><<><<<>>>><<><>>>><<><>>><<<><<><>>>><<<>>><>>>><<<>>><<>>>><<<><<<>>><><"
  ">><>>>><<>>><<<><>>>><<>><<<>>>><<>><>>><<>><<>>>><<>><<<>><<>>><<>>>><><<>><<<<><<>>><<<<>><>><>>>><<><<<>><<<<>>><"
  ">>><<<<>>><><<<>>><<>><>>>><<>>>><<<>>>><<>><<>>>><<<>><<>><<<><<<<>>>><<<><<<><>>><><<<<>>>><><<<<><<<<>><>>><<>><<"
  "<<><<<<><<><<<<>><<<<>><<>><<<>><<<>><>><<<>>><<>><<<<>>>><<<<>>><<>>>><<<>><>>>><<<<>>>><<<<>>>><<>><>><<>><<<<>>><"
  "<<>>><<<<>>><<<>>>><<>>>><<>>><>><<>><<<><><><<<<>><><<<<>><<>><<>>><<<<>>>><>>><<>>>><<<>>><<<>>><<<><><>>><<<<>>><"
  ">>><<<>>><><<>>>><<<<>>>><<<<><>>><<>>>><<<<>><<<>>><><<<<>>><<<>><<>>>><<>>><<<>>><<<>>><<>>><<>>>><<<<>><<><>>><<<"
  "<>>><<<>>><>>>><<<>>><<<>>><<<<>>>><<<<>><<<>>>><<<><<<<><<>><<<>>><<<>>>><<<>><<<<>>><<<<><>>>><<<>>><><<<>>>><<><<"
  "<><<<<>>>><<>><>>>><<<>>><<<>>><<>>>><>>><><<<<>>>><<<>>><<<<>><>>>><>><>>>><>>>><<<<>>>><<>>><<>><<>><<<><>><<<>>><"
  "<>>>><<<>><<><<<<>>>><<<<><><<>>><>><><<>>><<>>>><<><<<>><<<>>><>>>><<<><<>>>><<>>>><<>>><<<>>>><>>><<<<>><<<<><<<<>"
  "<<<<>><<<<><<<>>>><>><<<>><<><<>><<<>><>><<<<>>><>><<<<>><<<<>>>><<<>><<<<>>>><<<>><<>>>><<<<>>><>>>><<>>>><<<<>>><<"
  ">>>><>>><<<<>>>><<<>>><<<><<<>><><>><<<<>>>><<<<>><<<><<<>>>><<>>><<><<>>><<>><<<>>><<<>>><<<>>>><>>>><<<<>>>><<<>><"
  "<<<><<<>><<>>><>>><<<<>>><<<>><>>><<<><<<<><<>><<<<><>><<<>><<><<<>>><<<<>>><<<>>><<><>><<<<><><<<<>>>><<>>><<<<>>><"
  ">><<<<><<>><<<<>>><<>>><<<<>><<>>><<<<>><>>>><<>><<<<>>>><<>>>><<<<>><<<>><<>>>><<<>>>><>>>><>>>><<<><<>>><<>>><<<<>"
  ">><>><><<<<>>>><<<><<<>>><>>><<>><<>>><<>>><<<>>><<<<>><<>>>><<<><<>>>><>><<<>>>><>>>><<<>>><<><<<<>><<<>>><<>>><<<<"
  ">>><<<<><<<<>>><<<>>><>>><<>><<<<>>><>><<><<<<>>><<>>><>>><<>>><><<>>>><<<<>>><>>><>>><>><>>><<<<>><<<><<>><<<<>>><>"
  ">>><<<<>><<<<><<<<>>>><<>>><<<>>>><<<>><>>><<<<>>><<>>>><<>>><<>>>><>>>><<><<<>>><<<>><<>><<<<>><>>>><<>>>><<><<<>>>"
  "<>>>><<<<><<<>><>>><<<<>><>><<>>>><<<>><<<>><>>>><><>>>><>>><<<<>>>><<<>>>><<<>>><><<<<><>>><<<<>>>><<<><<><<<><<>>>"
  "><><>>>><<<<>><<<<>>><><><<<>><<>>>><<<><<<<>>><<<<><<<<>>><<<>>><<><<>>><<<>>><<<>>><<>>><<>><>><<<<>>><<<>>><<<>>>"
  "<>>><<>><<>>><<<<>>>><<><<>>>><<><<<<>>><<<><<<>>><><<<>>><<>>><<<>>>><<<<>>>><><<>>><<<>>>><<<>>>><<><<>>>><<<>><<>"
  "><<>>>><<>><<<><><<<><>>><><>>><<>>>><<<>><>>><>><>>>><<<<>>><>>><>><<<><<>>>><<<>>>><<><><<<>><<>><<<><<<><<><<<<>>"
  "><>><<<<>>><<<>>><<>>><<>>><<><<<<>>>><<>>><>><<<<>><<>>><>><<>><<<><>>><<><<<>>>><<<><>>>><<<<>><<<<><>><>><><<<>><"
  "<<>>>><><<<><<>><<>><>>><<<>>>><<>><<<<>>>><>>>><<>><<<>><<<>><<<<><<<<>>><<>>>><<<<>><<>>>><<<<>>><<<>><<>><<>><<<<"
  "><>>>><<<<><<<<>>><<<>><<>><<<>><><><<<<><<>><<><<<<>>><<>><<<>>><<<<><<>><>><<<><>>><>>><>><<>>>><<<<>>>><<<<><<<>>"
  ">><<<>><>>>><>>>><>>><<>>>><<<<>><<<>>><<<<>>><<>>><>>>><<<<>>><<<<>>>><<><<>><>>>><<><>>>><>>><<<><<>>><<<<>>>><>><"
  "<<<><<<<><>>><><<<<>>><>>><<<>>>><<>>>><<><<<>>><<<>><<>>><>><>>><<>>><<<<>>>><<<>><><<<<>>><<<<>><<>>>><<>>>><>>><<"
  "<>>>><<<>>>><<>>><<>><<<>>>><<<<>><<>>>><<<>>>><>>><<>>><<>><<<>>>><<>><><>><<<><<<>>><>>>><<><><<<<>><>>><<<<>><<<>"
  "<>>>><>>>><<><<<<>>>><>>>><>><<<<>>>><<>><<<<><<<><<<<><<>><>><>>><>><<>><<<>><<<<><<<>>>><<<<>><<><<<<>>>><<<>>><<>"
  ">>><<>><<<<><<<<>>>><<<>><<<>>>><<>>><>>>><<><<>>>><<<><>>><>>><<<>>><<<><<>>><<>>><><>>>><<<>>><<<>>>><<<<>><><<<>>"
  "><<>><<<<>>>><<>>><<<<><>><<<<>>>><<<<>>>><><>>>><>><>>><<<<><><<<>>>><<<<>><>>><<<<>>>><<>>>><<<<>><>>><<>><<>><>>>"
  "><<<><<<<>>><<<><<>>>><<>><<<>><<<<>><<<>>>><>>>><<>>><<<>><<>><<<>>>><<<>><>>>><>><<<>>>><>><<>>>><>>>><<<>>>><<<<>"
  ">>><<<><>><<<<>>>><>>>><<>><<<>><<><<<<>>>><<<<>>><<>>><<<<>><<<>><<<><<<><<<<>><<<<>>>><<>>><>>><>>>><<<>>>><<<>>><"
  "<<>>>><<>>><<<>>><>><<>>>><<<>>>><<>><<<<><<><<><<>><>>><<<<>>><<<<>>>><<<>>>><<>>><<<>>>><<<>><<>>><>>><<<<><<<<>>>"
  "><<<>>>><><<<<>>><>><<><<>><<<>>><<<<><>><<<>><>><<<<>><<<><<><<<<>>>><<>>>><<<<>>><<<<><<><<<><<<<>><<<<>>>><<<<>><"
  "<<>>><<><<>><<<>><<<>>><>>>><<>>><<<><<<<><><<<>><<<<>>><<<><>>>><>><>><<<>><>>>><>>>><>>>><<>><<<<><<<>><<<>>><<<<>"
  ">><<<>>>><<<>><><<<>><><<<<>><<<<>>><<<><<<>>><<><<>>>><>><<<><<<>><>><<<<><<>>>><>>><<<<>>>><<<>><<>>>><>>><<<<>>><"
  "<>>><<<<>><<><<><<<>>>><<<><<<<>><>>>><<<<>><>><<<>><<<>>><<<<>><>>>><<<<>><<>><<<<>>><<<<><<<<>><>><<<>>>><<<>>>><<"
  "<>>>><>>>><<<>>>><<>><<><>>><<<<>>><<<<>>>><>><<<<><<<<>>><>><<<<>>>><<<>>>><<>>><>><<><<<><<<><><<<<>>>><<<>>><<<<>"
  ">>><<>><<<>>>><<><>>><<<>><>><<<><>><<>><><<<><<><>>>><<<<>>>><>><<>><><<<>><<<>>><<>>>><<<<>><<<<><>>><>>>><<<<>>><"
  "<<<>>><<>><<>><>>>><<<>>>><<<>><<><<<<>><>>>><<<<>>>><<><<<><<<<><>><<<>><<<>>><<>><<<<>>><<<<>><<><>>>><<<>>>><<<>>"
  ">><>><<<<>>><<>><<<><<>>><<>><<<<>>><<<<>>>><<<>><<<<>><<>>>><<>>>><<<>>>><><<>>>><<<<>><<><>>>><><<<>>>><<>>><<>>><"
  "<>><<<><><>>>><><<><<<<>>>><>><<<<>>>><<<>><<<<>>>><><<<>>>><<<>><<<<>>>><>><><<<<>>>><<<>>><<<>><<<<>>>><<<>>><<<>>"
  "><<<>>>><<<><<<>>><<<>><<<>>><>>>><<<<>><><<<><<><><<<><<<>><<>><>>><>>>><<><<<><<<<><><<>>><<>>>><<>><>><<>>>><<<>>"
  "<<>>>><<<><<>><><<<<>><<<<>>>><<>>>><<<>>>><>>>><<><<<>>>><<<>>>><<<<>><<>>>><<<<><<<<>>>><<<<>>>><<<>><<>>><<<>><<<"
  "><<>>>><<<<>><<<<><<><>>><<<>><<><>>>><<<<>>>><>>>><<><<<>>><<<<><<>><<<<>>>><<<><>><<>>><>>>><<><<>>><><<<><<<<>>>>"
  "<><<<<>>><<>><<<>>>><><<<<><<<<>>>><>><<><<<>>><<<><>>><<<>>><<<<>><<<>>><<<>>><<<>>>><<<>>>><<<<><>>>><<>>><><<<<><"
  "<<><>><<<><<<><><<<>><<<>>>><<<<>><<>>>><<<<>>><<<>><<<>>><<<<>>><<>>><<>>>><<>>><<<><>>>><<>>>><>>><<>>>><>><<<>><<"
  ">>><<>><<<><<<<>><<<>>><<<<>><<<<><><<<<>>>><>>>><><>><<<>>>><<><<<><<><>><<>><<><<<>><>>><<>>>><<>>>><<>>><<<>>><<>"
  ">>><<<>><<<<>>><<>><>>><<<<>>>><<>><<<<><<<<><>><>>>><>>><<<<>>>><<<<>>><<<>><><<>>>><<>><>>>><<<><>>>><<>>>><<>>><>"
  ">><>>><<><><<<>><<>>><<>>><<>><<<><>>><<<>>>><<>>><>><<>>><><<<>>>><<>>>><>>><<<<>><><<<<><<<<>><<<>>><<>><>>>><<>>>"
  "<<<>>><<<<>><<<>>>><<<>>>><<<>>>><<><>>><<>>>><<<>>>><<<<><<<<>>>><<>>>><<<>>><<<<>>>><<<<>><<<<>><<<<><><<<>>>><<>>"
  "<<<>>>><>><>>>><<>>>><<<>><>><<<<><>><<>>><<>>><<<<>>>><<<><><><<>>>><<<<>><<<<>><><<><<<>>><<<>><<<<>><<<><<<>>>><<"
  "<><<<>>>><<><<<>><<<<>>><<<>>>><><<<>>><>>>><>>><<>>><>>>><><>>><<<<><><<>><>>>><<<><<<>>>><<<<>>>><<>>><>><<<><<><<"
  "<>>><<>><<<>>>><<<<>><<<>>><<<>>>><>><>><>><><<<>><<>>><<<>><>><>>><<<>><><<<>>>><<><><<<<><<>><<<><>><<<<>>>><<<>>>"
  "<<>>>><<<>>>><<<><<>>>><<<<>>>><<<>><<<<>>>><<><<>>><<>>>><>><>><<<<><>>><>>>><<>><>>>><<>>><<>>><<>>><>><>>>><<<>>>"
  "><<>>>><<<>><<>><<>>>><<>>><<<<><<><<<<><<>><<<<>>>><<<>><><><><>>><<<<>><><<<<>>>><>>>><<>>><<<<><<><<>><><<>>>><<<"
  ">>><<<<>>>><<<>><<<<>><<<>><>>>><<<><>>>><>>>><<<<>><>>>><<<>>>><>>>><<<<><>><>>>><<<<>>><>><<<>>>><<<><<<><<<<>>><<"
  "<>><>><>>>><<<>><<<>>>><<>>>><<>>>><<><<><<>><>><>><<<>>>><<><<<<>>><>>><>><<>><<<>>>><>>>><<<<>><><<>>><<><>>><>>><"
  "<>>><<<<><<<>>>><<<<>>>><<<><>>><<>>>><>>>><<<<>><><<><<<<>><<>><>>>><<<<><<<><>><<<<>><<>><<<<>><><<<<>><<<<><<<<><"
  "<><<<><>>>><<>>>><<<<>>><>>>><<>><><<>>>><>>>><<<>>><<<>>><<><<>>>><<>>>><<<>><<<><<<<>>><<>>>><<>>>><><<>>><<<><<<>"
  "<<><<<<>>><<<<>><<<><<<>>>><>><<>><<<><<<>>>><<><<<<>>><<<<>>><<<<>><<<>>><>>><<>>><<<<>><<<<><<<>><<<>><<<<>>><<<<>"
  ">>><<><>>><<<<>><<<>><<>>>><><>>><>><<><<<<>><<<<>>><<>>><<<<><<>>>><<<<><><<>><<<>><<>>><<<<><<<><>><>><>>><<<>><><"
  ">>>><<<><<<<>>>><>><<<<>><<<<>><><>><>>><<<<>>>><<>><<>>>><<<>>><<<><<<>>><>>>><<<>>><<<>>><<<>>>><<<>>>><>>>><<>><<"
  "<<><>><<<>><<<<>>><>>><>><>>><>>>><<<<>>>><>><<<<>>>><>><<>>><<><<><><<<><><<<>><<>>>><<><<<><<<<>>><<<><<<<>><<<<>>"
  "><<><<>>><<<><>>><<><<<<>><<<<>>><>>>><<>>><<><><<><>>>><<><<<>>>><<>>><<<<>>>><>>><<<<>>>><<>>>><<>>><<<>><>><<><>>"
  "<<<<>>>><<<<>>><<><>>><<<<>>><>>><<<>>><<<<>>>><<<>><<>><<><<<<>>><<<><><<<>>><<<>>><<>>>><<>>>><<>>><<<<><<<<><><<>"
  "><><<<<>>>><<<>>>><<>><<<<><<<<>><><<<<>>><>><<<<>>>><<<<><>><<<>><<>>><><>>><<>><<<>>><<<<>><<<>>><<>><<>>>><<<<>>>"
  "><<>><>><<<>>><<<>>>><<<>>>><<<<>>><<><<><<<<>><<<>>><><<<<>>><>><<>>>><<<<>>>><<<<><<><<>><>>><<><<<>>>><<<><><><<>"
  "><>>><<<<><<<<>>><<><<>><>><<<>>>><<<>>>><<><<<>>><>>><<>><<>><<<>>>><<>><<<<>><<<<>><<<<>><<<>>><<><<<<>>><<<<><<<<"
  ">>><>>><<<<>><<<<>>>><<>>>><<>>>><><<<><<<<>>><<<>>><<<<>>>><<>>><<>>><<<<>><<<<><<<<><<>>>><<<<>>>><><<<>>><<>><><>"
  ">><<<>>><<<><>>><>><<>>>><>>>><>>>><<<<>>><<<>>><>>><>>><<<><>><<<<>>><<>>><>>><<<>>><<>>><>><>><<>><<<>><<<>>>><<<<"
  ">><<<><<<>>>><>><>><<<>>>><<<>>><<>>><><<<<><<<<>>>><<<<>>>><><>><>><<<>><<<>><<<<>>>><<>><<<<>>><><<<>><>>>><<>>><<"
  "<<>>>><<>><<<<>><>>>><<<<><<<>><<<>>><<<<>><><<<<>>><><<>>>><<<><>>><<<<>>>><<<>><<<<>>>><<<<>>><<>><<<<>><<<<>><<>>"
  ">><><<>>>><><<<>>>><<<>>><<<<><<<>><>>><>>><<<>>>><<>>><<<<>><<>>><><<<>><><<<>>>><<<><>><<>>><<<<>><<<>>>><<>>>><<<"
  "<>>>><<>><<>>><<<>>><<<>>><>><>>>><<<<>>><<><><<<>>><<<<>>><>>><<<>><<>>>><<<<>><><<<>><<>>><<><>>><<<><<<<>><>><<<>"
  ">>><<>>>><>>><>>><<<<>>>><<<<>>><>><><<><<>><>>>><>>><<<<><<>>><<<<>>>><<<>>><<<>><<<<>>>><<>>>><><<><>>><<>>>><<>>>"
  "<<<>>>><<><<>>>><<<<>>>><>>>><>><<>>><<<>>>><><<>><<><<<<>>><<<>>><<<<>><<>><<<>>>><<<>>>><<<<>>><<<><<<>>><<><<<>>>"
  "<<<>>>><<<<>><>><<>>><<>>>><<<>>>><<<>>>><>>><<><>>>><>>><<<><<<<><><<<<>>><<<>><<<<>>><><<>><<<<>><>>><>><<<<>>>><>"
  "<<>><<>>><<<<>><<<><<<>>>><<<>>><<<<><><<<<>>>><>><<><<<<><>><<><<<<>>><<<>>>><<>>>><<<<><<<>>><<<<>>>><<<<><<>>><<<"
  "<>>><<<<><<<>><<><<<<><<>>>><><<<<>>>><<<>><>>>><>>><<<<>><><<<<>><<<><<>><<<>>>><<><<><<<><<>>>><<<<>>><<<>>><<<<><"
  "<<>>>><<<<>><>><<>>>><>>><<<<><<<<>>>><<<<>>><>>>><<<><<<<>>><>>>><>><<><>><<<<>>><<<>>><<<<><<<<>>><<<><<>>><<<<>>>"
  "<><<<<><<>><>>><<<<>>><><<<>>>><<<<><<<>><<<<>>><>>>><>><<>><<<<><><>>><<>><>>><<>>>><<<>>><<>>>><<<>><<<<>>><>><<>>"
  "<<<>>><<<>>>><<<<>>><<<>>><>>>><<>><<<<>>>><<>>>><><>>>><<<><<<>>>><<><>>>><<<<>><>>><<<><><<<<><><<<<>>>><<>><<<>>>"
  "><>>>><>>><>>><<<><<<>><>>><>>>><<<<>><<<<><<>>><<<><<<<>><<<<><<<<>><<>><<<><>><<<><><<>><<>>><><<>><>>><>>>><>>>><"
  "<<>><<>>><>>><<<<>><<<<>>>><<>>><<<>><<>>><<<><<>><<>>><<><<<<><<<>>><<<<>>><<<<><>>>><>><<<>>>><>>>><<<>>>><<<<>>>>"
  "<<>>>><>>>><><<>><<<<>>><><<><<<><<<<>>>><<<<>>><<<<>>><>>>><<><<<><<>>><<>>><<<<>>>><><<><>>>><<<<>>>><>>><<<<>>><>"
  "><<<<>>><>><<<>>>><<<><>><<<>>><<<>><<<<><<>>><<<>><><<<<>>>><<<>>>><<>>><<>>>><<<<>><<<><<>><<>><><<><><<<<>>>><>>"
};

auto can_move_left(std::vector<std::string> const& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x) -> bool
{
  // can't move left if next to the shape we would get out of bounds
  if (int64_t(shape_pos_x) - 1 <= -1) { return false; }

  for (auto y_offset = size_t{ 0 }; y_offset < shape.size(); ++y_offset) {
    auto x_offset = size_t{ 0 };
    for (; x_offset < shape[0].size(); ++x_offset) {
      if (shape[y_offset][x_offset] == '#') { break; }
    }
    auto y = shape_pos_y + y_offset;
    auto x = int64_t(shape_pos_x + x_offset);
    if (x - 1 <= -1 || tower[y][x - 1] == '#') { return false; }
  }
  return true;
}

auto can_move_right(std::vector<std::string> const& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x) -> bool
{
  // can't move right if next to the shape we would get out of bounds
  if (shape_pos_x + 1 >= tower[0].size()) { return false; }

  for (auto y_offset = size_t{ 0 }; y_offset < shape.size(); ++y_offset) {
    auto x_offset = static_cast<int>(shape[0].size());
    for (; x_offset >= 0; --x_offset) {
      if (shape[y_offset][x_offset] == '#') { break; }
    }
    auto y = shape_pos_y + y_offset;
    auto x = shape_pos_x + x_offset;
    if (x + 1 >= tower[0].size() || tower[y][x + 1] == '#') { return false; }
  }
  return true;
}

auto can_move_down(std::vector<std::string> const& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x) -> bool
{
  if (int64_t(shape_pos_y) - 1 <= -1) { return false; }

  for (auto x_offset = size_t{ 0 }; x_offset < shape[0].size(); ++x_offset) {
    auto y_offset = size_t{ 0 };
    for (; y_offset < shape.size(); ++y_offset) {
      if (shape[y_offset][x_offset] == '#') { break; }
    }
    auto y = int64_t(shape_pos_y + y_offset);
    auto x = shape_pos_x + x_offset;
    if (y - 1 <= -1 || tower[y - 1][x] == '#') { return false; }
  }
  return true;
}

auto delete_shape(std::vector<std::string>& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x)
{
  for (auto y = shape_pos_y; y < shape_pos_y + shape.size(); ++y) {
    for (auto x = shape_pos_x; x < shape_pos_x + shape[0].size(); ++x) {
      if (shape[y - shape_pos_y][x - shape_pos_x] == '#') { tower[y][x] = '.'; }
    }
  }
}

auto insert_shape(std::vector<std::string>& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x)
{
  for (auto y = shape_pos_y; y < shape_pos_y + shape.size(); ++y) {
    for (auto x = shape_pos_x; x < shape_pos_x + shape[0].size(); ++x) {
      if (shape[y - shape_pos_y][x - shape_pos_x] == '#') { tower[y][x] = '#'; }
    }
  }
}

auto move_left(std::vector<std::string>& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x)
{
  delete_shape(tower, shape, shape_pos_y, shape_pos_x);
  insert_shape(tower, shape, shape_pos_y, shape_pos_x - 1);
}

auto move_right(std::vector<std::string>& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x)
{
  delete_shape(tower, shape, shape_pos_y, shape_pos_x);
  insert_shape(tower, shape, shape_pos_y, shape_pos_x + 1);
}

auto move_down(std::vector<std::string>& tower,
  std::vector<std::string> const& shape,
  size_t shape_pos_y,
  size_t shape_pos_x)
{
  delete_shape(tower, shape, shape_pos_y, shape_pos_x);
  insert_shape(tower, shape, shape_pos_y - 1, shape_pos_x);
}

auto print_dir(size_t jet_index)
{
  std::cout << jets << std::endl;
  for (auto i = size_t{ 0 }; i < jets.size(); ++i) {
    if (i == jet_index) {
      std::cout << '|';
      continue;
    }
    std::cout << ' ';
  }
  std::cout << std::endl;
}

auto print_tower(std::vector<std::string> const& tower, unsigned sleep_for_millis)
{
  for (auto it = tower.crbegin(); it != tower.crend(); ++it) { std::cout << '|' << *it << '|' << std::endl; }
  std::cout << std::string("+-------+") << std::endl << std::endl;
  std::this_thread::sleep_for(std::chrono::milliseconds(sleep_for_millis));
}

constexpr auto spawn_in_distance_left = 2;
constexpr auto spawn_in_distance_bottom = 3;

auto concat_last_5_patterns(std::vector<std::string> const& patterns) -> std::optional<std::string>
{
  if (patterns.size() < 5) { return std::nullopt; }
  auto concatenated = std::string{};
  auto count = 0;
  for (auto pattern = patterns.crbegin(); pattern != patterns.crend() && count < 5; ++pattern, ++count) {
    concatenated += *pattern;
  }
  return std::make_optional(concatenated);
}

struct CycleWatch
{
  std::vector<std::string> _patterns;
  std::vector<size_t> _heights;
  // 5 patterns -> _heights index
  std::unordered_map<std::string, size_t> _pattern_cycle_tracker;

  // will return true if with the current entry we would detect a pattern
  // so the length of the cycle is: _patterns.size() - _pattern_cycle_tracker[last_5_pattern]
  auto add_pattern(std::string_view pattern, size_t current_height) -> std::pair<bool, std::string>
  {
    _patterns.emplace_back(pattern);
    _heights.emplace_back(current_height);
    auto cycle_pattern = concat_last_5_patterns(_patterns);
    if (cycle_pattern.has_value()) {
      if (_pattern_cycle_tracker.contains(cycle_pattern.value())) {
        return std::make_pair(true, cycle_pattern.value());
      }
      _pattern_cycle_tracker.insert({ cycle_pattern.value(), _heights.size() - 1 });
    }
    return std::make_pair(false, cycle_pattern.has_value() ? cycle_pattern.value() : std::string{});
  }
};

auto simulate(size_t max_shape_count, bool render) -> size_t
{
  auto tower = std::vector<std::string>{ ".......", ".......", "......." };
  auto shape_count = size_t{ 0 };

  auto jet_index = size_t{ 0 };
  auto shape_index = size_t{ 0 };
  auto shape_pos_y = tower.size();
  auto shape_pos_x = size_t{ spawn_in_distance_left };
  auto highest_y = size_t(0);

  auto check_for_cycles_after = size_t{ 200 };
  auto cycle_detected = false;
  auto cycle_pattern = std::string{};
  CycleWatch cycle_watch;

  while (shape_count < max_shape_count && !cycle_detected) {
    // add the new shape
    while (tower.size() < (shape_pos_y + shapes[shape_index].size())) { tower.emplace_back("......."); }
    if (render) {
      insert_shape(tower, shapes[shape_index], shape_pos_y, shape_pos_x);
      std::cout << "Initial" << std::endl;
      print_tower(tower, 1000);
    }

    auto shape_move_pattern = std::to_string(shape_index);
    // move shape
    while (true) {
      auto should_move_left = true;
      if (jets[jet_index] == '>') { should_move_left = false; }

      if (should_move_left) {
        if (can_move_left(tower, shapes[shape_index], shape_pos_y, shape_pos_x)) {
          if (render) { move_left(tower, shapes[shape_index], shape_pos_y, shape_pos_x); }
          --shape_pos_x;
          shape_move_pattern += ">";
        } else {
          shape_move_pattern += "X";
        }
      } else {
        if (can_move_right(tower, shapes[shape_index], shape_pos_y, shape_pos_x)) {
          if (render) { move_right(tower, shapes[shape_index], shape_pos_y, shape_pos_x); }
          ++shape_pos_x;
          shape_move_pattern += "<";
        } else {
          shape_move_pattern += "X";
        }
      }
      jet_index = (jet_index + 1) % jets.size();

      if (!can_move_down(tower, shapes[shape_index], shape_pos_y, shape_pos_x)) {
        if (render) {
          std::cout << "Moving" << std::endl;
          print_tower(tower, 300);
        } else {// even if we aren't rendering we must place the shape into it's final position
          insert_shape(tower, shapes[shape_index], shape_pos_y, shape_pos_x);
        }
        break;
      }

      if (render) {
        move_down(tower, shapes[shape_index], shape_pos_y, shape_pos_x);
        std::cout << "Moving" << std::endl;
        print_tower(tower, 300);
      }
      --shape_pos_y;
      shape_move_pattern += "|";
    }

    // restore shape spawn in position
    shape_pos_x = size_t{ spawn_in_distance_left };
    highest_y = std::max(highest_y, shape_pos_y + shapes[shape_index].size());
    shape_pos_y = highest_y + spawn_in_distance_bottom;

    shape_index = (shape_index + 1) % shapes.size();
    ++shape_count;
    if (shape_count > check_for_cycles_after) {
      // check for cycles
      auto [cycle_found, cycle_marker] = cycle_watch.add_pattern(shape_move_pattern, highest_y);
      // bah, can't auto unpack without making new variables....
      cycle_detected = cycle_found;
      cycle_pattern = cycle_marker;
    }
  }
  auto remaining_elements_to_simulate = max_shape_count - shape_count;
  auto additional_height = size_t{ 0 };
  // got a cycle and have elements to calculate
  if (remaining_elements_to_simulate > 0) {
    auto const cycle_start = cycle_watch._pattern_cycle_tracker.find(cycle_pattern);
    // the -1 is necessary because _heights gets extended with the final height where the new cycle would start
    auto const elements_in_cycle = cycle_watch._heights.size() - 1 - cycle_start->second;
    auto const full_cycle_count = remaining_elements_to_simulate / elements_in_cycle;
    auto const base_height = cycle_watch._heights[cycle_start->second];

    additional_height += full_cycle_count * (highest_y - base_height);
    remaining_elements_to_simulate = remaining_elements_to_simulate % elements_in_cycle;
    additional_height += cycle_watch._heights[cycle_start->second + remaining_elements_to_simulate] - base_height;
  }

  return highest_y + additional_height;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  // constexpr auto max_shape_count = size_t{2022};
  constexpr auto max_shape_count = size_t{ 1000000000000 };
  auto total_height = simulate(max_shape_count, false);
  std::cout << "Height of the tower is: " << total_height << std::endl;

  return 0;
}


// Without pattern detection 2022 with long jet input
// mrai@potatoland:~/tinker/cp/advent_of_code/2022$ time ./day_17
// Height of the tower is: 3219

// real    0m0.001s
// user    0m0.001s
// sys     0m0.000s