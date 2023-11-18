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

// as given by the exercise description, only shifted to the right by
// one, so it neatly aligns with the rest of the directions
constexpr auto offsets = std::array<std::pair<int, int>, 4>{ std::make_pair(-1, 0),
  std::make_pair(0, 1),
  std::make_pair(1, 0),
  std::make_pair(0, -1) };
constexpr auto direction_markers = std::array<char, 4>{ '^', '>', 'v', '<' };

auto parse_matrix() -> std::vector<std::string>
{
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
  return matrix;
}

auto get_face_width(std::vector<std::string> const& matrix) -> size_t
{
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
  return face_width;
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

class Plane
{
public:
  Plane(std::vector<std::string> matrix, int i, int j, int offset_index)
    : _matrix{ std::move(matrix) }, _i{ i }, _j{ j }, _offset_index{ offset_index }
  {}

  auto simulate_move(uint32_t count_until) -> void
  {
    auto offset = offsets[_offset_index];
    auto direction_marker = direction_markers[_offset_index];

    auto next_i = _i;
    auto next_j = _j;
    _matrix[_i][_j] = direction_marker;
    for (auto count = uint32_t{ 0 }; count < count_until; ++count) {
      next_i = next_valid_i(_matrix, _i, _j, offset.first);
      next_j = next_valid_j(_matrix, _i, _j, offset.second);
      if (next_i == _i && next_j == _j) { break; }
      _matrix[next_i][next_j] = direction_marker;
      _i = next_i;
      _j = next_j;
    }
  }

  auto get_i() const -> int { return _i; }

  auto get_j() const -> int { return _j; }

  auto get_offset_index() const -> int { return _offset_index; }

  auto change_direction(char rotation)
  {
    if (rotation == 'R') {
      _offset_index = (_offset_index + 1) % 4;
    } else if (rotation == 'L') {
      _offset_index = (_offset_index + 3) % 4;
    }
  }

private:
  std::vector<std::string> _matrix;
  int _i;
  int _j;
  int _offset_index;
};

// Direction in the 2D plane
enum class Direction {
  NONE,
  UP,
  RIGHT,
  DOWN,
  LEFT,
};

auto operator<<(std::ostream& os, Direction dir) -> std::ostream&
{
  switch (dir) {
  case Direction::UP: {
    os << "U";
    break;
  }
  case Direction::RIGHT: {
    os << "R";
    break;
  }
  case Direction::DOWN: {
    os << "D";
    break;
  }
  case Direction::LEFT: {
    os << "L";
    break;
  }
  default: {
    os << "X";
    break;
  }
  }
  return os;
}

auto int_to_direction(int direction) -> Direction
{
  switch (direction) {
  case 0: {
    return Direction::UP;
  }
  case 1: {
    return Direction::RIGHT;
  }
  case 2: {
    return Direction::DOWN;
  }
  case 3: {
    return Direction::LEFT;
  }
  default: {
    return Direction::NONE;
  }
  }
}

// in the numbering scheme we keep the clockwise ordering
// so that it is easy to serialize directions into an array as well
auto direction_to_int(Direction dir) -> int
{
  switch (dir) {
  case Direction::UP: {
    return 0;
  }
  case Direction::RIGHT: {
    return 1;
  }
  case Direction::DOWN: {
    return 2;
  }
  case Direction::LEFT: {
    return 3;
  }
  default: {
    return 4;
  }
  }
}

// a cube side in the 2D plane
struct Face
{
  int top_left_i;
  int top_left_j;
  // 4 elements total for each cardinal direction
  // always in this order: top, right, bottom, left
  std::array<char, 4> neighbors;

  // we ask the face, where is the argument relative to you, on which side
  auto face_connected_in_direction(char from_face_id) const -> Direction
  {
    auto connected_in_direction = 0;
    for (; connected_in_direction < 4; ++connected_in_direction) {
      if (neighbors[connected_in_direction] == from_face_id) { break; }
    }
    return int_to_direction(connected_in_direction);
  }
};

auto on_face_sides(int face_width, int i, int j) -> std::vector<Direction>
{
  auto on_sides = std::vector<Direction>{};
  if (i == 0 && 0 <= j && j < face_width) { on_sides.push_back(Direction::UP); }
  if (j == face_width - 1 && 0 <= i && i < face_width) { on_sides.push_back(Direction::RIGHT); }
  if (i == face_width - 1 && 0 <= j && j < face_width) { on_sides.push_back(Direction::DOWN); }
  if (j == 0 && 0 <= i && i < face_width) { on_sides.push_back(Direction::LEFT); }
  return on_sides;
}

auto calculate_offset_index(Direction dir) -> int
{
  switch (dir) {
  case Direction::UP: {
    return direction_to_int(Direction::DOWN);
  }
  case Direction::LEFT: {
    return direction_to_int(Direction::RIGHT);
  }
  case Direction::RIGHT: {
    return direction_to_int(Direction::LEFT);
  }
  case Direction::DOWN: {
    return direction_to_int(Direction::UP);
  }
  default: {
    return direction_to_int(Direction::NONE);
  }
  }
}

// ahh a bit confusing names, we are essentially identifying the edge at which two faces are touching
// in the fully assembled cube. this edge is identified here by the direction we are stepping off from the
// 'from' face, and by the direction we would have to step off from the target face to arrive back at 'from'
// PS.: My apologies this code is disturbing even to me.
// there are no checks whether or not from_i, from_j are valid for a given face_width
auto simulate_stepping_between_faces(Direction const from_face_dir,
  int from_i,
  int from_j,
  int face_width,
  Direction const direction_from_to_face_dir) -> std::tuple<int, int, int>// i, j on new face
{
  auto candidate_i = from_i;
  auto candidate_j = from_j;
  // the offset index is fully specified by direction_from_to_face_dir, as it has to be it's opposite
  auto candidate_offset = (direction_to_int(direction_from_to_face_dir) + 2) % 4;

  // the directions are the same, ex.: UP <-> UP
  if (from_face_dir == direction_from_to_face_dir) {
    if (from_face_dir == Direction::UP || from_face_dir == Direction::DOWN) {
      candidate_j = face_width - 1 - from_j;
    } else if (from_face_dir == Direction::LEFT || from_face_dir == Direction::RIGHT) {
      candidate_i = face_width - 1 - from_i;
    }
  } else {
    // if the directions are perpendicular, ex.: TOP <-> RIGHT
    if ((direction_to_int(from_face_dir) + 2) % 4 != direction_to_int(direction_from_to_face_dir)) {
      if (from_face_dir == Direction::UP && direction_from_to_face_dir == Direction::RIGHT) {
        candidate_i = face_width - 1 - from_j;
      } else if (from_face_dir == Direction::UP && direction_from_to_face_dir == Direction::LEFT) {
        candidate_i = from_j;
      } else if (from_face_dir == Direction::LEFT && direction_from_to_face_dir == Direction::UP) {
        candidate_j = from_i;
      } else if (from_face_dir == Direction::LEFT && direction_from_to_face_dir == Direction::DOWN) {
        candidate_j = face_width - 1 - from_i;
      } else if (from_face_dir == Direction::RIGHT && direction_from_to_face_dir == Direction::UP) {
        candidate_j = face_width - 1 - from_i;
      } else if (from_face_dir == Direction::RIGHT && direction_from_to_face_dir == Direction::DOWN) {
        candidate_j = from_i;
      } else if (from_face_dir == Direction::DOWN && direction_from_to_face_dir == Direction::LEFT) {
        candidate_i = face_width - 1 - from_j;
      } else if (from_face_dir == Direction::DOWN && direction_from_to_face_dir == Direction::RIGHT) {
        candidate_i = from_j;
      }
    }
    // else directions are opposing in which case i, j doesn't change
  }
  // then solely on direction_from_to_face_dir we can set one of the indexes
  if (direction_from_to_face_dir == Direction::UP) {
    candidate_i = 0;
  } else if (direction_from_to_face_dir == Direction::RIGHT) {
    candidate_j = face_width - 1;
  } else if (direction_from_to_face_dir == Direction::DOWN) {
    candidate_i = face_width - 1;
  } else if (direction_from_to_face_dir == Direction::LEFT) {
    candidate_j = 0;
  }

  return std::make_tuple(candidate_i, candidate_j, candidate_offset);
}

class Cube
{
public:
  // we are assuming we always start from the bottom face from position 0, 0
  // this is okay, because that is how we always parse the cube shape
  Cube(std::vector<std::string> matrix, size_t face_width, int offset_index)
    : _matrix{ matrix }, _face_width{ face_width }, _i{ 0 }, _j{ 0 }, _offset_index{ offset_index }, _on_face{ 'B' }
  {}

  auto set_face(char face_id, Face face) { _faces[face_id] = std::move(face); }

  auto simulate_move(uint32_t count_until) -> void
  {
    // For nostalgic purposes, here is where I made a mistake.
    // initially wrote prev_i = _i; which is super wrong!
    // _i, _j is in facet coordinates not global plane coordinates, so if we manage to transition between faces such
    // that our previous and current face local coordinates match we would stop thinking we hit a rock. A very nice bug
    // too, because it takes a specific fold and transition to even notice it. Took me just about 10 hours to fix, or
    // more but I am not gonna admit to that. :D
    auto prev_plane_i = get_plane_i();
    auto prev_plane_j = get_plane_j();
    auto direction_marker = direction_markers[_offset_index];

    _matrix[get_plane_i()][get_plane_j()] = direction_marker;
    for (auto count = uint32_t{ 0 }; count < count_until; ++count) {
      next_valid();
      auto current_plane_i = get_plane_i();
      auto current_plane_j = get_plane_j();
      if (current_plane_i == prev_plane_i && current_plane_j == prev_plane_j) { break; }
      prev_plane_i = current_plane_i;
      prev_plane_j = current_plane_j;
      direction_marker = direction_markers[_offset_index];
      _matrix[get_plane_i()][get_plane_j()] = direction_marker;
    }
    _matrix[get_plane_i()][get_plane_j()] = 'X';
  }

  auto get_face_i() const -> int { return _i; }
  auto get_face_j() const -> int { return _j; }

  auto get_plane_i() const -> int
  {
    auto current_face = _faces.at(_on_face);
    return get_plane_i(current_face, _i);
  }

  auto get_plane_j() const -> int
  {
    auto current_face = _faces.at(_on_face);
    return get_plane_j(current_face, _j);
  }

  auto get_offset_index() const -> int { return _offset_index; }

  auto change_direction(char rotation)
  {
    if (rotation == 'R') {
      _offset_index = (_offset_index + 1) % 4;
    } else if (rotation == 'L') {
      _offset_index = (_offset_index + 3) % 4;
    }
  }

  auto get_matrix() const -> std::vector<std::string> { return _matrix; }

  auto get_faces() const -> std::map<char, Face> const& { return _faces; }

private:
  auto get_plane_i(Face const& face, int i) const -> int { return face.top_left_i + i; }

  auto get_plane_j(Face const& face, int j) const -> int { return face.top_left_j + j; }

  auto next_valid() -> void
  {
    auto offset = offsets[_offset_index];
    auto const& current_face = _faces.at(_on_face);

    auto on_edges = on_face_sides(_face_width, _i, _j);

    auto offset_direction = int_to_direction(_offset_index);
    auto moving_off_the_edge = std::any_of(
      on_edges.cbegin(), on_edges.cend(), [&offset_direction](auto const& dir) { return dir == offset_direction; });
    // if we aren't moving off the edge
    if (!moving_off_the_edge) {
      auto candidate_i = _i + offset.first;
      auto candidate_j = _j + offset.second;
      // these have to be valid coordinates at this point
      if (_matrix[get_plane_i(current_face, candidate_i)][get_plane_j(current_face, candidate_j)] != '#') {
        _i = candidate_i;
        _j = candidate_j;
      }
      return;
    }
    // crossing an edge
    auto target_face_id = current_face.neighbors[_offset_index];

    auto target_face = _faces.at(target_face_id);
    auto direction_from_target = target_face.face_connected_in_direction(_on_face);
    assert(direction_from_target != Direction::NONE);

    // remember offset_direction at this point is the same as the side of the face we are stepping off from
    auto [candidate_i, candidate_j, candidate_offset_index] =
      simulate_stepping_between_faces(offset_direction, _i, _j, _face_width, direction_from_target);
    if (_matrix[get_plane_i(target_face, candidate_i)][get_plane_j(target_face, candidate_j)] != '#') {
      _i = candidate_i;
      _j = candidate_j;
      _offset_index = candidate_offset_index;
      _on_face = target_face_id;
    }
  }

private:
  std::vector<std::string> _matrix;
  size_t _face_width;
  int _i;
  int _j;
  int _offset_index;
  char _on_face;
  // for each face of the cube we only need to know
  // where it's top left corner is on the 2D plane
  std::map<char, Face> _faces;
};

// these represent neighbors for each side, always in a clockwise order
// for a valid cube folding in
const auto side_neighbor_mapping =
  std::map<char, std::array<char, 4>>{ std::make_pair('B', std::array<char, 4>{ 'N', 'E', 'S', 'W' }),
    std::make_pair('N', std::array<char, 4>{ 'T', 'E', 'B', 'W' }),
    std::make_pair('T', std::array<char, 4>{ 'S', 'E', 'N', 'W' }),
    std::make_pair('W', std::array<char, 4>{ 'N', 'B', 'S', 'T' }),
    std::make_pair('E', std::array<char, 4>{ 'T', 'S', 'B', 'N' }),
    std::make_pair('S', std::array<char, 4>{ 'T', 'W', 'B', 'E' }) };

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

// this relative simple function uses the assumption that after
// we have aligned the neighbor arrays, we can fill in the missing spaces
// because these neighbor arrays always represent the neighbors of a face
// in a clockwise order
// the baseline neighbors always starts from the top, but actuals could
// be rotated, that is why we need alignment
auto align_and_fill_neighbors(char face_id, std::array<char, 4>& actual)
{
  auto const& baseline_neighbors = side_neighbor_mapping.at(face_id);
  auto source_i = int{ 0 };
  for (; source_i < actual.size(); ++source_i) {
    if (actual[source_i] != ' ') { break; }
  }
  auto target_i = int{ 0 };
  if (source_i != actual.size()) {
    for (; target_i < baseline_neighbors.size(); ++target_i) {
      if (baseline_neighbors[target_i] == actual[source_i]) { break; }
    }
  }
  // after alignment, fill actual with the rest of the sides
  for (auto count = size_t{ 0 }; count < actual.size(); ++count) {
    actual[source_i] = baseline_neighbors[target_i];
    source_i = (source_i + 1) % actual.size();
    target_i = (target_i + 1) % baseline_neighbors.size();
  }
}

// this function will assign the cube side identifiers to the 2D plane sectors
auto diagram_into_cube(std::vector<std::string> const& matrix,
  std::vector<std::vector<size_t>> const& diagram,
  size_t face_width) -> std::pair<Cube, std::vector<std::vector<char>>>
{
  constexpr auto directions = std::array<std::pair<int, int>, 4>{
    std::make_pair(-1, 0), std::make_pair(0, 1), std::make_pair(1, 0), std::make_pair(0, -1)
  };

  auto i = int{ 0 };
  auto j = int{ 0 };
  // find first cube face
  for (; j < diagram[0].size(); ++j) {
    if (diagram[i][j] != 0) { break; }
  }

  auto cube_projection = std::vector<std::vector<char>>(diagram.size(), std::vector<char>(diagram[0].size(), ' '));

  auto cube = Cube{ matrix, face_width, 1 };

  std::queue<std::pair<char, Face>> to_check;
  to_check.emplace(std::make_pair('B', Face(i * face_width, j * face_width, side_neighbor_mapping.at('B'))));
  while (!to_check.empty()) {
    auto checking = to_check.front();
    to_check.pop();

    i = checking.second.top_left_i / face_width;
    j = checking.second.top_left_j / face_width;
    if (i < 0 || i >= cube_projection.size() || j < 0 || j >= cube_projection[0].size() || diagram[i][j] == 0
        || cube_projection[i][j] != ' ') {
      continue;
    }

    align_and_fill_neighbors(checking.first, checking.second.neighbors);
    cube.set_face(checking.first, checking.second);
    cube_projection[i][j] = checking.first;

    for (auto direction_index = size_t{ 0 }; direction_index < directions.size(); ++direction_index) {
      auto next_i = i + directions[direction_index].first;
      auto next_j = j + directions[direction_index].second;

      auto neighbors = std::array<char, 4>{ ' ', ' ', ' ', ' ' };
      neighbors[(direction_index + 2) % 4] = checking.first;

      to_check.push(std::make_pair(
        checking.second.neighbors[direction_index], Face(next_i * face_width, next_j * face_width, neighbors)));
    }
  }

  return std::make_pair(cube, cube_projection);
}

auto calculate_password(int i, int j, int offset_index) -> int
{
  // plus 1, because the exercise is 1 indexed not 0
  // we have to move the direction offset back by one, because we shifted it in the beginning
  return ((i + 1) * 1000) + ((j + 1) * 4) + ((offset_index + 3) % 4);
}

auto test_simulate_stepping_between_faces(Direction const from_face_dir,
  int from_i,
  int from_j,
  int face_width,
  Direction const direction_from_to_face_dir,
  std::tuple<int, int, int> expected) -> void
{
  auto got = simulate_stepping_between_faces(from_face_dir, from_i, from_j, face_width, direction_from_to_face_dir);
  std::cout << from_i << " " << from_j << ": " << from_face_dir << " -> " << direction_from_to_face_dir << ": ";
  if (got == expected) {
    std::cout << "PASS";
  } else {
    std::cout << "FAIL - got: (" << std::get<0>(got) << "," << std::get<1>(got) << "," << std::get<2>(got)
              << ") expected: (" << std::get<0>(expected) << "," << std::get<1>(expected) << ","
              << std::get<2>(expected) << ")";
  }
  std::cout << std::endl;
}

auto main() -> int
{
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);

  auto matrix = parse_matrix();
  auto face_width = get_face_width(matrix);

  std::cout << "TEST simulate_stepping_between_faces" << std::endl;
  // Resolved to this, I was so stuck on the coordinates bug!!!
  // same
  test_simulate_stepping_between_faces(Direction::UP, 0, 0, 5, Direction::UP, std::make_tuple(0, 4, 2));
  test_simulate_stepping_between_faces(Direction::UP, 0, 4, 5, Direction::UP, std::make_tuple(0, 0, 2));
  test_simulate_stepping_between_faces(Direction::UP, 0, 1, 5, Direction::UP, std::make_tuple(0, 3, 2));

  test_simulate_stepping_between_faces(Direction::RIGHT, 0, 4, 5, Direction::RIGHT, std::make_tuple(4, 4, 3));
  test_simulate_stepping_between_faces(Direction::RIGHT, 4, 4, 5, Direction::RIGHT, std::make_tuple(0, 4, 3));
  test_simulate_stepping_between_faces(Direction::RIGHT, 1, 4, 5, Direction::RIGHT, std::make_tuple(3, 4, 3));

  test_simulate_stepping_between_faces(Direction::DOWN, 4, 0, 5, Direction::DOWN, std::make_tuple(4, 4, 0));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 4, 5, Direction::DOWN, std::make_tuple(4, 0, 0));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 1, 5, Direction::DOWN, std::make_tuple(4, 3, 0));

  test_simulate_stepping_between_faces(Direction::LEFT, 0, 0, 5, Direction::LEFT, std::make_tuple(4, 0, 1));
  test_simulate_stepping_between_faces(Direction::LEFT, 4, 0, 5, Direction::LEFT, std::make_tuple(0, 0, 1));
  test_simulate_stepping_between_faces(Direction::LEFT, 1, 0, 5, Direction::LEFT, std::make_tuple(3, 0, 1));
  // opposing
  test_simulate_stepping_between_faces(Direction::UP, 0, 0, 5, Direction::DOWN, std::make_tuple(4, 0, 0));
  test_simulate_stepping_between_faces(Direction::UP, 0, 4, 5, Direction::DOWN, std::make_tuple(4, 4, 0));
  test_simulate_stepping_between_faces(Direction::UP, 0, 1, 5, Direction::DOWN, std::make_tuple(4, 1, 0));

  test_simulate_stepping_between_faces(Direction::RIGHT, 0, 4, 5, Direction::LEFT, std::make_tuple(0, 0, 1));
  test_simulate_stepping_between_faces(Direction::RIGHT, 4, 4, 5, Direction::LEFT, std::make_tuple(4, 0, 1));
  test_simulate_stepping_between_faces(Direction::RIGHT, 1, 4, 5, Direction::LEFT, std::make_tuple(1, 0, 1));

  test_simulate_stepping_between_faces(Direction::DOWN, 4, 0, 5, Direction::UP, std::make_tuple(0, 0, 2));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 4, 5, Direction::UP, std::make_tuple(0, 4, 2));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 1, 5, Direction::UP, std::make_tuple(0, 1, 2));

  test_simulate_stepping_between_faces(Direction::LEFT, 0, 0, 5, Direction::RIGHT, std::make_tuple(0, 4, 3));
  test_simulate_stepping_between_faces(Direction::LEFT, 4, 0, 5, Direction::RIGHT, std::make_tuple(4, 4, 3));
  test_simulate_stepping_between_faces(Direction::LEFT, 1, 0, 5, Direction::RIGHT, std::make_tuple(1, 4, 3));
  // perpendicular
  test_simulate_stepping_between_faces(Direction::UP, 0, 0, 5, Direction::RIGHT, std::make_tuple(4, 4, 3));
  test_simulate_stepping_between_faces(Direction::UP, 0, 4, 5, Direction::RIGHT, std::make_tuple(0, 4, 3));
  test_simulate_stepping_between_faces(Direction::UP, 0, 1, 5, Direction::RIGHT, std::make_tuple(3, 4, 3));

  test_simulate_stepping_between_faces(Direction::UP, 0, 0, 5, Direction::LEFT, std::make_tuple(0, 0, 1));
  test_simulate_stepping_between_faces(Direction::UP, 0, 4, 5, Direction::LEFT, std::make_tuple(4, 0, 1));
  test_simulate_stepping_between_faces(Direction::UP, 0, 1, 5, Direction::LEFT, std::make_tuple(1, 0, 1));

  test_simulate_stepping_between_faces(Direction::LEFT, 0, 0, 5, Direction::UP, std::make_tuple(0, 0, 2));
  test_simulate_stepping_between_faces(Direction::LEFT, 4, 0, 5, Direction::UP, std::make_tuple(0, 4, 2));
  test_simulate_stepping_between_faces(Direction::LEFT, 1, 0, 5, Direction::UP, std::make_tuple(0, 1, 2));

  test_simulate_stepping_between_faces(Direction::LEFT, 0, 0, 5, Direction::DOWN, std::make_tuple(4, 4, 0));
  test_simulate_stepping_between_faces(Direction::LEFT, 4, 0, 5, Direction::DOWN, std::make_tuple(4, 0, 0));
  test_simulate_stepping_between_faces(Direction::LEFT, 1, 0, 5, Direction::DOWN, std::make_tuple(4, 3, 0));

  test_simulate_stepping_between_faces(Direction::RIGHT, 0, 4, 5, Direction::UP, std::make_tuple(0, 4, 2));
  test_simulate_stepping_between_faces(Direction::RIGHT, 4, 4, 5, Direction::UP, std::make_tuple(0, 0, 2));
  test_simulate_stepping_between_faces(Direction::RIGHT, 1, 4, 5, Direction::UP, std::make_tuple(0, 3, 2));

  test_simulate_stepping_between_faces(Direction::RIGHT, 0, 4, 5, Direction::DOWN, std::make_tuple(4, 0, 0));
  test_simulate_stepping_between_faces(Direction::RIGHT, 4, 4, 5, Direction::DOWN, std::make_tuple(4, 4, 0));
  test_simulate_stepping_between_faces(Direction::RIGHT, 1, 4, 5, Direction::DOWN, std::make_tuple(4, 1, 0));

  test_simulate_stepping_between_faces(Direction::DOWN, 4, 0, 5, Direction::LEFT, std::make_tuple(4, 0, 1));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 4, 5, Direction::LEFT, std::make_tuple(0, 0, 1));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 1, 5, Direction::LEFT, std::make_tuple(3, 0, 1));

  test_simulate_stepping_between_faces(Direction::DOWN, 4, 0, 5, Direction::RIGHT, std::make_tuple(0, 4, 3));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 4, 5, Direction::RIGHT, std::make_tuple(4, 4, 3));
  test_simulate_stepping_between_faces(Direction::DOWN, 4, 1, 5, Direction::RIGHT, std::make_tuple(1, 4, 3));

  std::cout << "Face width is: " << face_width << std::endl;
  auto cube_face_diagram = compress_map_into_cube_face_diagram(matrix, face_width);
  auto [cube, cube_projection] = diagram_into_cube(matrix, cube_face_diagram, face_width);

  for (auto line : cube_projection) {
    dump_to(std::cout, line.cbegin(), line.cend());
    std::cout << std::endl;
  }
  std::cout << std::endl << "Cube face mapping: " << std::endl;
  for (auto const& [key, face] : cube.get_faces()) {
    std::cout << "Face: " << key << ", top left: " << face.top_left_i << " " << face.top_left_j << std::endl;
    std::cout << "\tNeighbors: ";
    dump_to(std::cout, face.neighbors.cbegin(), face.neighbors.cend());
    std::cout << std::endl;
  }

  auto start_i = int{ 0 };
  auto start_j = int{ 0 };
  auto first_empty_at = std::find(matrix[0].cbegin(), matrix[0].cend(), '.');
  start_j = std::distance(matrix[0].cbegin(), first_empty_at);

  auto plane = Plane{ matrix, start_i, start_j, 1 };

  auto directions = std::string{};
  std::getline(std::cin, directions);
  auto dir_data = std::istringstream{ directions };
  auto count = uint32_t{ 0 };

  auto turn = ' ';
  auto read_number = true;
  while (!dir_data.eof()) {
    if (read_number) {
      dir_data >> count;
      plane.simulate_move(count);
      cube.simulate_move(count);

      // std::cout << "Cube state: " << std::endl;
      // for (auto const& line : cube.get_matrix()) {
      //   dump_to(std::cout, line.cbegin(), line.cend());
      //   std::cout << std::endl;
      // }
    } else {
      dir_data >> turn;
      plane.change_direction(turn);
      cube.change_direction(turn);
    }
    read_number = !read_number;
  }

  std::cout << std::endl;
  std::cout << "Part 1" << std::endl;
  std::cout << "Final password: " << calculate_password(plane.get_i(), plane.get_j(), plane.get_offset_index())
            << std::endl;
  std::cout << std::endl;
  std::cout << "Part 2" << std::endl;
  std::cout << cube.get_plane_i() << " " << cube.get_plane_j() << " " << cube.get_offset_index() << std::endl;
  std::cout << "Final password: " << calculate_password(cube.get_plane_i(), cube.get_plane_j(), cube.get_offset_index())
            << std::endl;
  std::cout << std::endl;

  return 0;
}
