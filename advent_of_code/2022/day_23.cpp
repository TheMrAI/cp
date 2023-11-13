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

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto get_rectangle_area(std::vector<std::string> const& buffered_matrix) -> size_t
{
  auto min_i = buffered_matrix.size();
  auto max_i = size_t{ 0 };
  auto min_j = buffered_matrix[0].size();
  auto max_j = size_t{ 0 };
  auto elf_count = size_t{ 0 };
  for (auto i = size_t{ 0 }; i < buffered_matrix.size(); ++i) {
    for (auto j = size_t{ 0 }; j < buffered_matrix[0].size(); ++j) {
      if (buffered_matrix[i][j] == '#') {
        min_i = std::min(min_i, i);
        max_i = std::max(max_i, i);
        min_j = std::min(min_j, j);
        max_j = std::max(max_j, j);

        ++elf_count;
      }
    }
  }

  auto rectangle_area = (1 + max_i - min_i) * (1 + max_j - min_j);
  return rectangle_area - elf_count;
}

struct ProposalInfo
{
  bool valid;
  size_t elf_i;
  size_t elf_j;
};

auto main() -> int
{
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);

  auto matrix = std::vector<std::string>{};
  auto line = std::string{};
  while (std::getline(std::cin, line)) { matrix.push_back(line); }

  auto buffered_matrix = std::vector<std::string>(matrix.size() + 400, std::string(matrix[0].size() + 400, '.'));
  for (auto i = size_t{ 0 }; i < matrix.size(); ++i) {
    for (auto j = size_t{ 0 }; j < matrix[0].size(); ++j) { buffered_matrix[i + 20][j + 20] = matrix[i][j]; }
  }

  constexpr auto considerations = std::array<char, 4>{ 'N', 'S', 'W', 'E' };
  auto considerations_start = size_t{ 0 };
  constexpr auto check_for_direction = std::array<std::array<std::pair<int, int>, 3>, 4>{
    std::array<std::pair<int, int>, 3>{ std::make_pair(-1, -1), std::make_pair(-1, 0), std::make_pair(-1, 1) },
    std::array<std::pair<int, int>, 3>{ std::make_pair(1, -1), std::make_pair(1, 0), std::make_pair(1, 1) },
    std::array<std::pair<int, int>, 3>{ std::make_pair(-1, -1), std::make_pair(0, -1), std::make_pair(1, -1) },
    std::array<std::pair<int, int>, 3>{ std::make_pair(-1, 1), std::make_pair(0, 1), std::make_pair(1, 1) }
  };
  constexpr auto proposal_offsets = std::array<std::pair<int, int>, 4>{
    std::make_pair(-1, 0), std::make_pair(1, 0), std::make_pair(0, -1), std::make_pair(0, 1)
  };

  constexpr auto check_for_action = std::array<std::pair<int, int>, 8>{ std::make_pair(0, -1),
    std::make_pair(-1, -1),
    std::make_pair(-1, 0),
    std::make_pair(-1, 1),
    std::make_pair(0, 1),
    std::make_pair(1, 1),
    std::make_pair(1, 0),
    std::make_pair(1, -1) };

  auto simulate = true;
  auto simulation_round = size_t{ 0 };
  while (simulate) {
    // oh yes make that map as ugly it can be
    auto proposals = std::map<std::pair<int, int>, ProposalInfo>{};

    for (auto i = size_t{ 0 }; i < buffered_matrix.size(); ++i) {
      for (auto j = size_t{ 0 }; j < buffered_matrix[0].size(); ++j) {
        if (buffered_matrix[i][j] == '.') { continue; }
        auto skip =
          std::all_of(check_for_action.cbegin(), check_for_action.cend(), [&buffered_matrix, &i, &j](auto offset) {
            return buffered_matrix[i + offset.first][j + offset.second] == '.';
          });
        if (skip) { continue; }

        for (auto dir_index = considerations_start; dir_index < considerations_start + considerations.size();
             ++dir_index) {
          auto index = dir_index % considerations.size();
          auto all_clear = std::all_of(check_for_direction[index].cbegin(),
            check_for_direction[index].cend(),
            [&buffered_matrix, i, j](
              auto offset) { return buffered_matrix[i + offset.first][j + offset.second] == '.'; });
          if (all_clear) {
            auto proposal = std::make_pair(i + proposal_offsets[index].first, j + proposal_offsets[index].second);
            if (proposals.contains(proposal)) {
              proposals[proposal].valid = false;
              break;
            }
            proposals[proposal] = ProposalInfo{ true, i, j };
            break;
          }
        }
      }
    }

    simulate = false;
    for (auto const& proposal : proposals) {
      if (!proposal.second.valid) { continue; }
      simulate = true;
      buffered_matrix[proposal.second.elf_i][proposal.second.elf_j] = '.';
      buffered_matrix[proposal.first.first][proposal.first.second] = '#';
    }

    considerations_start = (considerations_start + 1) % 4;

    ++simulation_round;
    if (simulation_round == 10) {
      std::cout << "Rectangle area after 10 rounds: " << get_rectangle_area(buffered_matrix) << std::endl;
    }
  }

  std::cout << "Elfs stopped moving on round: " << simulation_round << std::endl;

  return 0;
}
