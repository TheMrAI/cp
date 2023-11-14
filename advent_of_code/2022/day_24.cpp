#include <bits/stdc++.h>

auto get_horizontal_blizzards(std::vector<std::string> const& board) -> std::vector<std::string>
{
  auto horizontal_blizzards = std::vector<std::string>(board.size() - 2, std::string{});
  for (auto i = size_t{ 1 }; i < board.size() - 1; ++i) {
    horizontal_blizzards[i - 1] = std::string{ std::next(board[i].cbegin()), std::prev(board[i].cend()) };
    std::replace_if(
      horizontal_blizzards[i - 1].begin(),
      horizontal_blizzards[i - 1].end(),
      [](auto c) { return c == 'v' || c == '^'; },
      '.');
  }

  return horizontal_blizzards;
}

auto get_vertical_blizzards(std::vector<std::string> const& board) -> std::vector<std::string>
{
  auto vertical_blizzards = std::vector<std::string>{ board[0].size() - 2, std::string(board.size() - 2, '.') };
  for (auto j = size_t{ 1 }; j < board[0].size() - 1; ++j) {
    for (auto i = size_t{ 1 }; i < board.size() - 1; ++i) {
      if (board[i][j] == '<' || board[i][j] == '>') { continue; }
      vertical_blizzards[j - 1][i - 1] = board[i][j];
    }
  }

  return vertical_blizzards;
}

auto horizontal_blizzard_hits(int i, int j, size_t minute, std::vector<std::string> const& horizontal_blizzards) -> bool
{
  // if out of bounds no blizzard can hit no?
  if (i < 0 || i >= horizontal_blizzards.size() || j < 0 || j >= horizontal_blizzards[0].size()) { return false; }
  auto left_drifting_simulated_blizzard_index = (j + minute) % horizontal_blizzards[0].size();
  minute = minute % horizontal_blizzards[0].size();
  minute = horizontal_blizzards[0].size() - minute;
  auto right_drifting_simulated_blizzard_index = (j + minute) % horizontal_blizzards[0].size();
  return horizontal_blizzards[i][right_drifting_simulated_blizzard_index] == '>'
         || horizontal_blizzards[i][left_drifting_simulated_blizzard_index] == '<';
}

auto vertical_blizzard_hits(int j, int i, size_t minute, std::vector<std::string> const& vertical_blizzards) -> bool
{
  if (j < 0 || j >= vertical_blizzards.size() || i < 0 || i >= vertical_blizzards[0].size()) { return false; }
  auto up_drifting_simulated_blizzard_index = (i + minute) % vertical_blizzards[0].size();
  minute = minute % vertical_blizzards[0].size();
  minute = vertical_blizzards[0].size() - minute;
  auto down_drifting_simulated_blizzard_index = (i + minute) % vertical_blizzards[0].size();
  return vertical_blizzards[j][down_drifting_simulated_blizzard_index] == 'v'
         || vertical_blizzards[j][up_drifting_simulated_blizzard_index] == '^';
}

struct State
{
  int i;
  int j;
  int minute;
  int priority;// manhattan distance to goal - minutes
  // the goal is to priority states that are closer to the goal, but if they fly away in time
  // we should prioritize states that are farther away
  int absolute_best_arrival_time;
};

auto get_state_identifier(const State& s) -> std::string
{
  return std::to_string(s.i) + "," + std::to_string(s.j) + "," + std::to_string(s.minute);
}

auto get_manhattan_distance(int i, int j, std::pair<int, int> const& goal) -> int
{
  return std::abs(goal.first - i) + std::abs(goal.second - j);
}

auto initialize_state(int i, int j, int minute, std::pair<int, int> const& goal) -> State
{
  auto distance = get_manhattan_distance(i, j, goal);
  auto arrival_time = distance + minute;
  auto priority = (distance * 30) + minute;
  return State{ i, j, minute, priority, arrival_time };
}

auto find_shortest_path_in_minutes(std::vector<std::string> const& board,
  std::vector<std::string> const& horizontal_blizzards,
  std::vector<std::string> const& vertical_blizzards,
  int starting_at_minute,
  std::pair<int, int> start,
  std::pair<int, int> goal) -> size_t
{
  auto state_comparator = [](State const& lhs, State const& rhs) { return lhs.priority > rhs.priority; };
  auto states = std::priority_queue<State, std::vector<State>, decltype(state_comparator)>{ state_comparator };
  states.push(initialize_state(start.first, start.second, starting_at_minute, goal));

  constexpr auto offsets = std::array<std::pair<int, int>, 5>{
    std::make_pair(0, 0), std::make_pair(1, 0), std::make_pair(-1, 0), std::make_pair(0, 1), std::make_pair(0, -1)
  };
  auto current_best_time = std::numeric_limits<int>::max();
  // identifies states by i, j and minutes
  // we only ever need to queue a state only once
  auto queued_states = std::unordered_set<std::string>{};

  while (!states.empty()) {
    auto to_check = states.top();
    states.pop();

    if (to_check.i == goal.first && to_check.j == goal.second) {
      current_best_time = std::min(current_best_time, to_check.minute);
      // std::cout << current_best_time << std::endl;
      continue;
    }

    for (auto const& offset : offsets) {
      // usually I put these at the top of the while loop for easier readability, but in this case
      // we have to save as much memory as we can, so the states won't even get queued if they are incorrect or already
      // have been processed
      auto state_candidate =
        initialize_state(to_check.i + offset.first, to_check.j + offset.second, to_check.minute + 1, goal);
      auto state_candidate_id = get_state_identifier(state_candidate);

      if (state_candidate.i < 0 || state_candidate.i >= board.size() || state_candidate.j < 0
          || state_candidate.j >= board[0].size() || board[state_candidate.i][state_candidate.j] == '#'
          || queued_states.contains(state_candidate_id)
          || current_best_time <= state_candidate.absolute_best_arrival_time) {
        continue;
      }

      auto adjusted_i = state_candidate.i - 1;
      auto adjusted_j = state_candidate.j - 1;

      if (horizontal_blizzard_hits(adjusted_i, adjusted_j, state_candidate.minute, horizontal_blizzards)
          || vertical_blizzard_hits(adjusted_j, adjusted_i, state_candidate.minute, vertical_blizzards)) {
        continue;
      }
      states.push(state_candidate);
      queued_states.insert(state_candidate_id);
    }
  }
  // this of course shouldn't happen
  return current_best_time;
}

auto main() -> int
{
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);

  auto board = std::vector<std::string>{};
  auto line = std::string{};
  while (std::getline(std::cin, line)) { board.push_back(line); }

  for (auto line : board) { std::cout << line << std::endl; }

  auto horizontal_blizzards = get_horizontal_blizzards(board);
  auto vertical_blizzards = get_vertical_blizzards(board);

  auto start = std::make_pair(0, 1);
  auto goal = std::make_pair(static_cast<int>(board.size() - 1), static_cast<int>(board[0].size() - 2));

  std::cout << "Part 1" << std::endl;
  std::cout << "Time to reach exit: "
            << find_shortest_path_in_minutes(board, horizontal_blizzards, vertical_blizzards, 0, start, goal)
            << std::endl;

  std::cout << std::endl;

  std::cout << "Part 2, there and back again, and there again, a hobbits tale" << std::endl;
  auto time = find_shortest_path_in_minutes(board, horizontal_blizzards, vertical_blizzards, 0, start, goal);
  time = find_shortest_path_in_minutes(board, horizontal_blizzards, vertical_blizzards, time, goal, start);
  time = find_shortest_path_in_minutes(board, horizontal_blizzards, vertical_blizzards, time, start, goal);
  std::cout << "Total time: " << time << std::endl;

  return 0;
}

// simple brute force queue
// time ./day_24 < day_24.long_e_in
// Time to reach exit: 47

// real    0m36.489s
// user    0m32.468s
// sys     0m4.000s

// with priority queue (still too slow)
// time ./day_24 < day_24.long_e_in

// real    0m0.595s
// user    0m0.595s
// sys     0m0.000s

// with tracking which state was visited
// time ./day_24 < day_24.long_e_in

// real    0m0.003s
// user    0m0.003s
// sys     0m0.000s

// dangit! the above was already enough, but I hacked an example from the final one, which was impossible
// i just didn't know about it, so it would run forever :'(
// because of that I would keep looking for bugs
// anyways:
// mrai@potatoland:~/tinker/cp/advent_of_code/2022$ time ./day_24 < day_24.in

// Time to reach exit: 322

// real    0m0.077s
// user    0m0.077s
// sys     0m0.000s