#include <bits/stdc++.h>

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

auto rollover(int value, int max) -> int
{
  if (value >= max) {
    return 0;
  } else if (value < 0) {
    return max - 1;
  }
  return value;
}

auto next_valid_i(std::vector<std::string> const& matrix, int i, int j, int offset_i) -> int
{
  auto next_i = i + offset_i;
  next_i = rollover(next_i, matrix.size());
  for (; next_i < matrix.size() && matrix[next_i][j] == ' '; next_i += offset_i) {}
  next_i = rollover(next_i, matrix.size());
  for (; next_i < matrix.size() && matrix[next_i][j] == ' '; next_i += offset_i) {}
  if (matrix[next_i][j] == '#') { return i; }
  return next_i;
}

// yes, yes code duplication, sue me
auto next_valid_j(std::vector<std::string> const& matrix, int i, int j, int offset_j) -> int
{
  auto next_j = j + offset_j;
  next_j = rollover(next_j, matrix[0].size());
  for (; next_j < matrix[0].size() && matrix[i][next_j] == ' '; next_j += offset_j) {}
  next_j = rollover(next_j, matrix[0].size());
  for (; next_j < matrix[0].size() && matrix[i][next_j] == ' '; next_j += offset_j) {}
  if (matrix[i][next_j] == '#') { return j; }
  return next_j;
}

auto simulate(std::vector<std::string>& matrix,
  int i,
  int j,
  std::pair<int, int> offset,
  char direction_marker,
  uint32_t count_until) -> std::pair<int, int>
{
  auto next_i = i;
  auto next_j = j;
  matrix[i][j] = direction_marker;
  for (auto count = uint32_t{ 0 }; count < count_until; ++count) {
    next_i = next_valid_i(matrix, i, j, offset.first);
    next_j = next_valid_j(matrix, i, j, offset.second);
    if (next_i == i && next_j == j) { break; }
    matrix[next_i][next_j] = direction_marker;
    i = next_i;
    j = next_j;
  }
  return std::make_pair(next_i, next_j);
}

// this function will assign a unique id to each face, between [1-6], the non face sectors
// will be marked by 0
auto compress_map_into_cube_face_diagram(std::vector<std::string> const& matrix, size_t face_width)
  -> std::vector<std::vector<size_t>>
{
  assert(matrix.size() % face_width == 0);
  assert(matrix[0].size() % face_width == 0);

  auto row_count = matrix.size() / face_width;
  auto column_count = matrix[0].size() / face_width;
  auto compressed_map = std::vector<std::vector<size_t>>(row_count, std::vector<size_t>(column_count, 0));
  auto face_id = size_t{ 1 };
  for (auto i = size_t{ 0 }; i < row_count; ++i) {
    for (auto j = size_t{ 0 }; j < column_count; ++j) {
      if (matrix[i * face_width][j * face_width] != ' ') {
        compressed_map[i][j] = face_id;
        ++face_id;
      }
    }
  }
  return compressed_map;
}

auto main() -> int
{
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);

  auto matrix = std::vector<std::string>{};
  auto longest_line = size_t{ 0 };
  auto line = std::string{};
  while (std::getline(std::cin, line)) {
    if (line.empty()) { break; }
    longest_line = std::max(longest_line, line.size());
    matrix.push_back(line);
  }

  // normalize line widths
  for (auto& line : matrix) { line.resize(longest_line, ' '); }
  // figure out the width of a face, by finding the shortest non space filled line
  // there has to be one such either by scanning horizontally or vertically otherwise
  // the cube cannot be assembled
  auto face_width = std::numeric_limits<size_t>::max();
  for (auto const& line : matrix) {
    auto count = size_t{ 0 };
    for (auto c : line) {
      if (c != ' ') {
        ++count;
      } else if (count != size_t{ 0 }) {
        face_width = std::min(face_width, count);
        count = size_t{ 0 };
      }
    }
    if (count != size_t{ 0 }) { face_width = std::min(face_width, count); }
  }

  for (auto j = size_t{ 0 }; j < matrix[0].size(); ++j) {
    auto count = size_t{ 0 };
    for (auto i = size_t{ 0 }; i < matrix.size(); ++i) {
      if (matrix[i][j] != ' ') {
        ++count;
      } else if (count != size_t{ 0 }) {
        face_width = std::min(face_width, count);
        count = size_t{ 0 };
      }
    }
    if (count != size_t{ 0 }) { face_width = std::min(face_width, count); }
  }

  std::cout << "Face width is: " << face_width << std::endl;
  auto cube_face_diagram = compress_map_into_cube_face_diagram(matrix, face_width);
  for (auto const& line : cube_face_diagram) {
    dump_to(std::cout, line.cbegin(), line.cend());
    std::cout << std::endl;
  }

  auto directions = std::string{};
  std::getline(std::cin, directions);

  auto start_i = int{ 0 };
  auto start_j = int{ 0 };
  auto first_empty_at = std::find(matrix[0].cbegin(), matrix[0].cend(), '.');
  start_j = std::distance(matrix[0].cbegin(), first_empty_at);

  auto dir_data = std::istringstream{ directions };
  auto count = uint32_t{ 0 };
  auto offsets = std::array<std::pair<int, int>, 4>{
    std::make_pair(0, 1), std::make_pair(1, 0), std::make_pair(0, -1), std::make_pair(-1, 0)
  };
  auto direction_marker = std::array<char, 4>{ '>', 'v', '<', '^' };
  auto offset_index = uint8_t{ 0 };
  auto turn = ' ';
  auto read_number = true;
  while (!dir_data.eof()) {
    if (read_number) {
      dir_data >> count;
      auto [next_i, next_j] =
        simulate(matrix, start_i, start_j, offsets[offset_index], direction_marker[offset_index], count);

      start_i = next_i;
      start_j = next_j;
    } else {
      dir_data >> turn;
      if (turn == 'R') {
        offset_index = (offset_index + 1) % 4;
      } else if (turn == 'L') {
        offset_index = (offset_index + 3) % 4;
      }
    }
    read_number = !read_number;
  }
  // plus 1, because the exercise is 1 indexed not 0
  std::cout << "Final password: " << ((start_i + 1) * 1000) + ((start_j + 1) * 4) + offset_index << std::endl;

  return 0;
}
